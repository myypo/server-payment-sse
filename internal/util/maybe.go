package u

import (
	"bytes"
	"encoding/json"
)

type Maybe[T any] map[bool]T

func MaybeFrom[T any](v T) Maybe[T] {
	return map[bool]T{true: v}
}

func (m Maybe[T]) MarshalJSON() ([]byte, error) {
	v, ok := m.Some()
	if !ok {
		return []byte("null"), nil
	}
	return json.Marshal(v)
}

func (m Maybe[T]) Some() (T, bool) {
	v, ok := m[true]
	return v, ok
}

func (m Maybe[T]) None() bool {
	if m == nil {
		return true
	}
	_, ok := m[false]
	return ok
}

func (m *Maybe[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		var def T
		*m = map[bool]T{false: def}
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*m = map[bool]T{true: v}

	return nil
}

func (m Maybe[T]) IsSome() bool {
	_, ok := m[true]
	return ok
}

func (m Maybe[T]) tryResolve() (T, bool) {
	return m.Some()
}

type maybeResolvable[T any] interface {
	tryResolve() (T, bool)
}

func OR[T any, M maybeResolvable[T]](m M, def T) T {
	if v, ok := m.tryResolve(); ok {
		return v
	}

	return def
}
