// package kvs

// import (
// 	"fmt"
// 	"math/rand"
// )

// const MaxLevel int = 16
// const P float32 = 0.5

// //A Node in the SkipList
// type Node struct {
// 	key int
// 	//value interface{}
// 	value string
// 	forward []*Node
// }

// type SkipList struct {
// 	header *Node
// 	level int
// }

// func NewNode(level int, key int, value string) *Node{
// 	return &Node{
// 		key: key,
// 		value: value,
// 		forward: make([]*Node, level),
// 	}
// }

// func NewSkipList() *SkipList{
// 	return &SkipList{
// 		header: NewNode(MaxLevel, -1, ""),
// 		level: 1,
// 	}
// }

// func RandomLevel() int {
// 	level := 1
// 	for rand.Float32() < P && level < MaxLevel {
// 		level++
// 	}
// 	return level
// }

// func (sl *SkipList) Search(key int) (string, bool) {
// 	x := sl.header
// 	for i := sl.level-1; i>=0; i-- {
// 		for x.forward[i]!=nil && x.forward[i].key < key {
// 			x = x.forward[i]
// 		}
// 	}
// 	x=x.forward[0]
// 	if x!=nil && x.key == key{
// 		return x.value, true
// 	}
// 	return "", false
// }

// func (sl *SkipList) Insert(key int, value string) {
// 	update := make([]*Node, MaxLevel)
// 	x := sl.header

// 	for i := sl.level - 1; i >= 0; i-- {
// 		for x.forward[i] != nil && x.forward[i].key < key {
// 			x = x.forward[i]
// 		}
// 		update[i] = x
// 	}

// 	x = x.forward[0]

// 	if x != nil && x.key == key {
// 		x.value = value
// 	} else {
// 		lvl := RandomLevel()
// 		if lvl > sl.level {
// 			for i := sl.level; i < lvl; i++ {
// 				update[i] = sl.header
// 			}
// 			sl.level = lvl
// 		}
// 		x = NewNode(lvl, key, value)
// 		for i := 0; i < lvl; i++ {
// 			x.forward[i] = update[i].forward[i]
// 			update[i].forward[i] = x
// 		}
// 	}
// }

// func (sl *SkipList) Delete(key int) {
// 	update := make([]*Node, MaxLevel)
// 	x := sl.header

// 	for i := sl.level - 1; i >= 0; i-- {
// 		for x.forward[i] != nil && x.forward[i].key < key {
// 			x = x.forward[i]
// 		}
// 		update[i] = x
// 	}

// 	x = x.forward[0]

// 	if x != nil && x.key == key {
// 		for i := 0; i < sl.level; i++ {
// 			if update[i].forward[i] != x {
// 				break
// 			}
// 			update[i].forward[i] = x.forward[i]
// 		}
// 		for sl.level > 1 && sl.header.forward[sl.level-1] == nil {
// 			sl.level--
// 		}
// 	}
// }
// func (sl *SkipList) DeleteRange(startKey, endKey int) {
// 	update := make([]*Node, MaxLevel)
// 	x := sl.header

// 	for i := sl.level - 1; i >= 0; i-- {
// 		for x.forward[i] != nil && x.forward[i].key < startKey {
// 			x = x.forward[i]
// 		}
// 		update[i] = x
// 	}

// 	x = x.forward[0]

// 	for x != nil && x.key <= endKey {
// 		next := x.forward[0]
// 		sl.Delete(x.key)
// 		x = next
// 	}
// }
// func (sl *SkipList) PrintLevels() {
// 	for i := sl.level - 1; i >= 0; i-- {
// 		fmt.Printf("Level %d: ", i+1)
// 		x := sl.header.forward[i]
// 		for x != nil {
// 			fmt.Printf("%d ", x.key)
// 			x = x.forward[i]
// 		}
// 		fmt.Println()
// 	}
// }

// func (sl *SkipList) Rank(key int) int {
// 	count := 0
// 	x := sl.header

// 	// Traverse from the top level to the bottom level
// 	for i := sl.level - 1; i >= 0; i-- {
// 		for x.forward[i] != nil && x.forward[i].key < key {
// 			count += 1 // Count this element
// 			x = x.forward[i]
// 		}
// 	}

// 	return count
// }

package kvs

import (
	"fmt"
	"math/rand"
)

const MaxLevel int = 16
const P float32 = 0.5

// A Node in the SkipList
type Node struct {
	key     int
	value   string
	forward []*Node
	span    []int
}

type SkipList struct {
	header *Node
	level  int
}

func NewNode(level int, key int, value string) *Node {
	return &Node{
		key:     key,
		value:   value,
		forward: make([]*Node, level),
		span:    make([]int, level),
	}
}

func NewSkipList() *SkipList {
	return &SkipList{
		header: NewNode(MaxLevel, -1, ""),
		level:  1,
	}
}

func RandomLevel() int {
	level := 1
	for rand.Float32() < P && level < MaxLevel {
		level++
	}
	return level
}

func (sl *SkipList) Search(key int) (string, bool) {
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
	}
	x = x.forward[0]
	if x != nil && x.key == key {
		return x.value, true
	}
	return "", false
}

func (sl *SkipList) Insert(key int, value string) {
	update := make([]*Node, MaxLevel)
	rank := make([]int, MaxLevel)
	x := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for x.forward[i] != nil && x.forward[i].key < key {
			rank[i] += x.span[i]
			x = x.forward[i]
		}
		update[i] = x
	}

	level := RandomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.header
			update[i].span[i] = sl.header.span[i]
		}
		sl.level = level
	}

	node := NewNode(level, key, value)
	for i := 0; i < level; i++ {
		node.forward[i] = update[i].forward[i]
		update[i].forward[i] = node

		node.span[i] = update[i].span[i] - (rank[0] - rank[i])
		update[i].span[i] = (rank[0] - rank[i]) + 1
	}

	for i := level; i < sl.level; i++ {
		update[i].span[i]++
	}
}

func (sl *SkipList) Delete(key int) {
	update := make([]*Node, MaxLevel)
	x := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[0]
	if x != nil && x.key == key {
		for i := 0; i < sl.level; i++ {
			if update[i].forward[i] != x {
				break
			}
			update[i].span[i] += x.span[i] - 1
			update[i].forward[i] = x.forward[i]
		}

		for sl.level > 1 && sl.header.forward[sl.level-1] == nil {
			sl.level--
		}
	}
}

func (sl *SkipList) DeleteRange(startKey, endKey int) {
	update := make([]*Node, MaxLevel)
	x := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < startKey {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[0]
	for x != nil && x.key <= endKey {
		next := x.forward[0]
		sl.Delete(x.key)
		x = next
	}
}

func (sl *SkipList) PrintLevels() {
	for i := sl.level - 1; i >= 0; i-- {
		fmt.Printf("Level %d: ", i+1)
		x := sl.header.forward[i]
		for x != nil {
			fmt.Printf("%d ", x.key)
			x = x.forward[i]
		}
		fmt.Println()
	}
}

// func (sl *SkipList) Rank(key int) int {
// 	rank := 0
// 	x := sl.header
// 	for i := sl.level - 1; i >= 0; i-- {
// 		for x.forward[i] != nil && x.forward[i].key < key {
// 			rank += x.span[i]
// 			x = x.forward[i]
// 		}
// 	}
// 	return rank
// }
func (sl *SkipList) Rank(key int) int {
	count := 0
	x := sl.header

	// Traverse from the top level to the bottom level
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			count += x.span[i] // Accumulate the span
			x = x.forward[i]
		}
	}

	// Add the rank based on the position in the last level
	if x != nil && x.key == key {
		count += 1 // The rank of the element itself
	}

	return count
}
