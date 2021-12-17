package types

type Range struct {
	Min, Max float64

	Minex, Maxex bool
}

func (r *Range) IsValid() bool {
	if r.Min > r.Max {
		return false
	}
	if r.Minex || r.Maxex {
		return r.Max != r.Min
	}
	return true
}

func (r *Range) GteMin(value float64) bool {
	if r.Minex {
		return value > r.Min
	}
	return value >= r.Min
}

func (r *Range) LteMax(value float64) bool {
	if r.Maxex {
		return value < r.Max
	}
	return value <= r.Max
}
