package omniworkboard

import (
	"context"
	"testing"
)

func TestBoard_CreateAndMoveCard(t *testing.T) {
	board := NewBoard(BoardConfig{Name: "Test Board"})
	ctx := context.Background()

	// Create a card
	card, err := board.CreateCard(ctx, "Test Task", "Description", PriorityNormal)
	if err != nil {
		t.Fatalf("CreateCard failed: %v", err)
	}

	if card.Column != ColumnBacklog {
		t.Errorf("Expected column %s, got %s", ColumnBacklog, card.Column)
	}

	// Move to todo
	card, err = board.MoveCard(ctx, card.ID, ColumnTodo)
	if err != nil {
		t.Fatalf("MoveCard to todo failed: %v", err)
	}

	if card.Column != ColumnTodo {
		t.Errorf("Expected column %s, got %s", ColumnTodo, card.Column)
	}

	// Move to in_progress
	card, err = board.MoveCard(ctx, card.ID, ColumnInProgress)
	if err != nil {
		t.Fatalf("MoveCard to in_progress failed: %v", err)
	}

	// Move to done
	card, err = board.MoveCard(ctx, card.ID, ColumnDone)
	if err != nil {
		t.Fatalf("MoveCard to done failed: %v", err)
	}

	if card.CompletedAt == nil {
		t.Error("Expected CompletedAt to be set")
	}
}

func TestBoard_Dependencies(t *testing.T) {
	board := NewBoard(BoardConfig{Name: "Test Board"})
	ctx := context.Background()

	// Create two cards
	card1, err := board.CreateCard(ctx, "Task 1", "", PriorityNormal)
	if err != nil {
		t.Fatal(err)
	}

	card2, err := board.CreateCard(ctx, "Task 2", "", PriorityNormal)
	if err != nil {
		t.Fatal(err)
	}

	// Add dependency: card2 depends on card1
	if err := board.AddDependency(ctx, card2.ID, card1.ID); err != nil {
		t.Fatalf("AddDependency failed: %v", err)
	}

	// Get card2 and check it's blocked
	card2, err = board.GetCard(ctx, card2.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(card2.BlockedBy) != 1 || card2.BlockedBy[0] != card1.ID {
		t.Errorf("Expected card2 to be blocked by card1, got %v", card2.BlockedBy)
	}

	// Try to move card2 to in_progress - should fail
	_, err = board.MoveCard(ctx, card2.ID, ColumnInProgress)
	if err == nil {
		t.Error("Expected error when moving blocked card")
	}

	// Complete card1
	_, err = board.MoveCard(ctx, card1.ID, ColumnTodo)
	if err != nil {
		t.Fatal(err)
	}
	_, err = board.MoveCard(ctx, card1.ID, ColumnDone)
	if err != nil {
		t.Fatal(err)
	}

	// Now card2 should be unblocked
	card2, err = board.GetCard(ctx, card2.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(card2.BlockedBy) != 0 {
		t.Errorf("Expected card2 to be unblocked after card1 done, got %v", card2.BlockedBy)
	}

	// Now moving card2 should work
	_, err = board.MoveCard(ctx, card2.ID, ColumnTodo)
	if err != nil {
		t.Fatalf("MoveCard should succeed now: %v", err)
	}
}

func TestBoard_CyclicDependency(t *testing.T) {
	board := NewBoard(BoardConfig{Name: "Test Board"})
	ctx := context.Background()

	card1, _ := board.CreateCard(ctx, "Task 1", "", PriorityNormal)
	card2, _ := board.CreateCard(ctx, "Task 2", "", PriorityNormal)
	card3, _ := board.CreateCard(ctx, "Task 3", "", PriorityNormal)

	// Create chain: card3 -> card2 -> card1
	if err := board.AddDependency(ctx, card2.ID, card1.ID); err != nil {
		t.Fatal(err)
	}
	if err := board.AddDependency(ctx, card3.ID, card2.ID); err != nil {
		t.Fatal(err)
	}

	// Try to create cycle: card1 -> card3
	err := board.AddDependency(ctx, card1.ID, card3.ID)
	if err == nil {
		t.Error("Expected error for cyclic dependency")
	}
}

func TestBoard_ListCards(t *testing.T) {
	board := NewBoard(BoardConfig{Name: "Test Board"})
	ctx := context.Background()

	// Create cards with different priorities
	_, _ = board.CreateCard(ctx, "Low Priority", "", PriorityLow)
	_, _ = board.CreateCard(ctx, "High Priority", "", PriorityHigh)
	_, _ = board.CreateCard(ctx, "Normal Priority", "", PriorityNormal)

	// List all cards
	cards, err := board.ListCards(ctx, CardFilter{})
	if err != nil {
		t.Fatal(err)
	}

	if len(cards) != 3 {
		t.Fatalf("Expected 3 cards, got %d", len(cards))
	}

	// Should be sorted by priority (highest first)
	if cards[0].Priority != PriorityHigh {
		t.Errorf("Expected first card to be high priority, got %s", cards[0].Priority.String())
	}

	// Filter by priority
	cards, err = board.ListCards(ctx, CardFilter{Priorities: []Priority{PriorityHigh}})
	if err != nil {
		t.Fatal(err)
	}

	if len(cards) != 1 {
		t.Errorf("Expected 1 high priority card, got %d", len(cards))
	}
}
