package sortedset

import (
	"github.com/fioncat/datastc/pkg/skiplist"
	"github.com/fioncat/datastc/pkg/types"
)

type SortedSet struct {
	dict map[string]float64
	zsl  *skiplist.SkipList
}

func New(length int) *SortedSet {
	return &SortedSet{
		dict: make(map[string]float64, length),
		zsl:  skiplist.New(),
	}
}

func (set *SortedSet) Set(key string, score float64) bool {
	curScore, ok := set.dict[key]
	if ok && curScore == score {
		return false
	}
	set.dict[key] = score
	if !ok {
		set.zsl.Insert(score, key)
		return true
	}

	set.zsl.UpdateScore(curScore, key, score)
	return true
}

func (set *SortedSet) Slice() []types.ScoreKey {
	return types.ScoreValue2Key(set.zsl.ToSlice(false))
}

func (set *SortedSet) SliceDesc() []types.ScoreKey {
	return types.ScoreValue2Key(set.zsl.ToSlice(true))
}

func (set *SortedSet) Get(key string) (float64, bool) {
	score, ok := set.dict[key]
	return score, ok
}

func (set *SortedSet) Delete(key string) bool {
	score, ok := set.dict[key]
	if !ok {
		return false
	}
	delete(set.dict, key)
	set.zsl.Delete(score, key)
	return true
}
