package skiplist

import (
	"bytes"
	"ds/fastrand"
	"errors"
	"math"
)

const (
	MaxHeight = 16  // 最高层数，这是论文里的层数，当然也可以是其他值，我就不深入了
	PValue    = 0.5 // p = 1/2
)

var probabilities [MaxHeight]uint32

func init() {
	probability := 1.0

	for level := 0; level < MaxHeight; level++ {
		probabilities[level] = uint32(probability * float64(math.MaxUint32))
		probability *= PValue
	}
}

type node struct {
	key []byte
	val []byte
	// "tower"是一个前向指针的集合，它将节点与跳转列表中每个相应层的后续节点连接起来。
	tower [MaxHeight]*node
}

type SkipList struct {
	head   *node
	height int
}

func NewSkipList() *SkipList {
	sl := &SkipList{}
	sl.head = &node{} // 哨兵节点
	sl.height = 1
	return sl
}

func (sl *SkipList) search(key []byte) (*node, [MaxHeight]*node) {
	var next *node
	var journey [MaxHeight]*node

	prev := sl.head
	for level := sl.height - 1; level >= 0; level-- {
		for next = prev.tower[level]; next != nil; next = prev.tower[level] {
			if bytes.Compare(key, next.key) <= 0 {
				break
			}
			prev = next
		}
		journey[level] = prev
	}

	if next != nil && bytes.Equal(key, next.key) {
		return next, journey
	}
	return nil, journey
}

func (sl *SkipList) Find(key []byte) ([]byte, error) {
	found, _ := sl.search(key)

	if found == nil {
		return nil, errors.New("key not found")
	}

	return found.val, nil
}

func randomHeight() int {
	seed := fastrand.Uint32()

	height := 1
	for height < MaxHeight && seed <= probabilities[height] {
		height++
	}

	return height
}

func (sl *SkipList) Insert(key []byte, val []byte) {
	found, journey := sl.search(key)

	if found != nil {
		// update value of existing key
		found.val = val
		return
	}
	height := randomHeight()
	nd := &node{key: key, val: val}

	for level := 0; level < height; level++ {
		prev := journey[level]

		if prev == nil {
			// prev is nil if we are extending the height of the tree,
			// because that level did not exist while the journey was being recorded
			prev = sl.head
		}
		nd.tower[level] = prev.tower[level]
		prev.tower[level] = nd
	}

	if height > sl.height {
		sl.height = height
	}
}
