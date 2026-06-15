# OmniWorkboard

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/omniworkboard/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/omniworkboard/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/omniworkboard/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/omniworkboard/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/omniworkboard/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/omniworkboard/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/omniworkboard
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/omniworkboard
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/omniworkboard
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/omniworkboard
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/omniworkboard
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fomniworkboard
 [loc-svg]: https://tokei.rs/b1/github/plexusone/omniworkboard
 [repo-url]: https://github.com/plexusone/omniworkboard
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/omniworkboard/blob/main/LICENSE

Project management workboard for AI agents. Provides a kanban-style board with cards, columns, priorities, and dependencies.

## Features

- 📋 **Kanban Columns**: Backlog, Todo, In Progress, Review, Done
- 🗂️ **Card Management**: Create, update, move, and delete task cards
- 🎯 **Priorities**: Low, Normal, High, Critical
- 🔗 **Dependencies**: Link cards with dependency tracking and cycle detection
- ✅ **Workflow Validation**: Enforces valid column transitions
- 🤖 **Agent Tools**: JSON-schema validated tools for AI agent integration
- 💾 **Persistence**: Optional storage via omnistorage-core

## Installation

```bash
go get github.com/plexusone/omniworkboard
```

## Usage

### Standalone Board

```go
package main

import (
    "context"
    "fmt"

    "github.com/plexusone/omniworkboard"
)

func main() {
    ctx := context.Background()

    // Create a board
    board := omniworkboard.NewBoard(omniworkboard.BoardConfig{
        Name: "My Project",
    })

    // Create a card
    card, _ := board.CreateCard(ctx, "Implement feature", "Description here", omniworkboard.PriorityHigh)
    fmt.Printf("Created card: %s\n", card.ID)

    // Move through workflow
    board.MoveCard(ctx, card.ID, omniworkboard.ColumnTodo)
    board.MoveCard(ctx, card.ID, omniworkboard.ColumnInProgress)
    board.MoveCard(ctx, card.ID, omniworkboard.ColumnDone)

    // Get stats
    stats := board.Stats()
    fmt.Printf("Total cards: %d\n", stats.TotalCards)
}
```

### As an OmniSkill

```go
package main

import (
    "context"

    "github.com/plexusone/omniworkboard"
)

func main() {
    ctx := context.Background()

    // Create skill for agent runtime
    skill := omniworkboard.NewSkill()
    skill.Init(ctx)
    defer skill.Close()

    // Get available tools
    tools := skill.Tools()
    for _, tool := range tools {
        fmt.Printf("Tool: %s\n", tool.Name())
    }
}
```

## Tools

The workboard exposes these tools for agent integration:

| Tool | Description |
|------|-------------|
| `create_card` | Create a new task card |
| `move_card` | Move a card to a different column |
| `list_cards` | List cards with optional filters |
| `update_card` | Update card details |
| `add_dependency` | Add a dependency between cards |
| `get_card` | Get details of a specific card |
| `delete_card` | Delete a card from the board |
| `board_stats` | Get workboard statistics |

## Columns

Cards flow through these columns:

```
Backlog -> Todo -> In Progress -> Review -> Done
```

Backward transitions are allowed for rework (e.g., Review -> In Progress).

## Dependencies

Cards can depend on other cards. A card is blocked until all its dependencies are in the Done column.

```go
// card2 depends on card1
board.AddDependency(ctx, card2.ID, card1.ID)

// card2 cannot move forward until card1 is done
_, err := board.MoveCard(ctx, card2.ID, omniworkboard.ColumnInProgress)
// err: card is blocked
```

Circular dependencies are detected and rejected.

## License

MIT License - see [LICENSE](LICENSE) for details.
