package sortedset

import (
	"fmt"
	"testing"

	"github.com/fioncat/datastc/pkg/types"
)

func TestSet(t *testing.T) {
	ss := New(5)

	ss.Set(0.78, "c")
	ss.Set(0.22, "c++")
	ss.Set(0.98, "python")
	ss.Set(0.87, "c#")
	ss.Set(0.12, "HTML")
	ss.Set(0.67, "javascript")
	ss.Set(0.24, "java")
	ss.Set(0.87, "go")

	ss.Delete("HTML")

	ss.Scan(func(score float64, key string) bool {
		fmt.Printf("%.2f: %s\n", score, key)
		return true
	})
	fmt.Println("=========================")
	ss.ScanRange(types.Range{
		Min: 0.3,
		Max: 0.9,
	}, func(score float64, key string) bool {
		fmt.Printf("%.2f: %s\n", score, key)
		return true
	})

}
