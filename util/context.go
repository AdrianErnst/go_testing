package util

type ContextKey string

func (s ContextKey) String() string {
	return string(s)
}
