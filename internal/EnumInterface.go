package internal

type enum[T int | string] interface {
	IsMember(val T) bool
	Members() []T
}
