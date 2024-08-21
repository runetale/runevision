package utility

type HandleSet[T any] map[Handle]T

type Handle struct {
	v *byte
}

func (s *HandleSet[T]) Add(e T) Handle {
	h := Handle{new(byte)}
	if *s == nil {
		*s = make(HandleSet[T])
	}
	(*s)[h] = e
	return h
}
