package omniworkboard

// ColumnType represents the type of column on a workboard.
type ColumnType string

const (
	// ColumnBacklog is for items not yet planned.
	ColumnBacklog ColumnType = "backlog"

	// ColumnTodo is for planned items ready to start.
	ColumnTodo ColumnType = "todo"

	// ColumnInProgress is for items currently being worked on.
	ColumnInProgress ColumnType = "in_progress"

	// ColumnReview is for items awaiting review.
	ColumnReview ColumnType = "review"

	// ColumnDone is for completed items.
	ColumnDone ColumnType = "done"
)

// String returns the string representation.
func (c ColumnType) String() string {
	return string(c)
}

// IsValid returns true if the column type is valid.
func (c ColumnType) IsValid() bool {
	switch c {
	case ColumnBacklog, ColumnTodo, ColumnInProgress, ColumnReview, ColumnDone:
		return true
	default:
		return false
	}
}

// IsDone returns true if the column represents completion.
func (c ColumnType) IsDone() bool {
	return c == ColumnDone
}

// IsActive returns true if the column represents active work.
func (c ColumnType) IsActive() bool {
	return c == ColumnInProgress || c == ColumnReview
}

// ParseColumnType parses a string into a ColumnType.
func ParseColumnType(s string) ColumnType {
	ct := ColumnType(s)
	if ct.IsValid() {
		return ct
	}
	return ColumnBacklog
}

// DefaultColumns returns the default set of columns in order.
func DefaultColumns() []ColumnType {
	return []ColumnType{
		ColumnBacklog,
		ColumnTodo,
		ColumnInProgress,
		ColumnReview,
		ColumnDone,
	}
}

// ColumnIndex returns the position of a column in the workflow.
func ColumnIndex(col ColumnType) int {
	for i, c := range DefaultColumns() {
		if c == col {
			return i
		}
	}
	return -1
}

// CanTransition checks if a transition between columns is allowed.
func CanTransition(from, to ColumnType) bool {
	fromIdx := ColumnIndex(from)
	toIdx := ColumnIndex(to)

	if fromIdx < 0 || toIdx < 0 {
		return false
	}

	// Allow forward transitions (standard workflow)
	if toIdx > fromIdx {
		return true
	}

	// Allow moving back to todo (e.g., for re-work)
	if to == ColumnTodo && from != ColumnBacklog {
		return true
	}

	// Allow moving back to in_progress from review
	if to == ColumnInProgress && from == ColumnReview {
		return true
	}

	return false
}
