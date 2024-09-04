package u

func All[I any](s []I, f func(I) bool) bool {
	for _, v := range s {
		if !f(v) {
			return false
		}
	}
	return true
}

func FilterI[I any](s []I, f func(int, I) bool) []I {
	res := make([]I, 0, len(s))
	for i, v := range s {
		if f(i, v) {
			res = append(res, v)
		}
	}
	return res
}

func Fold[T any, R any](s []T, f func(acc R, curr T) R, neutral R) R {
	acc := neutral
	for _, elem := range s {
		acc = f(acc, elem)
	}
	return acc
}

func Map[I, O any](s []I, f func(I) O) []O {
	res := make([]O, 0, len(s))
	for _, v := range s {
		res = append(res, f(v))
	}
	return res
}

func MapI[I, O any](s []I, f func(int, I) O) []O {
	res := make([]O, 0, len(s))
	for i, v := range s {
		res = append(res, f(i, v))
	}
	return res
}

func MapE[I, O any](s []I, f func(I) (O, error)) ([]O, error) {
	res := make([]O, 0, len(s))
	for _, v := range s {
		r, err := f(v)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func Permut[T any](s []T) [][]T {
	if len(s) == 0 {
		return [][]T{}
	}

	if len(s) == 1 {
		return [][]T{{s[0]}}
	}

	res := [][]T{}

	for i, v := range s {
		remaining := make([]T, len(s)-1)
		copy(remaining[:i], s[:i])
		copy(remaining[i:], s[i+1:])

		sub := Permut(remaining)

		for _, p := range sub {
			res = append(res, append([]T{v}, p...))
		}
	}

	return res
}
