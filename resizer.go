package gentable

import (
	"math"
	"sort"

	"github.com/charmbracelet/lipgloss"
)

// resizer resizes the table to fit a specified width.
type resizer struct {
	tableWidth   int
	columns      []resizerColumn
	styleFunc    StyleFunc
	columnWidths []int
}

// resizerColumn is a column in the resizer.
type resizerColumn struct {
	index      int
	min        int
	max        int
	median     int
	cellWidths []int
	xPadding   int // horizontal padding
	fixedWidth int
}

// newResizer creates a new resizer.
func newResizer(headers ...string) *resizer {
	r := &resizer{}
	r.resize(true, headers)
	return r
}

// resize resizes the table to fit the specified width.
//
// Given a user defined table width, we must ensure the table is exactly that
// width. This must account for all borders, column, separators, and column
// data.
//
// In the case where the table is narrower than the specified table width,
// we simply expand the columns evenly to fit the width.
// For example, a table with 3 columns takes up 50 characters total, and the
// width specified is 80, we expand each column by 10 characters, adding 30
// to the total width.
//
// In the case where the table is wider than the specified table width, we
// _could_ simply shrink the columns evenly but this would result in data
// being truncated (perhaps unnecessarily). The naive approach could result
// in very poor cropping of the table data. So, instead of shrinking columns
// evenly, we calculate the median non-whitespace length of each column, and
// shrink the columns based on the largest median.
//
// For example,
//
//	┌──────┬───────────────┬──────────┐
//	│ Name │ Age of Person │ Location │
//	├──────┼───────────────┼──────────┤
//	│ Kini │ 40            │ New York │
//	│ Eli  │ 30            │ London   │
//	│ Iris │ 20            │ Paris    │
//	└──────┴───────────────┴──────────┘
//
// Median non-whitespace length  vs column width of each column:
//
// Name: 4 / 5
// Age of Person: 2 / 15
// Location: 6 / 10
//
// The biggest difference is 15 - 2, so we can shrink the 2nd column by 13.
func (t *resizer) resize(isHeaders bool, rows ...[]string) {
	// Update cell widths for each column and update the min and max cell
	// length.
	for _, row := range rows {
		for i, cell := range row {
			cellLen := lipgloss.Width(cell)

			if isHeaders {
				t.columns = append(t.columns, resizerColumn{
					index:  i,
					min:    cellLen,
					max:    cellLen,
					median: cellLen,
				})
				continue
			}

			t.columns[i].cellWidths = append(t.columns[i].cellWidths, cellLen)
			t.columns[i].min = min(t.columns[i].min, cellLen)
			t.columns[i].max = max(t.columns[i].max, cellLen)
		}
	}

	// Update median cell width for each column.
	//
	// TODO: this is inefficient because the median is being recalculated by
	// re-sorting all the cell widths each and every time. There is probably a
	// more efficient algorithm.
	for i := range t.columns {
		t.columns[i].median = median(t.columns[i].cellWidths)
	}

	// Update fixed-width and horizontal padding for each column.
	for _, row := range rows {
		for j := range row {
			column := &t.columns[j]

			style := lipgloss.NewStyle()
			if j != len(t.columns)-1 {
				style = style.PaddingRight(1)
			}

			column.xPadding = max(column.xPadding, style.GetHorizontalFrameSize())
			column.fixedWidth = max(column.fixedWidth, style.GetWidth())
		}
	}

	// A table width wasn't specified. In this case, detect according to
	// content width.
	if t.tableWidth <= 0 {
		t.tableWidth = t.detectTableWidth()
	}

	// Update column widths.
	t.columnWidths = t.optimizedWidths()
}

// Width sets the table width.
func (r *resizer) Width(width int) {
	r.tableWidth = width
	r.resize(false)
}

// detectTableWidth detects the table width.
func (r *resizer) detectTableWidth() int {
	return r.maxCharCount() + r.totalHorizontalPadding()
}

// maxCharCount returns the maximum character count.
func (r *resizer) maxCharCount() int {
	var count int
	for _, col := range r.columns {
		if col.fixedWidth > 0 {
			count += col.fixedWidth - r.xPaddingForCol(col.index)
		} else {
			count += col.max
		}
	}
	return count
}

// totalHorizontalPadding returns the total padding.
func (r *resizer) totalHorizontalPadding() (totalHorizontalPadding int) {
	for _, col := range r.columns {
		totalHorizontalPadding += col.xPadding
	}
	return
}

func (r *resizer) optimizedWidths() (colWidths []int) {
	if r.maxTotal() <= r.tableWidth {
		return r.expandTableWidth()
	}
	return r.shrinkTableWidth()
}

// expandTableWidth expands the table width.
func (r *resizer) expandTableWidth() (colWidths []int) {
	colWidths = r.maxColumnWidths()

	for {
		totalWidth := sum(colWidths)
		if totalWidth >= r.tableWidth {
			break
		}

		shorterColumnIndex := 0
		shorterColumnWidth := math.MaxInt32

		for j, width := range colWidths {
			if width == r.columns[j].fixedWidth {
				continue
			}
			if width < shorterColumnWidth {
				shorterColumnWidth = width
				shorterColumnIndex = j
			}
		}

		colWidths[shorterColumnIndex]++
	}

	return
}

// shrinkTableWidth shrinks the table width.
func (r *resizer) shrinkTableWidth() (colWidths []int) {
	colWidths = r.maxColumnWidths()

	// Cut width of columns that are way too big.
	shrinkBiggestColumns := func(veryBigOnly bool) {
		for {
			totalWidth := sum(colWidths)
			if totalWidth <= r.tableWidth {
				break
			}

			bigColumnIndex := -math.MaxInt32
			bigColumnWidth := -math.MaxInt32

			for j, width := range colWidths {
				if width == r.columns[j].fixedWidth {
					continue
				}
				if veryBigOnly {
					if width >= (r.tableWidth/2) && width > bigColumnWidth { //nolint:mnd
						bigColumnWidth = width
						bigColumnIndex = j
					}
				} else {
					if width > bigColumnWidth {
						bigColumnWidth = width
						bigColumnIndex = j
					}
				}
			}

			if bigColumnIndex < 0 || colWidths[bigColumnIndex] == 0 {
				break
			}
			colWidths[bigColumnIndex]--
		}
	}

	// Cut width of columns that differ the most from the median.
	shrinkToMedian := func() {
		for {
			totalWidth := sum(colWidths)
			if totalWidth <= r.tableWidth {
				break
			}

			biggestDiffToMedian := -math.MaxInt32
			biggestDiffToMedianIndex := -math.MaxInt32

			for j, width := range colWidths {
				if width == r.columns[j].fixedWidth {
					continue
				}
				diffToMedian := width - r.columns[j].median
				if diffToMedian > 0 && diffToMedian > biggestDiffToMedian {
					biggestDiffToMedian = diffToMedian
					biggestDiffToMedianIndex = j
				}
			}

			if biggestDiffToMedianIndex <= 0 || colWidths[biggestDiffToMedianIndex] == 0 {
				break
			}
			colWidths[biggestDiffToMedianIndex]--
		}
	}

	shrinkBiggestColumns(true)
	shrinkToMedian()
	shrinkBiggestColumns(false)

	return colWidths
}

func (r *resizer) maxColumnWidths() []int {
	maxColumnWidths := make([]int, len(r.columns))
	for i, col := range r.columns {
		if col.fixedWidth > 0 {
			maxColumnWidths[i] = col.fixedWidth
		} else {
			maxColumnWidths[i] = col.max + r.xPaddingForCol(col.index)
		}
	}
	return maxColumnWidths
}

// maxTotal returns the maximum total width.
func (r *resizer) maxTotal() (maxTotal int) {
	for j, column := range r.columns {
		if column.fixedWidth > 0 {
			maxTotal += column.fixedWidth
		} else {
			maxTotal += column.max + r.xPaddingForCol(j)
		}
	}
	return
}

// xPaddingForCol returns the horizontal padding for a column.
func (r *resizer) xPaddingForCol(j int) int {
	if j >= len(r.columns) {
		return 0
	}
	return r.columns[j].xPadding
}

// sum returns the sum of all integers in a slice.
func sum(n []int) int {
	var sum int
	for _, i := range n {
		sum += i
	}
	return sum
}

// median returns the median of a slice of integers.
func median(n []int) int {
	sort.Ints(n)

	if len(n) <= 0 {
		return 0
	}
	if len(n)%2 == 0 {
		h := len(n) / 2            //nolint:mnd
		return (n[h-1] + n[h]) / 2 //nolint:mnd
	}
	return n[len(n)/2]
}
