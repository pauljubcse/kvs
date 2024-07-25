package kvs

import (
	"fmt"
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


func TestSkipListRank(t *testing.T) {
	sl := NewSkipList()

	// Insert some elements
	sl.Insert(3, "Value3")
	sl.Insert(6, "Value6")
	sl.Insert(7, "Value7")
	sl.Insert(9, "Value9")
	sl.Insert(12, "Value12")
	sl.Insert(19, "Value19")
	sl.Insert(21, "Value21")
	sl.Insert(25, "Value25")
	sl.Insert(26, "Value26")

	// Print levels for visual inspection
	fmt.Println("SkipList Levels:")
	sl.PrintLevels()

	// Define edge case test scenarios
	testCases := []struct {
		key      int
		expected int
	}{
		{3, 0},   // Rank of the first element
		{6, 1},   // Rank of the second element
		{7, 2},   // Rank of the third element
		{9, 3},   // Rank of the fourth element
		{12, 4},  // Rank of the fifth element
		{19, 5},  // Rank of the sixth element
		{21, 6},  // Rank of the seventh element
		{25, 7},  // Rank of the eighth element
		{26, 8},  // Rank of the ninth element
		{1, 0},  // Rank of an element smaller than the smallest element
		{5, 1},  // Rank of a non-existent element between existing elements
		{15, 5}, // Rank of a non-existent element between existing elements
		{30, 9}, // Rank of an element larger than the largest element
		{1, 0}, // Rank of a negative element
	}

	// Function to check the rank and compare with the expected value
	checkRank := func(key, expected int) {
		rank := sl.Rank(key)
		if expected == -1 {
			if rank != -1 {
				t.Errorf("Rank(%d) = %d; want %d", key, rank, expected)
			}
		} else {
			if rank != expected {
				t.Errorf("Rank(%d) = %d; want %d", key, rank, expected)
			}
		}
		fmt.Printf("Rank(%d) = %d; expected %d\n", key, rank, expected)
	}

	// Run the test cases
	for _, tc := range testCases {
		checkRank(tc.key, tc.expected)
	}

	// Additional edge cases with an empty list
	emptySL := NewSkipList()
	rank := emptySL.Rank(3)
	if rank != 0 {
		t.Errorf("Rank(%d) in an empty list = %d; want -1", 3, rank)
	}
	fmt.Printf("Rank(%d) in an empty list = %d; expected 0\n", 3, rank)
}