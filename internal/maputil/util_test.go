package maputil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys_IntMap(t *testing.T) {
	tests := []struct {
		name string
		m    map[int]string
		want []int
	}{
		{
			name: "nil map",
			m:    nil,
			want: []int{},
		},
		{
			name: "single element",
			m:    map[int]string{42: "foo"},
			want: []int{42},
		},
		{
			name: "multiple elements",
			m:    map[int]string{1: "a", 2: "b", 3: "c"},
			want: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Keys(tt.m)

			assert.ElementsMatch(t, tt.want, got, "Keys(%v)", tt.m)
		})
	}
}
