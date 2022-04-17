package schema

import "math"

type NumberSchema struct {
	Precision int     `json:"precision"`
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
}

// Verify that `NumberSchema` implements the `Schema` interface
var _ Schema = NumberSchema{}

func (schema NumberSchema) Validate(value any) bool {
	_value, ok := value.(float64)
	if !ok {
		return false
	}

	if _value < schema.Min || _value > schema.Max {
		return false
	}

	places := math.Pow10(schema.Precision)
	if math.Round(_value*places)/places != _value {
		return false
	}

	return true
}
