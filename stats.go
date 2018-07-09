package logs

// Stats services provide optional statistics monitoring to the logger.
type Stats interface {
	// Count adds the given amount to the named stat.
	Count(name string, amount int64)
	// Value sets the given stat to a new value.
	Value(name string, amount float64)
}
