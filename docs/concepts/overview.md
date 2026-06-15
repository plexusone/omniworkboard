# Concepts Overview

OmniWorkboard implements a kanban-style workboard with boards, columns, cards, and dependencies.

## Architecture

```
Board
├── Columns: [Backlog, Todo, In Progress, Review, Done]
└── Cards
    ├── Card 1
    │   ├── Title, Description
    │   ├── Column: In Progress
    │   ├── Priority: High
    │   ├── Labels: [bug, urgent]
    │   └── DependsOn: []
    └── Card 2
        ├── Title, Description
        ├── Column: Todo
        ├── Priority: Normal
        ├── Labels: [feature]
        └── DependsOn: [Card 1]  ← Blocked by Card 1
```

## Core Concepts

### Board

A board is a container for cards organized into columns. Each board has:

- **Name** - Identifier (e.g., "Sprint 42")
- **Description** - Optional details
- **Columns** - Ordered list of workflow stages
- **Cards** - Task items on the board

### Columns

Columns represent workflow stages:

| Column | Purpose |
|--------|---------|
| `Backlog` | Future work, not yet planned |
| `Todo` | Planned work, ready to start |
| `In Progress` | Currently being worked on |
| `Review` | Completed, awaiting review |
| `Done` | Finished work |

### Cards

Cards represent individual tasks:

| Field | Description |
|-------|-------------|
| `Title` | Short task name |
| `Description` | Detailed description |
| `Column` | Current workflow stage |
| `Priority` | Low, Normal, High, Critical |
| `Labels` | Categorization tags |
| `DependsOn` | Cards that must complete first |
| `BlockedBy` | Computed list of blocking cards |

### Dependencies

Dependencies model "must complete before" relationships:

```
Card A ──depends on──> Card B
        (Card A blocked until Card B is Done)
```

## Workflow

### Card Lifecycle

```
Created
    │
    ▼
┌─────────┐    ┌──────┐    ┌─────────────┐    ┌────────┐    ┌──────┐
│ Backlog │ → │ Todo │ → │ In Progress │ → │ Review │ → │ Done │
└─────────┘    └──────┘    └─────────────┘    └────────┘    └──────┘
```

### Blocking Logic

1. Card B depends on Card A
2. While Card A is not Done, Card B shows as "blocked"
3. When Card A moves to Done, Card B becomes unblocked
4. Card B can then proceed through the workflow

## Priority System

| Priority | When to Use |
|----------|-------------|
| `Critical` | Production issues, urgent bugs |
| `High` | Important features, deadlines |
| `Normal` | Regular work items |
| `Low` | Nice-to-haves, cleanup |

## Labels

Labels categorize cards:

```
bug, feature, enhancement, documentation
frontend, backend, infrastructure
p0, p1, p2, p3
sprint-42, q3-release
```

## State Management

### Thread Safety

All board operations are thread-safe:

```go
// Safe for concurrent access
go board.CreateCard(ctx, ...)
go board.MoveCard(ctx, ...)
go board.ListCards(ctx, ...)
```

### Persistence

Boards can be persisted using kvs.Store:

```go
skill := omniworkboard.NewSkill()
skill.SetStorage(store)  // Auto-saves state
```

## Integration Points

### AI Agents

Via compiled skill:

```go
agent.WithCompiledSkill(omniworkboard.NewSkill())
```

### Direct API

Via Go package:

```go
board := omniworkboard.NewBoard(config)
card, _ := board.CreateCard(ctx, ...)
```

### CLI

Via omniagent CLI:

```bash
omniagent chat
> Create a card for "Fix bug"
```

## See Also

- [Boards](boards.md) - Board operations
- [Cards](cards.md) - Card operations
- [Dependencies](dependencies.md) - Dependency management
