package gentable

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/maps"
)

const (
	// Height of the table header
	headerHeight = 1
	// Minimum recommended height for the table widget. Respecting this minimum
	// ensures the header and the borders and the filter widget are visible.
	MinHeight = 6
)

// Model is a bubbletea model for a table component.
type Model[V any] struct {
	getID    func(V) ID
	cols     []Column
	rows     []V
	render   func(V) RenderedCells
	rendered map[ID]RenderedCells

	currentRowIndex int
	currentRowID    ID

	// items are the unfiltered set of items available to the table.
	items    map[ID]V
	sortFunc SortFunc[V]

	selected   map[ID]V
	selectable bool
	scrollbar  bool
	border     bool

	filter *filter

	// index of first visible row
	start int
	// width of table without borders
	width int
	// height of table without borders
	height int
}

// ID uniquely identifies a table item
type ID any

// Column defines the table structure.
type Column struct {
	Key ColumnKey
	// TODO: Default to upper case of key
	Title          string
	TruncationFunc func(s string, w int, tail string) string
	// RightAlign aligns content to the right. If false, content is aligned to
	// the left.
	RightAlign bool
}

type ColumnKey string

// RenderedCells provides the rendered string for each column in a row.
type RenderedCells map[ColumnKey]string

type SortFunc[V any] func(V, V) int

// New creates a new model for the table widget.
func New[V any](
	idfn func(V) ID,
	renderfn func(V) RenderedCells,
	width, height int,
	opts ...Option[V],
) Model[V] {
	m := Model[V]{
		render:          renderfn,
		items:           make(map[ID]V),
		rendered:        make(map[ID]RenderedCells),
		selected:        make(map[ID]V),
		getID:           idfn,
		selectable:      true,
		filter:          newFilter(),
		currentRowIndex: -1,
	}
	for _, fn := range opts {
		fn(&m)
	}

	m.SetDimensions(width, height)

	return m
}

// Init initializes the table.
//
// NOTE: this currently does nothing but is implemented in order to satisfy
// tea.Model for user's convenience.
func (m Model[V]) Init() tea.Cmd {
	return nil
}

// Update is the Bubble Tea update loop.
func (m Model[V]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.filter.Focused() {
			// keys are diverted to filter when it's focused
			break
		}
		switch {
		case key.Matches(msg, keys.LineUp):
			m.MoveUp(1)
		case key.Matches(msg, keys.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, keys.PageUp):
			m.MoveUp(m.maxVisibleRowsCount())
		case key.Matches(msg, keys.PageDown):
			m.MoveDown(m.maxVisibleRowsCount())
		case key.Matches(msg, keys.HalfPageUp):
			m.MoveUp(m.maxVisibleRowsCount() / 2)
		case key.Matches(msg, keys.HalfPageDown):
			m.MoveDown(m.maxVisibleRowsCount() / 2)
		case key.Matches(msg, keys.GotoTop):
			m.GotoTop()
		case key.Matches(msg, keys.GotoBottom):
			m.GotoBottom()
		case key.Matches(msg, keys.Select):
			m.ToggleSelection()
		case key.Matches(msg, keys.SelectAll):
			m.SelectAll()
		case key.Matches(msg, keys.SelectClear):
			m.DeselectAll()
		case key.Matches(msg, keys.SelectRange):
			m.SelectRange()
		case key.Matches(msg, keys.Filter):
			// Focus the filter widget
			blink := m.filter.Focus()
			// Start blinking the cursor.
			return m, blink
		}
	}
	updateItems, cmd := m.filter.Update(msg)
	if updateItems {
		m.setRows(maps.Values(m.items)...)
	}
	return m, cmd
}

// View renders the table.
func (m Model[V]) View() string {
	// Calculate widths of each column.
	//
	// TODO: It's debatable whether this should be done here at render-time, or
	// when the width is changed or when the visible rows are changed (since
	// it's based only on visible rows not all rows). I think it's a trade off
	// between performance and code-complexity.
	//
	// cellWidths are the widths of each cell, arranged column then by row.
	cellWidths := make([][]int, len(m.cols))
	// Gather cell widths
	for i, col := range m.cols {
		// Allocate nested slice with space for column title and rows.
		cellWidths[i] = make([]int, 1+len(m.visibleRows()))
		// Gather column title width
		cellWidths[i][0] = len(col.Title)
		// Gather width of cell for each row.
		for j, row := range m.visibleRows() {
			cells := m.rendered[m.getID(row)]
			if cell, ok := cells[col.Key]; ok {
				cellWidths[i][j+1] = lipgloss.Width(cell)
			}
		}
	}
	// Calculate widths and median string lengths for each column
	var (
		widths  = make([]int, len(m.cols))
		medians = make([]int, len(m.cols))
	)
	for i := range m.cols {
		widths[i] = slices.Max(cellWidths[i])
		medians[i] = median(cellWidths[i])
	}
	// Determine total width of all columns, using the largest string found in
	// each column, and account for padding too.
	totalWidth := sum(widths) + len(m.cols)*2
	// Now determine if that sum is narrower or wider than model width and make
	// adjustments to column widths accordingly.
	if m.width == 0 {
		// Perform no width adjustment if user hasn't specified a width
	} else if totalWidth < m.width {
		// There is additional width; expand columns evenly until their sum
		// matches the model width
		var i int
		for totalWidth < m.width {
			widths[i]++
			totalWidth++
			i = (i + 1) % len(widths)
		}
	} else if totalWidth > m.width {
		// Sum of widest strings for each column is wider than model width.
		//
		// Find the biggest differences between the median and the column width.
		// Shrink the columns based on the largest difference.
		differences := make([]int, len(widths))
		for i := range widths {
			differences[i] = widths[i] - medians[i]
		}
		for totalWidth > m.width {
			i, _ := largest(differences)
			if differences[i] < 1 {
				// No difference between largest and median, so quit this
				// approach.
				break
			}
			shrink := min(differences[i], totalWidth-m.width)
			widths[i] -= shrink
			m.width -= shrink
			differences[i] = 0
		}
		// Table is still too wide, begin shrinking the columns based on the
		// largest column.
		for m.width > totalWidth {
			i, _ := largest(widths)
			if widths[i] < 1 {
				break
			}
			widths[i]--
			totalWidth--
		}
	}

	var builder strings.Builder
	// Filter is at the top
	if m.filter.isVisible() {
		builder.WriteString(
			Regular.Margin(0, 1).Render(m.filter.View()),
		)
		builder.WriteRune('\n')
		// Add horizontal rule between filter widget and table
		builder.WriteString(
			strings.Repeat("─", m.width),
		)
	}
	// Then the headers
	{
		headers := make([]string, len(m.cols))
		for i, col := range m.cols {
			style := lipgloss.NewStyle().Width(widths[i]).MaxWidth(widths[i]).Inline(true)
			if col.RightAlign {
				style = style.AlignHorizontal(lipgloss.Right)
			}
			renderedCell := style.Render(TruncateRight(col.Title, widths[i], "…"))
			headers[i] = Regular.Padding(0, 1).Render(renderedCell)
		}
		builder.WriteString(
			lipgloss.NewStyle().
				MaxWidth(m.width).
				Render(lipgloss.JoinHorizontal(lipgloss.Left, headers...)),
		)
	}
	// Then the rows
	builder.WriteRune('\n')
	{
		rendered := make([]string, len(m.visibleRows()))
		for i, row := range m.visibleRows() {
			var (
				id         = m.getID(row)
				background lipgloss.Color
				foreground lipgloss.Color
				current    bool
				selected   bool
			)
			if _, ok := m.selected[id]; ok {
				selected = true
			}
			current = id == m.currentRowID
			if current && selected {
				background = CurrentAndSelectedBackground
				foreground = CurrentAndSelectedForeground
			} else if current {
				background = CurrentBackground
				foreground = CurrentForeground
			} else if selected {
				background = SelectedBackground
				foreground = SelectedForeground
			}

			cells := m.rendered[id]
			styledCells := make([]string, len(m.cols))
			for i, col := range m.cols {
				content := cells[col.Key]
				// Truncate content if it is wider than column
				truncated := col.TruncationFunc(content, widths[i], "…")
				// Ensure content is all on one line.
				style := lipgloss.NewStyle().
					Width(widths[i]).
					MaxWidth(widths[i]).
					Inline(true)
				if col.RightAlign {
					style = style.AlignHorizontal(lipgloss.Right)
				}
				inlined := style.Render(truncated)
				// Apply block-styling to content
				boxed := lipgloss.NewStyle().
					Padding(0, 1).
					Render(inlined)
				styledCells[i] = boxed
			}

			// Join cells together to form a row, ensuring it doesn't exceed maximum
			// table width
			renderedRow := lipgloss.JoinHorizontal(lipgloss.Left, styledCells...)
			// Join cells together to form a row, ensuring it doesn't exceed maximum
			// table width
			renderedRow = lipgloss.NewStyle().
				MaxWidth(m.width).
				Render(renderedRow)

			// If current row or selected rows, strip colors and apply background color
			if current || selected {
				renderedRow = stripAnsi(renderedRow)
				renderedRow = lipgloss.NewStyle().
					Foreground(foreground).
					Background(background).
					Render(renderedRow)
			}
			rendered[i] = renderedRow
		}
		rowarea := lipgloss.NewStyle().Width(m.width - scrollbarWidth).Render(
			strings.Join(rendered, "\n"),
		)
		if m.scrollbar {
			// Generate scrollbar
			scrollbar := scrollbar(m.maxVisibleRowsCount(), len(m.rows), m.currentVisibleRowsCount(), m.start)
			// Put scrollbar to the right of rows
			rowarea = lipgloss.JoinHorizontal(lipgloss.Top, rowarea, scrollbar)
		}
		builder.WriteString(rowarea)
	}
	// Render table components, ensuring it is at least a min height
	content := lipgloss.NewStyle().
		Height(m.height).
		MaxHeight(m.height).
		Render(builder.String())
	if m.border {
		content = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Render(content)
	}
	return content
}

// SetDimensions sets the dimensions of the table.
func (m *Model[V]) SetDimensions(width, height int) {
	m.height = height
	m.width = width
	if m.border {
		m.height -= 2
		m.width -= 2
	}
	m.setStart()
}

// maxVisibleRowsCount returns the maximum number of visible rows.
func (m Model[V]) maxVisibleRowsCount() int {
	height := max(0, m.height-headerHeight)
	if m.filter.isVisible() {
		// Accommodate height of filter widget
		return max(0, height-filterHeight)
	}
	return height
}

// currentVisibleRowsCount returns the number of currently visible rows. Note
// this may be less than maxVisibleRowsCount().
//
// TODO: move this into visibleRows()
func (m Model[V]) currentVisibleRowsCount() int {
	return min(m.maxVisibleRowsCount(), len(m.rows)-m.start)
}

func (m Model[V]) visibleRows() []V {
	return m.rows[m.start:m.currentVisibleRowsCount()]
}

// Metadata renders a short string summarizing table row metadata.
func (m *Model[V]) Metadata() string {
	var metadata string
	// Calculate the top and bottom visible row ordinal numbers
	top := m.start + 1
	bottom := m.start + m.currentVisibleRowsCount()
	prefix := fmt.Sprintf("%d-%d of ", top, bottom)
	if m.filter.isVisible() {
		metadata = prefix + fmt.Sprintf("%d/%d", len(m.rows), len(m.items))
	} else {
		metadata = prefix + strconv.Itoa(len(m.rows))
	}
	return metadata
}

// CurrentRow returns the current row the user has highlighted.  If the table is
// empty then false is returned.
func (m Model[V]) CurrentRow() (V, bool) {
	if m.currentRowIndex < 0 || m.currentRowIndex >= len(m.rows) {
		return *new(V), false
	}
	return m.rows[m.currentRowIndex], true
}

// SelectedOrCurrent returns either the selected rows, or if there are no
// selections, the current row
func (m Model[V]) SelectedOrCurrent() []V {
	if len(m.selected) > 0 {
		rows := make([]V, len(m.selected))
		var i int
		for _, v := range m.selected {
			rows[i] = v
			i++
		}
		return rows
	}
	if row, ok := m.CurrentRow(); ok {
		return []V{row}
	}
	return nil
}

// ToggleSelection toggles the selection of the current row.
func (m *Model[V]) ToggleSelection() {
	if !m.selectable {
		return
	}
	current, ok := m.CurrentRow()
	if !ok {
		return
	}
	if _, isSelected := m.selected[current]; isSelected {
		delete(m.selected, current)
	} else {
		m.selected[m.getID(current)] = current
	}
}

// ToggleSelectionByID toggles the selection of the row with the given ID. If
// the ID does not exist no action is taken.
func (m *Model[V]) ToggleSelectionByID(id ID) {
	if !m.selectable {
		return
	}
	v, ok := m.items[id]
	if !ok {
		return
	}
	if _, isSelected := m.selected[id]; isSelected {
		delete(m.selected, id)
	} else {
		m.selected[id] = v
	}
}

// SelectAll selects all rows. Any rows not currently selected are selected.
func (m *Model[V]) SelectAll() {
	if !m.selectable {
		return
	}
	for _, row := range m.rows {
		m.selected[row] = row
	}
}

// DeselectAll de-selects any rows that are currently selected
func (m *Model[V]) DeselectAll() {
	if !m.selectable {
		return
	}
	m.selected = make(map[ID]V)
}

// SelectRange selects a range of rows. If the current row is *below* a selected
// row then rows between them are selected, including the current row.
// Otherwise, if the current row is *above* a selected row then rows between
// them are selected, including the current row. If there are no selected rows
// then no action is taken.
func (m *Model[V]) SelectRange() {
	if !m.selectable {
		return
	}
	if len(m.selected) == 0 {
		return
	}
	// Determine the first row to select, and the number of rows to select.
	first := -1
	n := 0
	for i, row := range m.rows {
		if i == m.currentRowIndex && first > -1 && first < m.currentRowIndex {
			// Select rows before and including current row
			n = m.currentRowIndex - first + 1
			break
		}
		if _, ok := m.selected[m.getID(row)]; !ok {
			// Ignore unselected rows
			continue
		}
		if i > m.currentRowIndex {
			// Select rows including current row and all rows up to but not
			// including next selected row
			first = m.currentRowIndex
			n = i - m.currentRowIndex
			break
		}
		// Start selecting rows after this currently selected row.
		first = i + 1
	}
	for _, row := range m.rows[first : first+n] {
		m.selected[m.getID(row)] = row
	}
}

// SetItems overwrites all existing items in the table with newItems.
func (m *Model[V]) SetItems(newItems ...V) {
	m.items = make(map[ID]V)
	m.rendered = make(map[ID]RenderedCells)
	m.AddItems(newItems...)
}

// AddItems idempotently adds items to the table, updating any items that exist
// on the table already with a matching ID.
func (m *Model[V]) AddItems(items ...V) {
	for _, item := range items {
		// Add/update item
		m.items[m.getID(item)] = item
		// (Re-)render item
		m.rendered[m.getID(item)] = m.render(item)
	}
	m.setRows(items...)
}

func (m *Model[V]) removeItem(item V) {
	delete(m.rendered, m.getID(item))
	delete(m.items, m.getID(item))
	delete(m.selected, m.getID(item))
	for i, row := range m.rows {
		if m.getID(row) == m.getID(item) {
			// TODO: this might well produce a memory leak. See note:
			// https://go.dev/wiki/SliceTricks#delete-without-preserving-order
			m.rows = append(m.rows[:i], m.rows[i+1:]...)
			break
		}
	}
	if m.getID(item) == m.currentRowID {
		// If item being removed is the current row the make the row above it
		// the new current row. (MoveUp also calls setStart, see below).
		m.MoveUp(1)
	} else {
		// Removing item may well affect index of first visible row, so
		// re-calculate just in case.
		m.setStart()
	}
}

func (m *Model[V]) setRows(items ...V) {
	selected := make(map[ID]V)
	m.rows = make([]V, 0, len(items))
	for _, item := range items {
		if m.filter.isVisible() && !m.matchFilter(item) {
			// Skip item that doesn't match filter
			continue
		}
		m.rows = append(m.rows, item)
		if m.selectable {
			if _, ok := m.selected[item]; ok {
				selected[item] = item
			}
		}
	}
	m.selected = selected
	// Sort rows in-place
	if m.sortFunc != nil {
		slices.SortFunc(m.rows, func(i, j V) int {
			return m.sortFunc(i, j)
		})
	}
	// Track current row index
	m.currentRowIndex = -1
	for i, row := range m.rows {
		if m.getID(row) == m.currentRowID {
			m.currentRowIndex = i
			break
		}
	}
	// Check if item corresponding to current row doesn't exist, which occurs
	// the very first time the table is populated. If so, set current row to the
	// first row.
	if len(m.rows) > 0 && m.currentRowIndex == -1 {
		m.currentRowIndex = 0
		m.currentRowID = m.getID(m.rows[m.currentRowIndex])
	}
	m.setStart()
}

// matchFilter returns true if the item with the given ID matches the filter
// value.
func (m *Model[V]) matchFilter(item V) bool {
	for _, col := range m.rendered[m.getID(item)] {
		// Remove ANSI escapes code before filtering
		stripped := stripAnsi(col)
		if strings.Contains(stripped, m.filter.Value()) {
			return true
		}
	}
	return false
}

// MoveUp moves the current row up by any number of rows.
// It can not go above the first row.
func (m *Model[V]) MoveUp(n int) {
	m.moveCurrentRow(-n)
}

// MoveDown moves the current row down by any number of rows.
// It can not go below the last row.
func (m *Model[V]) MoveDown(n int) {
	m.moveCurrentRow(n)
}

func (m *Model[V]) moveCurrentRow(n int) {
	if len(m.rows) > 0 {
		m.currentRowIndex = clamp(m.currentRowIndex+n, 0, len(m.rows)-1)
		m.currentRowID = m.getID(m.rows[m.currentRowIndex])
		m.setStart()
	}
}

func (m *Model[V]) setStart() {
	// Start index must be at least the current row index minus the max number
	// of visible rows.
	minimum := max(0, m.currentRowIndex-m.maxVisibleRowsCount()+1)
	// Start index is at most the lesser of:
	// (a) the current row index, or
	// (b) the number of rows minus the maximum number of visible rows (as many
	// rows as possible are rendered)
	maximum := max(0, min(m.currentRowIndex, len(m.rows)-m.maxVisibleRowsCount()))
	m.start = clamp(m.start, minimum, maximum)
}

// GotoTop makes the top row the current row.
func (m *Model[V]) GotoTop() {
	m.MoveUp(m.currentRowIndex)
}

// GotoBottom makes the bottom row the current row.
func (m *Model[V]) GotoBottom() {
	m.MoveDown(len(m.rows))
}

func clamp(v, low, high int) int {
	if high < low {
		low, high = high, low
	}
	return min(high, max(low, v))
}

// median returns the median of a slice of integers.
func median(n []int) int {
	sort.Ints(n)

	if len(n) <= 0 {
		return 0
	}
	if len(n)%2 == 0 {
		h := len(n) / 2            //nolint:gomnd
		return (n[h-1] + n[h]) / 2 //nolint:gomnd
	}
	return n[len(n)/2]
}

// sum returns the sum of all integers in a slice.
func sum(n []int) int {
	var sum int
	for _, i := range n {
		sum += i
	}
	return sum
}

// largest returns the largest element and it's index from a slice of integers.
func largest(n []int) (int, int) { //nolint:unparam
	var biggest, index int
	for i, e := range n {
		if n[i] > n[index] {
			biggest = e
			index = i
		}
	}
	return index, biggest
}
