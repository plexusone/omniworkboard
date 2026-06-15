// Package omniworkboard provides project management workboard functionality.
package omniworkboard

import (
	"time"
)

// Priority represents the priority level of a card.
type Priority int

const (
	// PriorityLow is for non-urgent tasks.
	PriorityLow Priority = iota
	// PriorityNormal is the default priority.
	PriorityNormal
	// PriorityHigh is for urgent tasks.
	PriorityHigh
	// PriorityCritical is for blocking or time-sensitive tasks.
	PriorityCritical
)

// String returns the string representation of the priority.
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityCritical:
		return "critical"
	default:
		return "normal"
	}
}

// ParsePriority parses a string into a Priority.
func ParsePriority(s string) Priority {
	switch s {
	case "low":
		return PriorityLow
	case "normal":
		return PriorityNormal
	case "high":
		return PriorityHigh
	case "critical":
		return PriorityCritical
	default:
		return PriorityNormal
	}
}

// Card represents a task card on the workboard.
type Card struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// Title is the card title.
	Title string `json:"title"`

	// Description provides detailed information.
	Description string `json:"description,omitempty"`

	// Column is the current column the card is in.
	Column ColumnType `json:"column"`

	// Priority is the card priority.
	Priority Priority `json:"priority"`

	// DependsOn is a list of card IDs this card depends on.
	DependsOn []string `json:"depends_on,omitempty"`

	// BlockedBy is computed from DependsOn - cards that block this one.
	// This is populated by the board when querying.
	BlockedBy []string `json:"blocked_by,omitempty"`

	// Labels are tags for categorization.
	Labels []string `json:"labels,omitempty"`

	// Assignee is the person or agent assigned to this card.
	Assignee string `json:"assignee,omitempty"`

	// Metadata contains additional key-value data.
	Metadata map[string]any `json:"metadata,omitempty"`

	// CreatedAt is when the card was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when the card was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// CompletedAt is when the card was moved to done.
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// IsBlocked returns true if the card has unfinished dependencies.
func (c *Card) IsBlocked() bool {
	return len(c.BlockedBy) > 0
}

// HasLabel returns true if the card has the given label.
func (c *Card) HasLabel(label string) bool {
	for _, l := range c.Labels {
		if l == label {
			return true
		}
	}
	return false
}

// AddLabel adds a label if not already present.
func (c *Card) AddLabel(label string) {
	if !c.HasLabel(label) {
		c.Labels = append(c.Labels, label)
	}
}

// RemoveLabel removes a label if present.
func (c *Card) RemoveLabel(label string) {
	for i, l := range c.Labels {
		if l == label {
			c.Labels = append(c.Labels[:i], c.Labels[i+1:]...)
			return
		}
	}
}

// AddDependency adds a dependency on another card.
func (c *Card) AddDependency(cardID string) {
	for _, d := range c.DependsOn {
		if d == cardID {
			return // Already exists
		}
	}
	c.DependsOn = append(c.DependsOn, cardID)
}

// RemoveDependency removes a dependency.
func (c *Card) RemoveDependency(cardID string) {
	for i, d := range c.DependsOn {
		if d == cardID {
			c.DependsOn = append(c.DependsOn[:i], c.DependsOn[i+1:]...)
			return
		}
	}
}
