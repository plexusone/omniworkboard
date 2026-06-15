# Dependencies

Dependencies model "must complete before" relationships between cards.

## Overview

When Card B depends on Card A:

- Card B is "blocked by" Card A
- Card A is "blocking" Card B
- Card B cannot be considered complete until Card A is Done

```
Card A (Login)
    │
    │ blocks
    ▼
Card B (Tests) ← blocked until Card A is Done
    │
    │ blocks
    ▼
Card C (Deploy) ← blocked until Card B is Done
```

## Adding Dependencies

```go
// Card B depends on Card A
err := board.AddDependency(ctx, cardB.ID, cardA.ID)

// Multiple dependencies
err = board.AddDependency(ctx, cardC.ID, cardA.ID)
err = board.AddDependency(ctx, cardC.ID, cardB.ID)
```

## Removing Dependencies

```go
err := board.RemoveDependency(ctx, cardB.ID, cardA.ID)
```

## Checking Dependencies

### DependsOn vs BlockedBy

```go
card, _ := board.GetCard(ctx, cardID)

// DependsOn - cards this card waits for
fmt.Printf("Depends on: %v\n", card.DependsOn)

// BlockedBy - cards that are blocking (not Done)
fmt.Printf("Blocked by: %v\n", card.BlockedBy)
```

**Key difference:**

- `DependsOn` - All declared dependencies
- `BlockedBy` - Dependencies that are not yet Done

### Example

```go
// Card A is Done
// Card B (depends on A, C) - A is Done, C is not
card, _ := board.GetCard(ctx, cardB.ID)

card.DependsOn  // ["cardA", "cardC"]
card.BlockedBy  // ["cardC"] - only C is blocking
```

## Blocked Cards

### List Blocked Cards

```go
blocked := board.ListCards(ctx, omniworkboard.ListOptions{
    BlockedOnly: true,
})

for _, c := range blocked {
    fmt.Printf("%s blocked by: %v\n", c.Title, c.BlockedBy)
}
```

### Check if Blocked

```go
card, _ := board.GetCard(ctx, cardID)
if len(card.BlockedBy) > 0 {
    fmt.Printf("Card is blocked by %d cards\n", len(card.BlockedBy))
}
```

## Dependency Chains

Dependencies can form chains:

```go
// A → B → C → D
board.AddDependency(ctx, cardB.ID, cardA.ID)
board.AddDependency(ctx, cardC.ID, cardB.ID)
board.AddDependency(ctx, cardD.ID, cardC.ID)
```

When A completes, B unblocks. When B completes, C unblocks. And so on.

## Cycle Detection

Circular dependencies are prevented:

```go
// This would create A → B → A cycle
board.AddDependency(ctx, cardA.ID, cardB.ID)
board.AddDependency(ctx, cardB.ID, cardA.ID) // Error: cycle detected
```

## Dependency Patterns

### Sequential Tasks

```
Design → Implement → Test → Deploy
```

```go
board.AddDependency(ctx, implement.ID, design.ID)
board.AddDependency(ctx, test.ID, implement.ID)
board.AddDependency(ctx, deploy.ID, test.ID)
```

### Parallel with Join

```
     ┌─ Frontend ─┐
API ─┤            ├─ Integration
     └─ Backend ──┘
```

```go
board.AddDependency(ctx, frontend.ID, api.ID)
board.AddDependency(ctx, backend.ID, api.ID)
board.AddDependency(ctx, integration.ID, frontend.ID)
board.AddDependency(ctx, integration.ID, backend.ID)
```

### Multiple Prerequisites

```
     ┌─ Requirement A ─┐
Task ┼─ Requirement B ─┼
     └─ Requirement C ─┘
```

```go
board.AddDependency(ctx, task.ID, reqA.ID)
board.AddDependency(ctx, task.ID, reqB.ID)
board.AddDependency(ctx, task.ID, reqC.ID)
```

## Unblocking

When a blocking card moves to Done, dependents are automatically unblocked:

```go
// B depends on A
board.AddDependency(ctx, cardB.ID, cardA.ID)

// B is blocked
card, _ := board.GetCard(ctx, cardB.ID)
fmt.Println(card.BlockedBy) // ["cardA"]

// Complete A
board.MoveCard(ctx, cardA.ID, omniworkboard.ColumnDone)

// B is unblocked
card, _ = board.GetCard(ctx, cardB.ID)
fmt.Println(card.BlockedBy) // []
```

## Visualizing Dependencies

### Get Dependency Graph

```go
graph := board.DependencyGraph(ctx)

for cardID, deps := range graph {
    card, _ := board.GetCard(ctx, cardID)
    fmt.Printf("%s depends on:\n", card.Title)
    for _, depID := range deps {
        dep, _ := board.GetCard(ctx, depID)
        fmt.Printf("  - %s\n", dep.Title)
    }
}
```

### Critical Path

Find the longest dependency chain:

```go
path := board.CriticalPath(ctx)
for i, cardID := range path {
    card, _ := board.GetCard(ctx, cardID)
    fmt.Printf("%d. %s\n", i+1, card.Title)
}
```

## Best Practices

### Keep Chains Short

Long chains increase risk:

```go
// Risky: long chain
A → B → C → D → E → F → G

// Better: parallel where possible
A → B ─┐
A → C ─┼→ G
A → D ─┘
```

### Explicit Dependencies

Only add necessary dependencies:

```go
// Don't add if not truly dependent
// ❌ "Write tests" depends on "Update README"

// ✓ Add if there's a real dependency
// "Write tests" depends on "Implement feature"
```

### Review Blocked Cards

Regularly check blocked cards:

```go
blocked := board.ListCards(ctx, omniworkboard.ListOptions{
    BlockedOnly: true,
})

if len(blocked) > 5 {
    fmt.Println("Warning: many blocked cards, review dependencies")
}
```

## Error Handling

```go
err := board.AddDependency(ctx, cardA.ID, cardB.ID)
if err != nil {
    switch {
    case errors.Is(err, omniworkboard.ErrCardNotFound):
        fmt.Println("Card not found")
    case errors.Is(err, omniworkboard.ErrCycleDetected):
        fmt.Println("Would create circular dependency")
    case errors.Is(err, omniworkboard.ErrSelfDependency):
        fmt.Println("Card cannot depend on itself")
    }
}
```

## See Also

- [Cards](cards.md) - Card operations
- [Boards](boards.md) - Board operations
- [Overview](overview.md) - Architecture overview
