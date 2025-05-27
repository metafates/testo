package assertions

import (
	"github.com/metafates/tego"
)

type Assertions struct {
	*tego.T
}

type Common struct {
	errorf func(msg string, args ...any)
}

func (a *Assertions) Require() Common {
	return Common{
		errorf: a.Fatalf,
	}
}

func (a *Assertions) Assert() Common {
	return Common{
		errorf: a.Errorf,
	}
}

func (c Common) True(b bool) {
	if !b {
		c.errorf("want true, got false")
	}
}
