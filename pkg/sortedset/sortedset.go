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

func (set *SortedSet) Set(score float64, key string) bool {
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

type ScanFunc func(score float64, key string) bool

func (set *SortedSet) Scan(f ScanFunc) {
	set.zsl.Scan(false, func(score float64, value interface{}) bool {
		return f(score, value.(string))
	})
}

func (set *SortedSet) ScanDesc(f ScanFunc) {
	set.zsl.Scan(true, func(score float64, value interface{}) bool {
		return f(score, value.(string))
	})
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

func (set *SortedSet) GetRange(r types.Range) []types.ScoreKey {
	return types.ScoreValue2Key(set.zsl.GetRange(r))
}

func (set *SortedSet) ScanRange(r types.Range, f ScanFunc) {
	set.zsl.ScanRange(r, func(score float64, value interface{}) bool {
		return f(score, value.(string))
	})
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

func (set *SortedSet) DeleteRange(r types.Range) int {
	return set.zsl.DeleteRange(r, func(score float64, value interface{}) {
		delete(set.dict, value.(string))
	})
}
