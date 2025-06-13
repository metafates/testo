package reflectutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeValue(t *testing.T) {
	type Mock struct {
		String *string
		Nested *struct {
			Nested *struct {
				N *int
			}
		}
	}

	value := Filled[Mock]()

	require.NotNil(t, value.Nested)
	require.NotNil(t, value.String)
	require.NotNil(t, value.Nested.Nested)
	require.NotNil(t, value.Nested.Nested.N)
}

func TestDeepClone(t *testing.T) {
	type Mock struct {
		Name    string
		private []string
	}

	original := Mock{Name: "Test", private: []string{"test"}}

	clone := DeepClone(original)

	require.Equal(t, original, clone)
}
