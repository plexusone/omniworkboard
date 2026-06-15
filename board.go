package omniworkboard

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Board represents a workboard with columns and cards.
type Board struct {
	// ID is the unique identifier.
	ID string `json:"id"`

	// Name is the board name.
	Name string `json:"name"`

	// Description is an optional description.
	Description string `json:"description,omitempty"`

	// Columns is the ordered list of columns.
	Columns []ColumnType `json:"columns"`

	// Cards is the map of card ID to card.
	Cards map[string]*Card `json:"cards"`

	// CreatedAt is when the board was created.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when the board was last updated.
	UpdatedAt time.Time `json:"updated_at"`

	// mu protects the board state.
	mu sync.RWMutex
}

// BoardConfig configures a new board.
type BoardConfig struct {
	// Name is the board name.
	Name string

	// Description is an optional description.
	Description string

	// Columns is the ordered list of columns (defaults to DefaultColumns).
	Columns []ColumnType
}

// NewBoard creates a new workboard.
func NewBoard(cfg BoardConfig) *Board {
	columns := cfg.Columns
	if len(columns) == 0 {
		columns = DefaultColumns()
	}

	now := time.Now()
	return &Board{
		ID:          uuid.New().String(),
		Name:        cfg.Name,
		Description: cfg.Description,
		Columns:     columns,
		Cards:       make(map[string]*Card),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// CreateCard creates a new card on the board.
func (b *Board) CreateCard(ctx context.Context, title, description string, priority Priority) (*Card, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	card := &Card{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Column:      ColumnBacklog,
		Priority:    priority,
		CreatedAt:   now,
		UpdatedAt:   now,
		Metadata:    make(map[string]any),
	}

	b.Cards[card.ID] = card
	b.UpdatedAt = now

	return card, nil
}

// GetCard retrieves a card by ID.
func (b *Board) GetCard(ctx context.Context, id string) (*Card, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	card, ok := b.Cards[id]
	if !ok {
		return nil, fmt.Errorf("card %q not found", id)
	}

	// Compute blocked_by
	cardCopy := *card
	cardCopy.BlockedBy = b.getBlockingCards(card)

	return &cardCopy, nil
}

// UpdateCard updates a card's fields.
func (b *Board) UpdateCard(ctx context.Context, id string, updates CardUpdate) (*Card, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	card, ok := b.Cards[id]
	if !ok {
		return nil, fmt.Errorf("card %q not found", id)
	}

	now := time.Now()

	if updates.Title != nil {
		card.Title = *updates.Title
	}
	if updates.Description != nil {
		card.Description = *updates.Description
	}
	if updates.Priority != nil {
		card.Priority = *updates.Priority
	}
	if updates.Assignee != nil {
		card.Assignee = *updates.Assignee
	}
	if updates.Labels != nil {
		card.Labels = updates.Labels
	}
	if updates.Metadata != nil {
		if card.Metadata == nil {
			card.Metadata = make(map[string]any)
		}
		for k, v := range updates.Metadata {
			if v == nil {
				delete(card.Metadata, k)
			} else {
				card.Metadata[k] = v
			}
		}
	}

	card.UpdatedAt = now
	b.UpdatedAt = now

	// Return copy with computed fields
	cardCopy := *card
	cardCopy.BlockedBy = b.getBlockingCards(card)

	return &cardCopy, nil
}

// CardUpdate contains optional fields for updating a card.
type CardUpdate struct {
	Title       *string
	Description *string
	Priority    *Priority
	Assignee    *string
	Labels      []string
	Metadata    map[string]any
}

// MoveCard moves a card to a different column.
func (b *Board) MoveCard(ctx context.Context, id string, toColumn ColumnType) (*Card, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	card, ok := b.Cards[id]
	if !ok {
		return nil, fmt.Errorf("card %q not found", id)
	}

	// Check if card is blocked
	blockers := b.getBlockingCards(card)
	if len(blockers) > 0 && toColumn != ColumnBacklog && toColumn != card.Column {
		return nil, fmt.Errorf("card %q is blocked by: %v", id, blockers)
	}

	// Check valid transition
	if !CanTransition(card.Column, toColumn) {
		return nil, fmt.Errorf("invalid transition from %s to %s", card.Column, toColumn)
	}

	now := time.Now()
	card.Column = toColumn
	card.UpdatedAt = now

	if toColumn.IsDone() {
		card.CompletedAt = &now
	} else {
		card.CompletedAt = nil
	}

	b.UpdatedAt = now

	// Return copy with computed fields
	cardCopy := *card
	cardCopy.BlockedBy = b.getBlockingCards(card)

	return &cardCopy, nil
}

// DeleteCard removes a card from the board.
func (b *Board) DeleteCard(ctx context.Context, id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.Cards[id]; !ok {
		return fmt.Errorf("card %q not found", id)
	}

	// Remove dependencies pointing to this card
	for _, card := range b.Cards {
		for i, dep := range card.DependsOn {
			if dep == id {
				card.DependsOn = append(card.DependsOn[:i], card.DependsOn[i+1:]...)
				break
			}
		}
	}

	delete(b.Cards, id)
	b.UpdatedAt = time.Now()

	return nil
}

// AddDependency adds a dependency between cards.
func (b *Board) AddDependency(ctx context.Context, cardID, dependsOnID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	card, ok := b.Cards[cardID]
	if !ok {
		return fmt.Errorf("card %q not found", cardID)
	}

	if _, ok := b.Cards[dependsOnID]; !ok {
		return fmt.Errorf("dependency card %q not found", dependsOnID)
	}

	if cardID == dependsOnID {
		return fmt.Errorf("card cannot depend on itself")
	}

	// Check for circular dependency
	if b.wouldCreateCycle(cardID, dependsOnID) {
		return fmt.Errorf("adding dependency would create a cycle")
	}

	card.AddDependency(dependsOnID)
	card.UpdatedAt = time.Now()
	b.UpdatedAt = time.Now()

	return nil
}

// RemoveDependency removes a dependency between cards.
func (b *Board) RemoveDependency(ctx context.Context, cardID, dependsOnID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	card, ok := b.Cards[cardID]
	if !ok {
		return fmt.Errorf("card %q not found", cardID)
	}

	card.RemoveDependency(dependsOnID)
	card.UpdatedAt = time.Now()
	b.UpdatedAt = time.Now()

	return nil
}

// ListCards returns cards matching the filter.
func (b *Board) ListCards(ctx context.Context, filter CardFilter) ([]*Card, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var cards []*Card
	for _, card := range b.Cards {
		if filter.matches(card) {
			cardCopy := *card
			cardCopy.BlockedBy = b.getBlockingCards(card)
			cards = append(cards, &cardCopy)
		}
	}

	// Sort by priority (highest first), then by creation time
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].Priority != cards[j].Priority {
			return cards[i].Priority > cards[j].Priority
		}
		return cards[i].CreatedAt.Before(cards[j].CreatedAt)
	})

	return cards, nil
}

// CardFilter specifies criteria for listing cards.
type CardFilter struct {
	// Columns filters by column (any of these).
	Columns []ColumnType

	// Priorities filters by priority (any of these).
	Priorities []Priority

	// Labels filters by labels (all must match).
	Labels []string

	// Assignee filters by assignee.
	Assignee string

	// IncludeBlocked includes blocked cards (default true).
	IncludeBlocked *bool
}

// matches checks if a card matches the filter.
func (f *CardFilter) matches(card *Card) bool {
	if len(f.Columns) > 0 {
		found := false
		for _, col := range f.Columns {
			if card.Column == col {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(f.Priorities) > 0 {
		found := false
		for _, p := range f.Priorities {
			if card.Priority == p {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(f.Labels) > 0 {
		for _, label := range f.Labels {
			if !card.HasLabel(label) {
				return false
			}
		}
	}

	if f.Assignee != "" && card.Assignee != f.Assignee {
		return false
	}

	return true
}

// getBlockingCards returns the IDs of cards that block the given card.
func (b *Board) getBlockingCards(card *Card) []string {
	var blockers []string
	for _, depID := range card.DependsOn {
		if dep, ok := b.Cards[depID]; ok {
			if !dep.Column.IsDone() {
				blockers = append(blockers, depID)
			}
		}
	}
	return blockers
}

// wouldCreateCycle checks if adding a dependency would create a cycle.
func (b *Board) wouldCreateCycle(fromID, toID string) bool {
	visited := make(map[string]bool)
	var dfs func(id string) bool
	dfs = func(id string) bool {
		if id == fromID {
			return true
		}
		if visited[id] {
			return false
		}
		visited[id] = true
		if card, ok := b.Cards[id]; ok {
			for _, depID := range card.DependsOn {
				if dfs(depID) {
					return true
				}
			}
		}
		return false
	}
	return dfs(toID)
}

// Stats returns statistics about the board.
func (b *Board) Stats() BoardStats {
	b.mu.RLock()
	defer b.mu.RUnlock()

	stats := BoardStats{
		TotalCards: len(b.Cards),
		ByColumn:   make(map[ColumnType]int),
		ByPriority: make(map[Priority]int),
	}

	for _, card := range b.Cards {
		stats.ByColumn[card.Column]++
		stats.ByPriority[card.Priority]++
		if len(b.getBlockingCards(card)) > 0 {
			stats.BlockedCards++
		}
	}

	return stats
}

// BoardStats contains board statistics.
type BoardStats struct {
	TotalCards   int
	ByColumn     map[ColumnType]int
	ByPriority   map[Priority]int
	BlockedCards int
}
