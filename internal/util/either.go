package u

type Either[L, R any] struct {
	lhs  L
	rhs  R
	side side
}

type side bool

const (
	Left  side = true
	Right side = false
)

func EitherFromLeft[L, R any](v L) *Either[L, R] {
	return &Either[L, R]{lhs: v, side: Left}
}

func EitherFromRight[L, R any](v R) *Either[L, R] {
	return &Either[L, R]{rhs: v, side: Right}
}

func (e Either[L, R]) Left() (L, bool) {
	if e.side == Left {
		return e.lhs, true
	}

	return e.lhs, false
}

func (e Either[L, R]) Right() (R, bool) {
	if e.side == Right {
		return e.rhs, true
	}

	return e.rhs, false
}
