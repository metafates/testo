package testman

type Plugin any

type TP[P Plugin] interface {
	T

	P() P
	RunP(name string, f func(t TP[P]))
}

type tp[P Plugin] struct {
	T

	p P
}

func (tt *tp[P]) P() P { return tt.p }

func (tt *tp[P]) RunP(name string, f func(t TP[P])) {
	tt.Run(name, func(t T) {
		f(&tp[P]{T: t, p: tt.p})
	})
}
