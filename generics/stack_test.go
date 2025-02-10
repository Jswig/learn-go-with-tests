package generics

import (
	"slices"
	"testing"
)

func TestPush(t *testing.T) {
	s := new(Stack[int])
	s.Push(1)

	want := []int{1}
	got := s.items
	if !slices.Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestIsEmpty(t *testing.T) {
	t.Run("empty stack", func(t *testing.T) {
		s := new(Stack[int])
		got := s.IsEmpty()
		want := true
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("stack with one item", func(t *testing.T) {
		s := new(Stack[int])
		s.Push(1)
		got := s.IsEmpty()
		want := false
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})
}

func TestPop(t *testing.T) {
	t.Run("empty stack", func(t *testing.T) {
		s := new(Stack[int])
		_, got := s.Pop()
		want := false
		if got != want {
			t.Errorf("got %t want %t", got, want)
		}
	})

	t.Run("stack with one item", func(t *testing.T) {
		s := new(Stack[int])
		s.Push(1)
		got, _ := s.Pop()
		want := 1
		if got != want {
			t.Errorf("got %v wwant %v", got, want)
		}
	})

	t.Run("stack with several items", func(t *testing.T) {
		s := new(Stack[int])
		s.Push(1)
		s.Push(2)
		got, _ := s.Pop()
		want := 2
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
