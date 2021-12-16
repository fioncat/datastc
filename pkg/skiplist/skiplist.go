package skiplist

import (
	"math/rand"
	"time"

	"github.com/fioncat/datastc/pkg/types"
)

const (
	maxLevel = 32

	levelRate int32 = 2
	randComp  int32 = 0xFFFF / levelRate
)

type node struct {
	level    []nodeLevel
	backward *node

	score float64
	value interface{}
}

type nodeLevel struct {
	forward *node
	span    int
}

type SkipList struct {
	header *node
	tail   *node

	length int
	level  int

	rand *rand.Rand
}

func New() *SkipList {
	return &SkipList{
		header: &node{
			level:    make([]nodeLevel, maxLevel),
			backward: nil,
		},

		tail: nil,

		length: 0,
		level:  1,

		rand: rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (zsl *SkipList) ToSlice(r bool) []types.ScoreValue {
	items := make([]types.ScoreValue, 0, zsl.length)
	var x *node
	if !r {
		x = zsl.header.level[0].forward
	} else {
		x = zsl.tail
	}
	for x != nil {
		items = append(items, types.ScoreValue{
			Score: x.score,
			Value: x.value,
		})
		if !r {
			x = x.level[0].forward
		} else {
			x = x.backward
		}
	}
	return items
}

func (zsl *SkipList) randLevel() int {
	level := 1
	for zsl.rand.Int31()&0xFFFF < randComp {
		level++
	}
	if level < maxLevel {
		return level
	}
	return maxLevel
}

func (zsl *SkipList) Insert(score float64, val interface{}) {
	var update [maxLevel]*node
	var rank [maxLevel]int

	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		if i == zsl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score &&
					types.Compare(x.level[i].forward.value, val) < 0)) {
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		update[i] = x
	}

	level := zsl.randLevel()
	if level > zsl.level {
		for i := zsl.level; i < level; i++ {
			rank[i] = 0
			update[i] = zsl.header
			update[i].level[i].span = zsl.length
		}
		zsl.level = level
	}

	x = &node{
		level:    make([]nodeLevel, level),
		backward: nil,

		score: score,
		value: val,
	}
	for i := 0; i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	for i := level; i < zsl.level; i++ {
		update[i].level[i].span++
	}

	if update[0] == zsl.header {
		x.backward = nil
	} else {
		x.backward = update[0]
	}

	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		zsl.tail = x
	}
	zsl.length++
}

func (zsl *SkipList) deleteNode(x *node, update [maxLevel]*node) {
	for i := 0; i < zsl.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span--
		}
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		zsl.tail = x.backward
	}

	for zsl.level > 1 && zsl.header.level[zsl.level-1].forward == nil {
		zsl.level--
	}
	zsl.length--
}

func (zsl *SkipList) Delete(score float64, value interface{}) bool {
	var update [maxLevel]*node
	x := zsl.header

	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score &&
					types.Compare(x.level[i].forward.value, value) < 0)) {
			x = x.level[i].forward
		}
		update[i] = x
	}

	x = x.level[0].forward
	if x != nil && x.score == score && types.Compare(x.value, value) == 0 {
		zsl.deleteNode(x, update)
		return true
	}
	return false
}

func (zsl *SkipList) UpdateScore(curScore float64, value interface{}, newScore float64) bool {
	var update [maxLevel]*node

	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < curScore ||
				(x.level[i].forward.score == curScore &&
					types.Compare(x.level[i].forward.value, value) < 0)) {
			x = x.level[i].forward
		}
		update[i] = x
	}

	x = x.level[0].forward
	if x == nil || x.score != curScore || types.Compare(x.value, value) != 0 {
		return false
	}

	if (x.backward == nil || x.backward.score < newScore) &&
		(x.level[0].forward == nil || x.level[0].forward.score > newScore) {
		x.score = newScore
		return true
	}

	zsl.deleteNode(x, update)
	zsl.Insert(newScore, value)
	return true
}

func (zsl *SkipList) GetRange(minScore, maxScore float64) []types.ScoreValue {
	if minScore > maxScore {
		return nil
	}

	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.score < minScore {
			x = x.level[i].forward
		}
	}

	x = x.level[0].forward
	items := make([]types.ScoreValue, 0)
	for x != nil {
		if x.score > maxScore {
			break
		}
		items = append(items, types.ScoreValue{
			Score: x.score,
			Value: x.value,
		})
		x = x.level[0].forward
	}

	return items
}
