package types

import (
	"bytes"
)

type Comparable interface {
	Compare(b interface{}) int
}

func Compare(a, b interface{}) int {
	switch va := a.(type) {
	case string:
		vb := b.(string)
		if va == vb {
			return 0
		}
		if va < vb {
			return -1
		}
		return 1

	case []byte:
		return bytes.Compare(va, b.([]byte))

	default:
		ca, ok := a.(Comparable)
		if !ok {
			panic("Compare: a must implement types.Comparable interface")
		}
		return ca.Compare(b)
	}
}
