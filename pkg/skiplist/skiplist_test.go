package skiplist

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/fioncat/datastc/pkg/types"
)

func items2strs(items []types.ScoreValue) []string {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = fmt.Sprintf("%.1f: %v", item.Score, item.Value)
	}
	return strs
}

func TestInsert(t *testing.T) {
	cases := []struct {
		data []types.ScoreValue

		results []string

		reverse bool
	}{
		{
			data:    []types.ScoreValue{},
			results: []string{},
		},
		{
			data:    []types.ScoreValue{},
			results: []string{},
			reverse: true,
		},
		{
			data: []types.ScoreValue{
				{Score: 0.3, Value: "c"},
				{Score: 0.2, Value: "b"},
				{Score: 0.4, Value: "d"},
				{Score: 0.6, Value: "f"},
				{Score: 0.5, Value: "e"},
				{Score: 0.1, Value: "a"},
			},

			results: []string{
				"0.1: a", "0.2: b", "0.3: c",
				"0.4: d", "0.5: e", "0.6: f",
			},
		},
		{
			data: []types.ScoreValue{
				{Score: 0.2, Value: "b"},
				{Score: 0.3, Value: "c"},
				{Score: 0.6, Value: "f"},
				{Score: 0.4, Value: "d"},
				{Score: 0.1, Value: "a"},
				{Score: 0.5, Value: "e"},
			},

			results: []string{
				"0.6: f", "0.5: e", "0.4: d",
				"0.3: c", "0.2: b", "0.1: a",
			},
			reverse: true,
		},
		{
			data: []types.ScoreValue{
				{Score: 0.2, Value: "aa"},
				{Score: 0.3, Value: "ss"},
				{Score: 0.1, Value: "a"},
				{Score: 0.1, Value: "z"},
				{Score: 0.2, Value: "zz"},
				{Score: 0.1, Value: "c"},
				{Score: 0.2, Value: "bb"},
				{Score: 0.1, Value: "b"},
			},

			results: []string{
				"0.1: a", "0.1: b", "0.1: c",
				"0.1: z", "0.2: aa", "0.2: bb",
				"0.2: zz", "0.3: ss",
			},
		},
		{
			data: []types.ScoreValue{
				{Score: 0.2, Value: "aa"},
				{Score: 0.3, Value: "aa"},
				{Score: 0.1, Value: "c"},
				{Score: 0.1, Value: "z"},
				{Score: 0.1, Value: "a"},
				{Score: 0.2, Value: "zz"},
				{Score: 0.3, Value: "bb"},
				{Score: 0.2, Value: "bb"},
				{Score: 0.1, Value: "b"},
			},

			results: []string{
				"0.3: bb", "0.3: aa", "0.2: zz",
				"0.2: bb", "0.2: aa", "0.1: z",
				"0.1: c", "0.1: b", "0.1: a",
			},

			reverse: true,
		},
	}

	for _, c := range cases {
		zsl := New()
		for _, item := range c.data {
			zsl.Insert(item.Score, item.Value)
		}
		slice := zsl.ToSlice(c.reverse)
		strs := items2strs(slice)
		if !reflect.DeepEqual(strs, c.results) {
			t.Fatalf("Expect: %v, get: %v", c.results, strs)
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		data []types.ScoreValue
		del  []types.ScoreValue

		results []string
	}{
		{
			data: []types.ScoreValue{},
			del: []types.ScoreValue{
				{Score: 0.1, Value: "a"},
				{Score: 0.2, Value: "b"},
			},
			results: []string{},
		},
		{
			data: []types.ScoreValue{
				{Score: 0.1, Value: "a"},
				{Score: 0.2, Value: "b"},
			},
			del: []types.ScoreValue{},
			results: []string{
				"0.1: a", "0.2: b",
			},
		},
		{
			data: []types.ScoreValue{
				{Score: 0.1, Value: "a"},
				{Score: 0.2, Value: "b"},
				{Score: 0.3, Value: "c"},
				{Score: 0.4, Value: "d"},
				{Score: 0.5, Value: "e"},
				{Score: 0.6, Value: "f"},
			},

			del: []types.ScoreValue{
				{Score: 0.1, Value: "a"},
				{Score: 0.3, Value: "c"},
				{Score: 0.5, Value: "zzz"}, // Not found
			},

			results: []string{
				"0.2: b", "0.4: d", "0.5: e", "0.6: f",
			},
		},
		{
			data: []types.ScoreValue{
				{Score: 0.1, Value: "aa"},
				{Score: 0.1, Value: "bb"},
				{Score: 0.1, Value: "cc"},
				{Score: 0.2, Value: "dd"},
				{Score: 0.2, Value: "ee"},
				{Score: 0.2, Value: "ff"},
			},

			del: []types.ScoreValue{
				{Score: 0.1, Value: "bb"},
				{Score: 0.1, Value: "cc"},
				{Score: 0.2, Value: "ff"},
			},

			results: []string{
				"0.1: aa", "0.2: dd", "0.2: ee",
			},
		},
		{
			data: []types.ScoreValue{
				{Score: 0.1, Value: "aa"},
				{Score: 0.1, Value: "cc"},
				{Score: 0.2, Value: "ee"},
				{Score: 0.1, Value: "bb"},
				{Score: 0.2, Value: "ff"},
				{Score: 0.2, Value: "dd"},
			},

			del: []types.ScoreValue{
				{Score: 0.1, Value: "aa"},
				{Score: 0.2, Value: "dd"},
				{Score: 0.2, Value: "ee"},
			},

			results: []string{
				"0.1: bb", "0.1: cc", "0.2: ff",
			},
		},
	}

	for _, c := range cases {
		zsl := New()
		for _, item := range c.data {
			zsl.Insert(item.Score, item.Value)
		}
		for _, item := range c.del {
			zsl.Delete(item.Score, item.Value)
		}
		slice := zsl.ToSlice(false)
		strs := items2strs(slice)
		if !reflect.DeepEqual(strs, c.results) {
			t.Fatalf("Expect: %v, get: %v", c.results, strs)
		}
	}
}

func TestUpdate(t *testing.T) {
	type updateNode struct {
		item  types.ScoreValue
		score float64
	}
	cases := []struct {
		data   []types.ScoreValue
		update []updateNode

		results []string
	}{
		{
			data: []types.ScoreValue{},
			update: []updateNode{
				{item: types.ScoreValue{Score: 0.3, Value: "c"}, score: 0.9},
				{item: types.ScoreValue{Score: 0.3, Value: "dd"}, score: 0.9},
			},

			results: []string{},
		},
		{
			data: []types.ScoreValue{
				{Score: 0.2, Value: "b"},
				{Score: 0.6, Value: "f"},
				{Score: 0.1, Value: "a"},
				{Score: 0.5, Value: "e"},
				{Score: 0.3, Value: "c"},
				{Score: 0.4, Value: "d"},
			},

			update: []updateNode{
				{item: types.ScoreValue{Score: 0.2, Value: "b"}, score: 0.1},
				{item: types.ScoreValue{Score: 0.3, Value: "c"}, score: 0.9},
				{item: types.ScoreValue{Score: 0.3, Value: "dd"}, score: 0.9},
			},

			results: []string{
				"0.1: a", "0.1: b", "0.4: d", "0.5: e", "0.6: f", "0.9: c",
			},
		},
		{
			data: []types.ScoreValue{
				{Score: 0.2, Value: "b"},
				{Score: 0.1, Value: "a"},
				{Score: 0.3, Value: "c"},
				{Score: 0.4, Value: "d"},
			},

			update: []updateNode{
				{item: types.ScoreValue{Score: 0.1, Value: "a"}, score: 0.0},
				{item: types.ScoreValue{Score: 0.2, Value: "b"}, score: 0.1},
			},

			results: []string{
				"0.0: a", "0.1: b", "0.3: c", "0.4: d",
			},
		},
	}

	for _, c := range cases {
		zsl := New()
		for _, item := range c.data {
			zsl.Insert(item.Score, item.Value)
		}
		for _, update := range c.update {
			zsl.UpdateScore(update.item.Score, update.item.Value, update.score)
		}
		slice := zsl.ToSlice(false)
		strs := items2strs(slice)
		if !reflect.DeepEqual(strs, c.results) {
			t.Fatalf("Expect: %v, get: %v", c.results, strs)
		}
	}
}

func TestGetByScoreRange(t *testing.T) {
	zsl := New()
	zsl.Insert(0.1, "a")
	zsl.Insert(0.2, "b")
	zsl.Insert(0.3, "c")
	zsl.Insert(0.4, "d")
	zsl.Insert(0.5, "e")

	cases := []struct {
		min, max float64

		zsl *SkipList

		results []string
	}{
		{
			min: 0.2, max: 0.1,
			zsl: zsl,

			results: []string{},
		},
		{
			min: 0.1, max: 0.9,
			zsl: New(),

			results: []string{},
		},
		{
			min: 0.0, max: 0.9,
			zsl: zsl,

			results: []string{
				"0.1: a", "0.2: b", "0.3: c",
				"0.4: d", "0.5: e",
			},
		},
		{
			min: 0.1, max: 0.9,
			zsl: zsl,

			results: []string{
				"0.1: a", "0.2: b", "0.3: c",
				"0.4: d", "0.5: e",
			},
		},
		{
			min: 0.0, max: 0.5,
			zsl: zsl,

			results: []string{
				"0.1: a", "0.2: b", "0.3: c",
				"0.4: d", "0.5: e",
			},
		},
		{
			min: 0.1, max: 0.5,
			zsl: zsl,

			results: []string{
				"0.1: a", "0.2: b", "0.3: c",
				"0.4: d", "0.5: e",
			},
		},
		{
			min: 0.1, max: 0.3,
			zsl: zsl,

			results: []string{
				"0.1: a", "0.2: b", "0.3: c",
			},
		},
		{
			min: 0.1, max: 0.35,
			zsl: zsl,

			results: []string{
				"0.1: a", "0.2: b", "0.3: c",
			},
		},
		{
			min: 0.15, max: 0.35,
			zsl: zsl,

			results: []string{
				"0.2: b", "0.3: c",
			},
		},
		{
			min: 0.15, max: 0.17,
			zsl: zsl,

			results: []string{},
		},
		{
			min: 0.09, max: 0.1,
			zsl: zsl,

			results: []string{"0.1: a"},
		},
	}
	cases = cases

}
