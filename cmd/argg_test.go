package cmd

import (
	"bytes"
	"math"
	"testing"

	"github.com/matryer/is"
)

func TestReadPipedArgs(t *testing.T) {
	is := is.New(t)
	tt := []struct {
		given    string
		expected []string
	}{
		{
			given:    "single",
			expected: []string{"single"},
		},
		{
			given:    "single\n",
			expected: []string{"single"},
		},
		{
			given:    "one\ndouble\n",
			expected: []string{"one", "double"},
		},
	}

	for _, tc := range tt {
		mockStdin := bytes.NewBufferString(tc.given)
		received := readPipedArgs(mockStdin)

		is.Equal(received, tc.expected)
	}
}

func TestMergeArgs(t *testing.T) {
	is := is.New(t)
	tt := []struct {
		a        []string
		b        []string
		expected []string
	}{
		{
			a:        []string{"ls", "-a"},
			b:        []string{"dir1", "dir2"},
			expected: []string{"-a", "dir1", "dir2"},
		},
		{
			a:        []string{"cat"},
			b:        []string{"dir1", "dir2"},
			expected: []string{"dir1", "dir2"},
		},
	}

	for _, tc := range tt {
		received := mergeArgs(tc.a, tc.b)
		is.Equal(received, tc.expected)
	}
}

func TestTakeN(t *testing.T) {
	is := is.New(t)
	tt := []struct {
		given    []string
		n        int
		expected [][]string
	}{
		{
			given:    []string{"arg1", "arg2", "arg3"},
			n:        1,
			expected: [][]string{{"arg1"}, {"arg2"}, {"arg3"}},
		},
		{
			given:    []string{"arg1", "arg2", "arg3"},
			n:        2,
			expected: [][]string{{"arg1", "arg2"}, {"arg3"}},
		},
		{
			given:    []string{"arg1", "arg2", "arg3", "arg4", "arg5"},
			n:        2,
			expected: [][]string{{"arg1", "arg2"}, {"arg3", "arg4"}, {"arg5"}},
		},
		{
			given:    []string{"arg1", "arg2", "arg3", "arg4", "arg5"},
			n:        3,
			expected: [][]string{{"arg1", "arg2", "arg3"}, {"arg4", "arg5"}},
		},
		{
			given:    []string{"arg1", "arg2", "arg3"},
			n:        math.MaxInt,
			expected: [][]string{{"arg1", "arg2", "arg3"}},
		},
	}

	for _, tc := range tt {
		received := splitArgsByN(tc.given, tc.n)
		is.Equal(received, tc.expected)
	}
}
