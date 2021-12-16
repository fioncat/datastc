package types

type ScoreKey struct {
	Key   string
	Score float64
}

type ScoreValue struct {
	Value interface{}
	Score float64
}

func ScoreValue2Key(values []ScoreValue) []ScoreKey {
	keys := make([]ScoreKey, len(values))
	for i, val := range values {
		keys[i] = ScoreKey{
			Key:   val.Value.(string),
			Score: val.Score,
		}
	}
	return keys
}
