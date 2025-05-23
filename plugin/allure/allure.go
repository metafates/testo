package allure

import (
	"testman"
)

type Allure struct {
	*testman.T

	title string
}

func (*Allure) New(t *testman.T) *Allure { return &Allure{T: t} }

func (a *Allure) Title(title string) {
	a.title = title
}

func (a *Allure) AfterAll() {
	a.Log("Allure.AfterAll")
}

func (a *Allure) BeforeAll() {
	a.Log("Allure.BeforeAll")
}

func (a *Allure) BeforeEach() {
	a.Log("Allure.BeforeEach " + a.Name())
}

func (a *Allure) AfterEach() {
	a.Log("Allure.AfterEach " + a.Name())
}
