# Boards

A board is the top-level container that organizes cards into columns.

## Creating Boards

### Default Configuration

```go
board := omniworkboard.NewBoard(omniworkboard.BoardConfig{
    Name: "My Project",
})
```

Default columns: Backlog → Todo → In Progress → Review → Done

### Custom Columns

```go
board := omniworkboard.NewBoard(omniworkboard.BoardConfig{
    Name: "Simple Board",
    Columns: []omniworkboard.ColumnType{
        omniworkboard.ColumnTodo,
        omniworkboard.ColumnInProgress,
        omniworkboard.ColumnDone,
    },
})
```

## Board Structure

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

## Column Types

```go
const (
    ColumnBacklog    ColumnType = "backlog"
    ColumnTodo       ColumnType = "todo"
    ColumnInProgress ColumnType = "in_progress"
    ColumnReview     ColumnType = "review"
    ColumnDone       ColumnType = "done"
)
```

### Default Columns

```go
func DefaultColumns() []ColumnType {
    return []ColumnType{
        ColumnBacklog,
        ColumnTodo,
        ColumnInProgress,
        ColumnReview,
        ColumnDone,
    }
}
```

## Board Operations

### Get Board Info

```go
fmt.Printf("Board: %s\n", board.Name)
fmt.Printf("Columns: %v\n", board.Columns)
fmt.Printf("Cards: %d\n", len(board.Cards))
```

### Get Cards by Column

```go
cards := board.ListCards(ctx, omniworkboard.ListOptions{
    Column: omniworkboard.ColumnInProgress,
})
```

### Board Statistics

```go
stats := board.Stats(ctx)
fmt.Printf("Total cards: %d\n", stats.Total)
fmt.Printf("In Progress: %d\n", stats.ByColumn[omniworkboard.ColumnInProgress])
fmt.Printf("Blocked: %d\n", stats.Blocked)
```

## Column Workflow

### Valid Transitions

Cards can move to any column, but typical workflow is:

```
Backlog → Todo → In Progress → Review → Done
```

### Skip Columns

Cards can skip columns if needed:

```go
// Direct to In Progress
board.MoveCard(ctx, cardID, omniworkboard.ColumnInProgress)

// Direct to Done (e.g., cancelled)
board.MoveCard(ctx, cardID, omniworkboard.ColumnDone)
```

### Move Backwards

Cards can move backwards:

```go
// Failed review, back to In Progress
board.MoveCard(ctx, cardID, omniworkboard.ColumnInProgress)
```

## Multiple Boards

For multiple projects, create separate boards:

```go
sprint42 := omniworkboard.NewBoard(omniworkboard.BoardConfig{
    Name: "Sprint 42",
})

sprint43 := omniworkboard.NewBoard(omniworkboard.BoardConfig{
    Name: "Sprint 43",
})
```

## Serialization

### JSON Export

```go
data, err := json.Marshal(board)
if err != nil {
    log.Fatal(err)
}
// Store or transmit data
```

### JSON Import

```go
var board omniworkboard.Board
err := json.Unmarshal(data, &board)
if err != nil {
    log.Fatal(err)
}
```

## Best Practices

### Board Naming

Use clear, descriptive names:

```go
// Good
"Sprint 42 - Q3 Features"
"Bug Triage - June 2026"
"Project Alpha - Phase 1"

// Less clear
"Board 1"
"My Tasks"
```

### Column Selection

Match columns to your workflow:

| Workflow | Columns |
|----------|---------|
| Simple | Todo, In Progress, Done |
| Agile | Backlog, Todo, In Progress, Review, Done |
| Support | New, Triaging, Investigating, Resolved |
| Content | Draft, Review, Editing, Published |

### Board Size

Keep boards focused:

- One board per project/sprint
- Archive completed boards
- Split large boards if > 50 active cards

## See Also

- [Cards](cards.md) - Card operations
- [Dependencies](dependencies.md) - Task dependencies
- [Overview](overview.md) - Architecture overview
