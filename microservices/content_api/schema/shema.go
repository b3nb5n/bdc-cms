package schema

type Schema interface {
	Validate(any) bool
}