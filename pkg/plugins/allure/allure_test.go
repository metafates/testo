package allure

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/metafates/testo"
	"github.com/stretchr/testify/require"
)

type T = *struct {
	*testo.T

	Allure
}

type Suite struct{}

func TestAllure(t *testing.T) {
	// TODO: use t.Chdir() when go is updated
	wd, err := os.Getwd()
	require.NoError(t, err)

	dir := t.TempDir()
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(wd) })
	t.Cleanup(func() { os.RemoveAll(filepath.Join(dir, "allure-results")) })

	testo.RunSuite[*Suite, T](t)

	require.DirExists(t, filepath.Join(dir, "allure-results"), "output dir does not exist")

	// TODO: other assertions
}

func (Suite) BeforeEach(t T) {
	Setup(t, "init", func(t T) {
		testo.Run(t, "nested", func(t T) {})
	})

	Setup(t, "extra init", func(t T) {})
}

func (Suite) AfterEach(t T) {
	Setup(t, "destroy", func(t T) {
		testo.Run(t, "nested", func(t T) {})
	})

	Setup(t, "extra destroy", func(t T) {})
}

func (Suite) TestFoo(t T) {
	t.Flaky()

	testo.Run(t, "subtest", func(t T) {
		t.Known()

		testo.Run(t, "nested", func(t T) {})
	})
}

func (Suite) CasesX() []int {
	return []int{1}
}

func (Suite) TestBar(t T, params struct{ X int }) {
}
