package kvs

import (
	"math/rand"
	"testing"
	"time"
)

// TestInsertion tests the insertion of elements into the skip list
func TestInsertion(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(3, "value3")
	sl.Insert(1, "value1")
	sl.Insert(7, "value7")
	sl.Insert(5, "value5")

	tests := []struct {
		key      int
		expected string
	}{
		{1, "value1"},
		{3, "value3"},
		{5, "value5"},
		{7, "value7"},
	}

	for _, test := range tests {
		value, found := sl.Search(test.key)
		if !found || value != test.expected {
			t.Errorf("Insert: expected %v for key %d, got %v", test.expected, test.key, value)
		}
	}
}

// TestDeletion tests the deletion of elements from the skip list
func TestDeletion(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(3, "value3")
	sl.Insert(1, "value1")
	sl.Insert(7, "value7")
	sl.Insert(5, "value5")

	sl.Delete(3)
	sl.Delete(1)

	tests := []struct {
		key      int
		expected string
		found    bool
	}{
		{1, "", false},
		{3, "", false},
		{5, "value5", true},
		{7, "value7", true},
	}

	for _, test := range tests {
		value, found := sl.Search(test.key)
		if found != test.found || (found && value != test.expected) {
			t.Errorf("Delete: expected %v for key %d, got %v", test.expected, test.key, value)
		}
	}
}

// TestSearch tests the search functionality of the skip list
func TestSearch(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(3, "value3")
	sl.Insert(1, "value1")
	sl.Insert(7, "value7")
	sl.Insert(5, "value5")

	tests := []struct {
		key      int
		expected string
		found    bool
	}{
		{1, "value1", true},
		{3, "value3", true},
		{5, "value5", true},
		{7, "value7", true},
		{2, "", false},
		{8, "", false},
	}

	for _, test := range tests {
		value, found := sl.Search(test.key)
		if found != test.found || (found && value != test.expected) {
			t.Errorf("Search: expected %v for key %d, got %v", test.expected, test.key, value)
		}
	}
}

// TestRandomLevel tests the random level generation
func TestRandomLevel(t *testing.T) {
	rand.Seed(1) // Fixed seed for predictable results

	levels := make([]int, MaxLevel)
	for i := 0; i < 10000; i++ {
		lvl := RandomLevel()
		if lvl < 1 || lvl > MaxLevel {
			t.Errorf("RandomLevel: generated level %d out of bounds", lvl)
		}
		levels[lvl-1]++
	}

	for i, count := range levels {
		t.Logf("Level %d: %d nodes", i+1, count)
	}
}

// TestPrintLevels tests the PrintLevels function (manual inspection required)
func TestPrintLevels(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(3, "value3")
	sl.Insert(1, "value1")
	sl.Insert(7, "value7")
	sl.Insert(5, "value5")

	t.Log("Skip List Level-wise:")
	sl.PrintLevels()

	sl.Delete(3)
	t.Log("Skip List Level-wise after deletion:")
	sl.PrintLevels()
}
func TestDeleteRange(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	sl := NewSkipList()

	sl.Insert(1, "value1")
	sl.Insert(2, "value2")
	sl.Insert(3, "value3")
	sl.Insert(4, "value4")
	sl.Insert(5, "value5")
	sl.Insert(6, "value6")
	sl.Insert(7, "value7")
	sl.Insert(8, "value8")
	sl.Insert(9, "value9")
	sl.PrintLevels()
	sl.DeleteRange(3, 7)
	sl.PrintLevels()
	tests := []struct {
		key   int
		exist bool
	}{
		{1, true},
		{2, true},
		{3, false},
		{4, false},
		{5, false},
		{6, false},
		{7, false},
		{8, true},
		{9, true},
	}

	for _, test := range tests {
		_, exist := sl.Search(test.key)
		if exist != test.exist {
			t.Errorf("Search(%d) = %v; want %v", test.key, exist, test.exist)
		}
	}
}
