# Getting Started

This guide walks through setting up OmniWorkboard for project management.

## Installation

```bash
go get github.com/plexusone/omniworkboard@latest
```

## Basic Usage

### Creating a Board

```go
import "github.com/plexusone/omniworkboard"

// Create with default columns
board := omniworkboard.NewBoard(omniworkboard.BoardConfig{
    Name:        "My Project",
    Description: "Project description",
})

// Create with custom columns
board := omniworkboard.NewBoard(omniworkboard.BoardConfig{
    Name: "Simple Board",
    Columns: []omniworkboard.ColumnType{
        omniworkboard.ColumnTodo,
        omniworkboard.ColumnInProgress,
        omniworkboard.ColumnDone,
    },
})
```

### Creating Cards

```go
ctx := context.Background()

// Create a high-priority card
card, err := board.CreateCard(ctx,
    "Implement feature X",           // Title
    "Detailed description here",      // Description
    omniworkboard.PriorityHigh,       // Priority
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created card: %s\n", card.ID)
```

### Moving Cards

```go
// Move to In Progress
err := board.MoveCard(ctx, card.ID, omniworkboard.ColumnInProgress)
if err != nil {
    log.Fatal(err)
}

// Move to Done
err = board.MoveCard(ctx, card.ID, omniworkboard.ColumnDone)
```

### Listing Cards

```go
// List all cards
cards := board.ListCards(ctx, omniworkboard.ListOptions{})

// List by column
inProgress := board.ListCards(ctx, omniworkboard.ListOptions{
    Column: omniworkboard.ColumnInProgress,
})

// List by priority
highPriority := board.ListCards(ctx, omniworkboard.ListOptions{
    Priority: omniworkboard.PriorityHigh,
})

// List by label
bugFixes := board.ListCards(ctx, omniworkboard.ListOptions{
    Label: "bug",
})
```

### Managing Dependencies

```go
// Create cards
taskA, _ := board.CreateCard(ctx, "Task A", "", omniworkboard.PriorityNormal)
taskB, _ := board.CreateCard(ctx, "Task B", "", omniworkboard.PriorityNormal)

// Task B depends on Task A (must complete A before B)
err := board.AddDependency(ctx, taskB.ID, taskA.ID)

// Check if card is blocked
card, _ := board.GetCard(ctx, taskB.ID)
if len(card.BlockedBy) > 0 {
    fmt.Println("Task B is blocked by:", card.BlockedBy)
}
```

## Workflow Example

A typical workflow:

```go
ctx := context.Background()

// 1. Create board
board := omniworkboard.NewBoard(omniworkboard.BoardConfig{
    Name: "Sprint 42",
})

// 2. Add cards to backlog
login, _ := board.CreateCard(ctx, "Implement OAuth login", "", omniworkboard.PriorityHigh)
tests, _ := board.CreateCard(ctx, "Write auth tests", "", omniworkboard.PriorityNormal)
docs, _ := board.CreateCard(ctx, "Update API docs", "", omniworkboard.PriorityLow)

// 3. Set dependencies
board.AddDependency(ctx, tests.ID, login.ID)  // tests need login first
board.AddDependency(ctx, docs.ID, tests.ID)   // docs need tests first

// 4. Start sprint - move to todo
board.MoveCard(ctx, login.ID, omniworkboard.ColumnTodo)

// 5. Work on login
board.MoveCard(ctx, login.ID, omniworkboard.ColumnInProgress)
// ... do work ...
board.MoveCard(ctx, login.ID, omniworkboard.ColumnReview)
// ... review passes ...
board.MoveCard(ctx, login.ID, omniworkboard.ColumnDone)

// 6. tests is now unblocked, move to todo
board.MoveCard(ctx, tests.ID, omniworkboard.ColumnTodo)

// 7. Check remaining work
remaining := board.ListCards(ctx, omniworkboard.ListOptions{})
for _, c := range remaining {
    if c.Column != omniworkboard.ColumnDone {
        fmt.Printf("Remaining: %s [%s]\n", c.Title, c.Column)
    }
}
```

## Persistence

### Using kvs.Store

```go
import (
    "github.com/plexusone/omnistorage-core/kvs"
    "github.com/plexusone/omnistorage-core/sqlite"
)

// Create storage
backend, _ := sqlite.Open("workboard.db")
store := kvs.New(backend)

// Create skill with storage
skill := omniworkboard.NewSkill()
skill.SetStorage(store)
```

### Manual Save/Load

```go
// Save board state
data, _ := json.Marshal(board)
store.Set(ctx, "board:sprint-42", data)

// Load board state
data, _ := store.Get(ctx, "board:sprint-42")
var board omniworkboard.Board
json.Unmarshal(data, &board)
```

## Agent Integration

See [Agent Integration](agent-integration.md) for details on using OmniWorkboard with AI agents.

## Next Steps

- [Concepts Overview](concepts/overview.md)
- [Boards](concepts/boards.md)
- [Cards](concepts/cards.md)
- [Dependencies](concepts/dependencies.md)
