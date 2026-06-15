package omniworkboard

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/plexusone/omniskill/skill"
	"github.com/plexusone/omnistorage-core/kvs"
)

const (
	// SkillName is the name of the workboard skill.
	SkillName = "workboard"

	// StorageKeyBoard is the key for persisting the board.
	StorageKeyBoard = "workboard:board"
)

// Skill implements compiled.Skill for the workboard.
type Skill struct {
	board   *Board
	tools   *WorkboardTools
	storage kvs.Store
}

// NewSkill creates a new workboard skill.
func NewSkill() *Skill {
	return &Skill{}
}

// NewSkillWithBoard creates a new workboard skill with an existing board.
func NewSkillWithBoard(board *Board) *Skill {
	s := &Skill{board: board}
	s.tools = NewWorkboardTools(board)
	return s
}

// Name implements compiled.Skill.
func (s *Skill) Name() string {
	return SkillName
}

// Description implements compiled.Skill.
func (s *Skill) Description() string {
	return "Project management workboard for tracking tasks with columns, priorities, and dependencies"
}

// SetStorage implements compiled.StorageAware.
func (s *Skill) SetStorage(store kvs.Store) {
	s.storage = store
}

// Init implements compiled.Skill.
func (s *Skill) Init(ctx context.Context) error {
	// Try to load existing board from storage
	if s.storage != nil && s.board == nil {
		if err := s.loadBoard(ctx); err != nil {
			// No existing board, create a new one
			s.board = NewBoard(BoardConfig{Name: "Default Board"})
		}
	}

	// Create board if still nil
	if s.board == nil {
		s.board = NewBoard(BoardConfig{Name: "Default Board"})
	}

	s.tools = NewWorkboardTools(s.board)
	return nil
}

// Close implements compiled.Skill.
func (s *Skill) Close() error {
	// Persist the board before closing
	if s.storage != nil && s.board != nil {
		ctx := context.Background()
		return s.saveBoard(ctx)
	}
	return nil
}

// Tools implements compiled.Skill.
func (s *Skill) Tools() []skill.Tool {
	if s.tools == nil {
		return nil
	}

	toolDefs := s.tools.Tools()
	result := make([]skill.Tool, len(toolDefs))

	for i, t := range toolDefs {
		result[i] = &workboardTool{
			name:        t.Name,
			description: t.Description,
			params:      convertInputSchema(t.InputSchema),
			handler:     s.tools.handlers[t.Name],
		}
	}

	return result
}

// Board returns the current board.
func (s *Skill) Board() *Board {
	return s.board
}

// loadBoard loads the board from storage.
func (s *Skill) loadBoard(ctx context.Context) error {
	data, err := s.storage.Get(ctx, StorageKeyBoard)
	if err != nil {
		return err
	}

	var board Board
	if err := json.Unmarshal(data, &board); err != nil {
		return fmt.Errorf("unmarshal board: %w", err)
	}

	s.board = &board
	return nil
}

// saveBoard persists the board to storage.
func (s *Skill) saveBoard(ctx context.Context) error {
	data, err := json.Marshal(s.board)
	if err != nil {
		return fmt.Errorf("marshal board: %w", err)
	}

	return s.storage.Set(ctx, StorageKeyBoard, data, 0)
}

// workboardTool wraps a workboard tool handler as a skill.Tool.
type workboardTool struct {
	name        string
	description string
	params      map[string]skill.Parameter
	handler     ToolHandler
}

func (t *workboardTool) Name() string {
	return t.name
}

func (t *workboardTool) Description() string {
	return t.description
}

func (t *workboardTool) Parameters() map[string]skill.Parameter {
	return t.params
}

func (t *workboardTool) Call(ctx context.Context, params map[string]any) (any, error) {
	// Convert params to JSON for the handler
	data, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal params: %w", err)
	}

	return t.handler(ctx, data)
}

// convertInputSchema converts a JSON Schema map to skill.Parameter map.
func convertInputSchema(schema map[string]any) map[string]skill.Parameter {
	result := make(map[string]skill.Parameter)

	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		return result
	}

	required := make(map[string]bool)
	if reqList, ok := schema["required"].([]string); ok {
		for _, r := range reqList {
			required[r] = true
		}
	}

	for name, propRaw := range properties {
		prop, ok := propRaw.(map[string]any)
		if !ok {
			continue
		}

		param := skill.Parameter{
			Required: required[name],
		}

		if typ, ok := prop["type"].(string); ok {
			param.Type = typ
		}
		if desc, ok := prop["description"].(string); ok {
			param.Description = desc
		}
		if enum, ok := prop["enum"].([]string); ok {
			param.Enum = make([]any, len(enum))
			for i, e := range enum {
				param.Enum[i] = e
			}
		}

		result[name] = param
	}

	return result
}
