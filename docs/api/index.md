# API Reference

This page provides an overview of the main types and functions in OmniWorkboard.

## Core Types

### Board

```go
type Board struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description,omitempty"`
    Columns     []ColumnType      `json:"columns"`
    Cards       map[string]*Card  `json:"cards"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}
```

### Card

```go
type Card struct {
    ID          string            `json:"id"`
    Title       string            `json:"title"`
    Description string            `json:"description"`
    Column      ColumnType        `json:"column"`
    Priority    Priority          `json:"priority"`
    DependsOn   []string          `json:"depends_on"`
    BlockedBy   []string          `json:"blocked_by"`
    Labels      []string          `json:"labels"`
    Metadata    map[string]any    `json:"metadata"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
    CompletedAt *time.Time        `json:"completed_at,omitempty"`
}
```

### ColumnType

```go
type ColumnType string

const (
    ColumnBacklog    ColumnType = "backlog"
    ColumnTodo       ColumnType = "todo"
    ColumnInProgress ColumnType = "in_progress"
    ColumnReview     ColumnType = "review"
    ColumnDone       ColumnType = "done"
)
```

### Priority

```go
type Priority int

const (
    PriorityLow      Priority = 0
    PriorityNormal   Priority = 1
    PriorityHigh     Priority = 2
    PriorityCritical Priority = 3
)
```

## Board Functions

### NewBoard

```go
func NewBoard(cfg BoardConfig) *Board
```

Creates a new workboard with the given configuration.

### Board Methods

```go
func (b *Board) CreateCard(ctx context.Context, title, description string, priority Priority) (*Card, error)
func (b *Board) GetCard(ctx context.Context, id string) (*Card, error)
func (b *Board) UpdateCard(ctx context.Context, id string, opts UpdateCardOptions) error
func (b *Board) DeleteCard(ctx context.Context, id string) error
func (b *Board) MoveCard(ctx context.Context, id string, column ColumnType) error
func (b *Board) ListCards(ctx context.Context, opts ListOptions) []*Card
func (b *Board) AddDependency(ctx context.Context, cardID, dependsOnID string) error
func (b *Board) RemoveDependency(ctx context.Context, cardID, dependsOnID string) error
func (b *Board) AddLabel(ctx context.Context, cardID, label string) error
func (b *Board) RemoveLabel(ctx context.Context, cardID, label string) error
```

## Configuration Types

### BoardConfig

```go
type BoardConfig struct {
    Name        string
    Description string
    Columns     []ColumnType
}
```

### UpdateCardOptions

```go
type UpdateCardOptions struct {
    Title       string
    Description string
    Priority    Priority
    Labels      []string
}
```

### ListOptions

```go
type ListOptions struct {
    Column      ColumnType
    Priority    Priority
    Label       string
    BlockedOnly bool
}
```

## Skill Integration

### Skill

```go
type Skill struct {
    // implements compiled.Skill
}

func NewSkill() *Skill
func (s *Skill) Name() string
func (s *Skill) Description() string
func (s *Skill) Tools() []skill.Tool
func (s *Skill) Init(ctx context.Context) error
func (s *Skill) Close() error
func (s *Skill) SetStorage(store kvs.Store)
```

### WorkboardTools

```go
type WorkboardTools struct {
    board *Board
}

func NewWorkboardTools(board *Board) *WorkboardTools
func (t *WorkboardTools) CreateCard(ctx context.Context, params map[string]any) (any, error)
func (t *WorkboardTools) MoveCard(ctx context.Context, params map[string]any) (any, error)
func (t *WorkboardTools) ListCards(ctx context.Context, params map[string]any) (any, error)
func (t *WorkboardTools) UpdateCard(ctx context.Context, params map[string]any) (any, error)
func (t *WorkboardTools) AddDependency(ctx context.Context, params map[string]any) (any, error)
func (t *WorkboardTools) GetBoard(ctx context.Context, params map[string]any) (any, error)
```

## Errors

```go
var (
    ErrCardNotFound    = errors.New("card not found")
    ErrCycleDetected   = errors.New("dependency would create cycle")
    ErrSelfDependency  = errors.New("card cannot depend on itself")
    ErrInvalidColumn   = errors.New("invalid column")
    ErrInvalidPriority = errors.New("invalid priority")
)
```

## Helper Functions

### DefaultColumns

```go
func DefaultColumns() []ColumnType
```

Returns the default column order: Backlog, Todo, In Progress, Review, Done.

### PriorityFromString

```go
func PriorityFromString(s string) (Priority, error)
```

Converts a string to Priority.

### ColumnFromString

```go
func ColumnFromString(s string) (ColumnType, error)
```

Converts a string to ColumnType.

## Full Documentation

For complete API documentation, see [pkg.go.dev/github.com/plexusone/omniworkboard](https://pkg.go.dev/github.com/plexusone/omniworkboard).
