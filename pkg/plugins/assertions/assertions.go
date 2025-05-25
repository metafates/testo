package assertions

import (
	"testman"
)

type Assertions struct {
	*testman.T
}

type common struct {
	errorf func(msg string, args ...any)
}

func (a *Assertions) Require() common {
	return common{
		errorf: a.Fatalf,
	}
}

func (a *Assertions) Assert() common {
	return common{
		errorf: a.Errorf,
	}
}

func (c common) True(b bool) {
	if !b {
		c.errorf("want true, got false")
	}
}
