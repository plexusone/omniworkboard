# Cards

Cards represent individual tasks or work items on a board.

## Card Structure

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

## Creating Cards

### Basic Card

```go
card, err := board.CreateCard(ctx,
    "Fix login bug",               // Title
    "Users can't login with SSO",  // Description
    omniworkboard.PriorityHigh,    // Priority
)
```

### Card with Options

```go
card, err := board.CreateCardWithOptions(ctx, omniworkboard.CreateCardOptions{
    Title:       "Implement feature",
    Description: "Add new feature X",
    Priority:    omniworkboard.PriorityNormal,
    Labels:      []string{"feature", "backend"},
    Metadata: map[string]any{
        "estimate": "3d",
        "assignee": "alice",
    },
})
```

## Priority Levels

```go
const (
    PriorityLow      Priority = 0
    PriorityNormal   Priority = 1
    PriorityHigh     Priority = 2
    PriorityCritical Priority = 3
)
```

### Priority Guidelines

| Priority | Response Time | Examples |
|----------|---------------|----------|
| Critical | Immediate | Production down, security issue |
| High | Same day | Important feature, blocking bug |
| Normal | This sprint | Regular work items |
| Low | When possible | Tech debt, nice-to-haves |

## Card Operations

### Get Card

```go
card, err := board.GetCard(ctx, cardID)
if err != nil {
    // Card not found
}
```

### Update Card

```go
err := board.UpdateCard(ctx, cardID, omniworkboard.UpdateCardOptions{
    Title:       "Updated title",
    Description: "Updated description",
    Priority:    omniworkboard.PriorityCritical,
})
```

### Move Card

```go
err := board.MoveCard(ctx, cardID, omniworkboard.ColumnInProgress)
```

### Delete Card

```go
err := board.DeleteCard(ctx, cardID)
```

### Add Labels

```go
err := board.AddLabel(ctx, cardID, "urgent")
err = board.AddLabel(ctx, cardID, "backend")
```

### Remove Labels

```go
err := board.RemoveLabel(ctx, cardID, "urgent")
```

## Listing Cards

### All Cards

```go
cards := board.ListCards(ctx, omniworkboard.ListOptions{})
```

### By Column

```go
inProgress := board.ListCards(ctx, omniworkboard.ListOptions{
    Column: omniworkboard.ColumnInProgress,
})
```

### By Priority

```go
urgent := board.ListCards(ctx, omniworkboard.ListOptions{
    Priority: omniworkboard.PriorityCritical,
})
```

### By Label

```go
bugs := board.ListCards(ctx, omniworkboard.ListOptions{
    Label: "bug",
})
```

### Blocked Cards

```go
blocked := board.ListCards(ctx, omniworkboard.ListOptions{
    BlockedOnly: true,
})
```

### Multiple Filters

```go
// High priority bugs in progress
cards := board.ListCards(ctx, omniworkboard.ListOptions{
    Column:   omniworkboard.ColumnInProgress,
    Priority: omniworkboard.PriorityHigh,
    Label:    "bug",
})
```

## Card Metadata

Store additional information in metadata:

```go
card.Metadata = map[string]any{
    "estimate":    "2d",
    "assignee":    "bob",
    "sprint":      42,
    "story_points": 5,
    "due_date":    "2026-06-15",
    "links": []string{
        "https://github.com/org/repo/issues/123",
    },
}
```

### Accessing Metadata

```go
if assignee, ok := card.Metadata["assignee"].(string); ok {
    fmt.Printf("Assigned to: %s\n", assignee)
}
```

## Completion

### Marking Complete

When a card moves to Done, it's automatically timestamped:

```go
board.MoveCard(ctx, cardID, omniworkboard.ColumnDone)

card, _ := board.GetCard(ctx, cardID)
fmt.Printf("Completed at: %v\n", card.CompletedAt)
```

### Moving Back

If moved back from Done, CompletedAt is cleared:

```go
board.MoveCard(ctx, cardID, omniworkboard.ColumnInProgress)

card, _ := board.GetCard(ctx, cardID)
// card.CompletedAt is nil
```

## Sorting

Cards are sorted by:

1. Priority (Critical first)
2. Created date (oldest first)

```go
cards := board.ListCards(ctx, omniworkboard.ListOptions{})
// Already sorted by priority, then creation time
```

## Best Practices

### Clear Titles

```go
// Good - specific and actionable
"Fix SSO login timeout error"
"Add pagination to user list API"
"Update Docker base image to Alpine 3.18"

// Bad - vague
"Fix bug"
"Update stuff"
"Work on feature"
```

### Useful Descriptions

Include:

- What needs to be done
- Why it's needed
- Acceptance criteria
- Links to related resources

```go
card, _ := board.CreateCard(ctx,
    "Add rate limiting to API",
    `## Description
Implement rate limiting to prevent abuse.

## Acceptance Criteria
- [ ] 100 requests per minute per user
- [ ] Return 429 status when exceeded
- [ ] Add X-RateLimit headers

## Links
- Design doc: https://...
- Issue: #123`,
    omniworkboard.PriorityHigh,
)
```

### Label Conventions

Establish consistent labels:

```
Types: bug, feature, enhancement, docs, chore
Areas: frontend, backend, api, database, infra
Priority: p0, p1, p2, p3
Status: blocked, needs-review, ready
```

## See Also

- [Boards](boards.md) - Board operations
- [Dependencies](dependencies.md) - Managing dependencies
- [Overview](overview.md) - Architecture overview
