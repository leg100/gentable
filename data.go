package gentable

// data provides access to the underlying rows
type data[V comparable] interface {
	// Rows returns the row values
	Rows() []V
	// Columns returns the number of columns in the table.
	Columns() int
	// Append rows to the unwindowed rows.
	Append(...row[V])
	getCells(v V) []string
}
