# Agent Integration

OmniWorkboard provides a compiled skill for seamless integration with OmniAgent and other AI agent frameworks.

## Quick Start

```go
import (
    "github.com/plexusone/omniagent/agent"
    "github.com/plexusone/omniworkboard"
)

// Create skill
skill := omniworkboard.NewSkill()

// Add to agent
a, err := agent.New(config,
    agent.WithCompiledSkill(skill),
)
```

## Available Tools

The workboard skill provides these tools to the agent:

### create_card

Create a new task card.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `title` | string | Yes | Card title |
| `description` | string | No | Card description |
| `priority` | string | No | low, normal, high, critical |
| `labels` | array | No | List of labels |

**Example:**
```
User: Create a high priority card for "Fix authentication bug"
Agent: [calls create_card with title="Fix authentication bug", priority="high"]
       Created card "Fix authentication bug" in Backlog with high priority (ID: abc123)
```

### move_card

Move a card to a different column.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `card_id` | string | Yes | Card ID |
| `column` | string | Yes | Target column |

**Example:**
```
User: Move the auth bug card to In Progress
Agent: [calls move_card with card_id="abc123", column="in_progress"]
       Moved "Fix authentication bug" to In Progress
```

### list_cards

List cards with optional filters.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `column` | string | No | Filter by column |
| `priority` | string | No | Filter by priority |
| `label` | string | No | Filter by label |

**Example:**
```
User: What's in progress?
Agent: [calls list_cards with column="in_progress"]
       In Progress:
       1. Fix authentication bug (High) - abc123
       2. Update API endpoints (Normal) - def456
```

### update_card

Update card details.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `card_id` | string | Yes | Card ID |
| `title` | string | No | New title |
| `description` | string | No | New description |
| `priority` | string | No | New priority |
| `labels` | array | No | New labels |

### add_dependency

Add a dependency between cards.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `card_id` | string | Yes | Card that depends |
| `depends_on` | string | Yes | Card that blocks |

**Example:**
```
User: The tests card depends on the login card
Agent: [calls add_dependency with card_id="tests123", depends_on="login456"]
       Added dependency: "Write tests" now depends on "Implement login"
```

### get_board

Get the current board state.

No parameters required.

## Skill Configuration

```go
skill := omniworkboard.NewSkill(omniworkboard.SkillConfig{
    // Board name (creates new or loads existing)
    BoardName: "My Project",

    // Board description
    BoardDescription: "Project management board",

    // Custom columns (optional)
    Columns: []omniworkboard.ColumnType{
        omniworkboard.ColumnTodo,
        omniworkboard.ColumnInProgress,
        omniworkboard.ColumnDone,
    },
})
```

## Storage Integration

The skill implements `compiled.StorageAware` for automatic storage injection:

```go
import (
    "github.com/plexusone/omnistorage-core/kvs"
    "github.com/plexusone/omnistorage-core/sqlite"
)

// Create storage
backend, _ := sqlite.Open("agent.db")
store := kvs.New(backend)

// Agent with storage
a, err := agent.New(config,
    agent.WithStorage(store),
    agent.WithCompiledSkill(omniworkboard.NewSkill()),
)
// Skill automatically receives storage
```

## Conversation Examples

### Project Setup

```
User: Let's set up a sprint board. Create cards for:
      1. Design API schema (high priority)
      2. Implement endpoints
      3. Write tests
      4. Deploy to staging

Agent: I'll create those cards for you.

[creates 4 cards]

Created 4 cards in Backlog:
- Design API schema (High) - abc123
- Implement endpoints (Normal) - def456
- Write tests (Normal) - ghi789
- Deploy to staging (Normal) - jkl012

User: The endpoints depend on the design, tests depend on endpoints,
      and deploy depends on tests.

Agent: [adds 3 dependencies]

Dependencies set:
- "Implement endpoints" depends on "Design API schema"
- "Write tests" depends on "Implement endpoints"
- "Deploy to staging" depends on "Write tests"
```

### Daily Standup

```
User: What's my board status?

Agent: [calls list_cards]

Here's your board status:

Done (2):
- Design API schema (High)
- Implement endpoints (Normal)

In Progress (1):
- Write tests (Normal)

Blocked (1):
- Deploy to staging - blocked by "Write tests"

Backlog (0): empty
```

### Moving Through Workflow

```
User: I finished writing tests, move it to review

Agent: [calls move_card with column="review"]

Moved "Write tests" to Review.

Note: "Deploy to staging" is now unblocked and ready to start.
```

## Error Handling

The skill handles common errors gracefully:

```
User: Move card xyz123 to done

Agent: I couldn't find a card with ID "xyz123". Would you like me to
       list all cards so you can identify the correct one?
```

```
User: Move the deploy card to done

Agent: I can't move "Deploy to staging" to Done because it's blocked
       by "Write tests" which is still in Review. Would you like me
       to move "Write tests" to Done first?
```

## Best Practices

### Clear Task Descriptions

```
User: Create a card for "login"  # Too vague

User: Create a high priority card for "Implement OAuth2 login with
      Google and GitHub providers" with labels "auth" and "feature"
```

### Use Dependencies

```
User: Make sure the testing card can't be done until the
      implementation is complete

Agent: [adds dependency]
       "Testing" now depends on "Implementation"
```

### Regular Updates

```
User: What's blocking progress?

Agent: [analyzes dependencies]
       2 cards are blocked:
       - "Deployment" is blocked by "Testing" (in Review)
       - "Documentation" is blocked by "API Design" (in Progress)
```

## See Also

- [Getting Started](getting-started.md)
- [Cards](concepts/cards.md)
- [Dependencies](concepts/dependencies.md)
