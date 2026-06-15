# OmniWorkboard

**Project management workboard for AI agents.**

OmniWorkboard provides a kanban-style workboard that AI agents can use to track tasks, manage projects, and coordinate work.

## Key Features

- **Kanban Board** - Organize tasks across columns (Backlog, Todo, In Progress, Review, Done)
- **Priority System** - Mark tasks as Low, Normal, High, or Critical priority
- **Dependencies** - Define task dependencies and track blockers
- **Labels** - Categorize tasks with custom labels
- **Agent Skills** - Built-in skill implementation for AI agent integration
- **Persistence** - Store boards using kvs.Store backend

## Quick Example

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
        Name:        "Sprint 42",
        Description: "Q3 Feature Development",
    })

    // Create cards
    card1, _ := board.CreateCard(ctx, "Implement login", "Add OAuth2 login flow", omniworkboard.PriorityHigh)
    card2, _ := board.CreateCard(ctx, "Add tests", "Unit tests for auth module", omniworkboard.PriorityNormal)

    // Set dependency
    board.AddDependency(ctx, card2.ID, card1.ID) // card2 depends on card1

    // Move card through workflow
    board.MoveCard(ctx, card1.ID, omniworkboard.ColumnInProgress)

    // List cards
    cards := board.ListCards(ctx, omniworkboard.ListOptions{
        Column: omniworkboard.ColumnInProgress,
    })

    for _, c := range cards {
        fmt.Printf("%s: %s [%s]\n", c.ID[:8], c.Title, c.Column)
    }
}
```

## Use Cases

| Use Case | Description |
|----------|-------------|
| **Agent Task Tracking** | Track tasks an AI agent is working on |
| **Project Planning** | Plan and organize project work |
| **Sprint Management** | Manage agile sprints with backlogs |
| **Personal Tasks** | Simple personal task management |
| **Multi-Agent Coordination** | Coordinate work across multiple agents |

## Package Structure

```
github.com/plexusone/omniworkboard
├── board.go     # Board struct and operations
├── card.go      # Card struct and operations
├── column.go    # Column types and transitions
├── tools.go     # Tool definitions for agents
├── skill.go     # Skill implementation for omniagent
└── doc.go       # Package documentation
```

## Installation

```bash
go get github.com/plexusone/omniworkboard@latest
```

## Integration with OmniAgent

OmniWorkboard provides a compiled skill for seamless integration:

```go
import (
    "github.com/plexusone/omniagent/agent"
    "github.com/plexusone/omniworkboard"
)

// Create workboard skill
skill := omniworkboard.NewSkill()

// Add to agent
a, err := agent.New(config,
    agent.WithCompiledSkill(skill),
)
```

Once integrated, the agent can use workboard tools:

```
User: Create a card for "Fix login bug" with high priority
Agent: [calls create_card] Created card "Fix login bug" in Backlog with high priority.

User: Move it to In Progress
Agent: [calls move_card] Moved "Fix login bug" to In Progress.

User: What's on my board?
Agent: [calls list_cards] You have 3 cards:
- In Progress: Fix login bug (High)
- Todo: Add unit tests (Normal)
- Backlog: Update documentation (Low)
```

## Getting Started

- [Getting Started Guide](getting-started.md)
- [Concepts Overview](concepts/overview.md)
- [Agent Integration](agent-integration.md)
