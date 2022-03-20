package async_test

import (
	"errors"
	"testing"

	"github.com/nilsbu/async"
)

// TODO test the rest

func TestPe(t *testing.T) {
	data := make([]int, 100)

	for _, c := range []struct {
		name string
		fs   []func() error
		data []int // when smaller than 100, remainder is assumed to be zero
		ok   bool
	}{
		{
			"empty",
			[]func() error{},
			[]int{},
			true,
		},
		{
			"no error",
			[]func() error{
				func() error { data[0] = 1; return nil },
				func() error { data[1] = 1; return nil },
				func() error { data[2] = 1; return nil },
			},
			[]int{1, 1, 1},
			true,
		},
		{
			"with error",
			[]func() error{
				func() error { data[0] = 1; return nil },
				func() error { data[1] = 1; return errors.New("") },
				func() error { data[2] = 1; return nil },
			},
			[]int{},
			false,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// reset
			for i := range data {
				data[i] = 0
			}

			err := async.Pe(c.fs)
			if c.ok && err != nil {
				t.Fatal("expected no error but got", err)
			} else if !c.ok && err == nil {
				t.Fatal("expected error but no occurred")
			}

			if err == nil {
				for i, d := range data {
					var v int
					if i < len(c.data) {
						v = c.data[i]
					}
					if d != v {
						t.Errorf("@%v: expected %v but got %v", i, v, d)
					}
				}
			}
		})
	}
}
