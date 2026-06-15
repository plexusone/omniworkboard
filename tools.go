package omniworkboard

import (
	"context"
	"encoding/json"
	"fmt"
)

// Tool represents a tool that can be invoked by an agent.
type Tool struct {
	// Name is the tool name.
	Name string `json:"name"`

	// Description describes what the tool does.
	Description string `json:"description"`

	// InputSchema is the JSON schema for the tool's input.
	InputSchema map[string]any `json:"input_schema"`
}

// ToolHandler handles tool invocations.
type ToolHandler func(ctx context.Context, input json.RawMessage) (any, error)

// WorkboardTools provides tools for interacting with a workboard.
type WorkboardTools struct {
	board    *Board
	handlers map[string]ToolHandler
}

// NewWorkboardTools creates tools for the given board.
func NewWorkboardTools(board *Board) *WorkboardTools {
	wt := &WorkboardTools{
		board:    board,
		handlers: make(map[string]ToolHandler),
	}

	wt.handlers["create_card"] = wt.handleCreateCard
	wt.handlers["move_card"] = wt.handleMoveCard
	wt.handlers["list_cards"] = wt.handleListCards
	wt.handlers["update_card"] = wt.handleUpdateCard
	wt.handlers["add_dependency"] = wt.handleAddDependency
	wt.handlers["get_card"] = wt.handleGetCard
	wt.handlers["delete_card"] = wt.handleDeleteCard
	wt.handlers["board_stats"] = wt.handleBoardStats

	return wt
}

// Tools returns the list of available tools.
func (wt *WorkboardTools) Tools() []Tool {
	return []Tool{
		{
			Name:        "create_card",
			Description: "Create a new task card on the workboard",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title": map[string]any{
						"type":        "string",
						"description": "The card title",
					},
					"description": map[string]any{
						"type":        "string",
						"description": "Optional detailed description",
					},
					"priority": map[string]any{
						"type":        "string",
						"enum":        []string{"low", "normal", "high", "critical"},
						"description": "Card priority (default: normal)",
					},
					"labels": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "Optional labels for categorization",
					},
				},
				"required": []string{"title"},
			},
		},
		{
			Name:        "move_card",
			Description: "Move a card to a different column",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"card_id": map[string]any{
						"type":        "string",
						"description": "The card ID to move",
					},
					"column": map[string]any{
						"type":        "string",
						"enum":        []string{"backlog", "todo", "in_progress", "review", "done"},
						"description": "Target column",
					},
				},
				"required": []string{"card_id", "column"},
			},
		},
		{
			Name:        "list_cards",
			Description: "List cards with optional filters",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"columns": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "Filter by columns",
					},
					"priorities": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "Filter by priorities",
					},
					"labels": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "Filter by labels (all must match)",
					},
					"assignee": map[string]any{
						"type":        "string",
						"description": "Filter by assignee",
					},
				},
			},
		},
		{
			Name:        "update_card",
			Description: "Update card details",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"card_id": map[string]any{
						"type":        "string",
						"description": "The card ID to update",
					},
					"title": map[string]any{
						"type":        "string",
						"description": "New title",
					},
					"description": map[string]any{
						"type":        "string",
						"description": "New description",
					},
					"priority": map[string]any{
						"type":        "string",
						"enum":        []string{"low", "normal", "high", "critical"},
						"description": "New priority",
					},
					"assignee": map[string]any{
						"type":        "string",
						"description": "New assignee",
					},
					"labels": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "New labels (replaces existing)",
					},
				},
				"required": []string{"card_id"},
			},
		},
		{
			Name:        "add_dependency",
			Description: "Add a dependency between cards",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"card_id": map[string]any{
						"type":        "string",
						"description": "The card that depends on another",
					},
					"depends_on": map[string]any{
						"type":        "string",
						"description": "The card that must be completed first",
					},
				},
				"required": []string{"card_id", "depends_on"},
			},
		},
		{
			Name:        "get_card",
			Description: "Get details of a specific card",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"card_id": map[string]any{
						"type":        "string",
						"description": "The card ID",
					},
				},
				"required": []string{"card_id"},
			},
		},
		{
			Name:        "delete_card",
			Description: "Delete a card from the board",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"card_id": map[string]any{
						"type":        "string",
						"description": "The card ID to delete",
					},
				},
				"required": []string{"card_id"},
			},
		},
		{
			Name:        "board_stats",
			Description: "Get workboard statistics",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
	}
}

// Invoke calls a tool with the given input.
func (wt *WorkboardTools) Invoke(ctx context.Context, name string, input json.RawMessage) (any, error) {
	handler, ok := wt.handlers[name]
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	return handler(ctx, input)
}

// createCardInput is the input for create_card.
type createCardInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    string   `json:"priority"`
	Labels      []string `json:"labels"`
}

func (wt *WorkboardTools) handleCreateCard(ctx context.Context, input json.RawMessage) (any, error) {
	var in createCardInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}

	priority := ParsePriority(in.Priority)
	card, err := wt.board.CreateCard(ctx, in.Title, in.Description, priority)
	if err != nil {
		return nil, err
	}

	for _, label := range in.Labels {
		card.AddLabel(label)
	}

	return map[string]any{
		"id":       card.ID,
		"title":    card.Title,
		"column":   card.Column,
		"priority": card.Priority.String(),
	}, nil
}

// moveCardInput is the input for move_card.
type moveCardInput struct {
	CardID string `json:"card_id"`
	Column string `json:"column"`
}

func (wt *WorkboardTools) handleMoveCard(ctx context.Context, input json.RawMessage) (any, error) {
	var in moveCardInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}

	column := ParseColumnType(in.Column)
	card, err := wt.board.MoveCard(ctx, in.CardID, column)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"id":     card.ID,
		"column": card.Column.String(),
	}, nil
}

// listCardsInput is the input for list_cards.
type listCardsInput struct {
	Columns    []string `json:"columns"`
	Priorities []string `json:"priorities"`
	Labels     []string `json:"labels"`
	Assignee   string   `json:"assignee"`
}

func (wt *WorkboardTools) handleListCards(ctx context.Context, input json.RawMessage) (any, error) {
	var in listCardsInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}

	filter := CardFilter{
		Labels:   in.Labels,
		Assignee: in.Assignee,
	}

	for _, col := range in.Columns {
		filter.Columns = append(filter.Columns, ParseColumnType(col))
	}
	for _, p := range in.Priorities {
		filter.Priorities = append(filter.Priorities, ParsePriority(p))
	}

	cards, err := wt.board.ListCards(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]any, len(cards))
	for i, card := range cards {
		result[i] = map[string]any{
			"id":         card.ID,
			"title":      card.Title,
			"column":     card.Column.String(),
			"priority":   card.Priority.String(),
			"blocked_by": card.BlockedBy,
		}
	}

	return result, nil
}

// updateCardInput is the input for update_card.
type updateCardInput struct {
	CardID      string   `json:"card_id"`
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Priority    *string  `json:"priority,omitempty"`
	Assignee    *string  `json:"assignee,omitempty"`
	Labels      []string `json:"labels,omitempty"`
}

func (wt *WorkboardTools) handleUpdateCard(ctx context.Context, input json.RawMessage) (any, error) {
	var in updateCardInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}

	updates := CardUpdate{
		Title:       in.Title,
		Description: in.Description,
		Assignee:    in.Assignee,
		Labels:      in.Labels,
	}

	if in.Priority != nil {
		p := ParsePriority(*in.Priority)
		updates.Priority = &p
	}

	card, err := wt.board.UpdateCard(ctx, in.CardID, updates)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"id":       card.ID,
		"title":    card.Title,
		"priority": card.Priority.String(),
	}, nil
}

// addDependencyInput is the input for add_dependency.
type addDependencyInput struct {
	CardID    string `json:"card_id"`
	DependsOn string `json:"depends_on"`
}

func (wt *WorkboardTools) handleAddDependency(ctx context.Context, input json.RawMessage) (any, error) {
	var in addDependencyInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}

	if err := wt.board.AddDependency(ctx, in.CardID, in.DependsOn); err != nil {
		return nil, err
	}

	return map[string]any{
		"success": true,
	}, nil
}

// getCardInput is the input for get_card.
type getCardInput struct {
	CardID string `json:"card_id"`
}

func (wt *WorkboardTools) handleGetCard(ctx context.Context, input json.RawMessage) (any, error) {
	var in getCardInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}

	card, err := wt.board.GetCard(ctx, in.CardID)
	if err != nil {
		return nil, err
	}

	return card, nil
}

// deleteCardInput is the input for delete_card.
type deleteCardInput struct {
	CardID string `json:"card_id"`
}

func (wt *WorkboardTools) handleDeleteCard(ctx context.Context, input json.RawMessage) (any, error) {
	var in deleteCardInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}

	if err := wt.board.DeleteCard(ctx, in.CardID); err != nil {
		return nil, err
	}

	return map[string]any{
		"success": true,
	}, nil
}

func (wt *WorkboardTools) handleBoardStats(ctx context.Context, _ json.RawMessage) (any, error) {
	stats := wt.board.Stats()
	return stats, nil
}
