package s3_client

import (
	"reflect"
	"testing"
)

func TestSplitKeyByDepth(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		depth    int
		expected []string
	}{
		{
			name:     "Normal case with valid depth",
			key:      "a/b/c/d",
			depth:    2,
			expected: []string{"a", "b"},
		},
		{
			name:     "Depth is zero, return full split",
			key:      "a/b/c/d",
			depth:    0,
			expected: []string{"a", "b", "c", "d"},
		},
		{
			name:     "Depth greater than length of split, return full split",
			key:      "a/b",
			depth:    5,
			expected: []string{"a", "b"},
		},
		{
			name:     "Depth equal to length of split, return full split",
			key:      "a/b/c",
			depth:    3,
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Empty key string",
			key:      "",
			depth:    2,
			expected: []string{""},
		},
		{
			name:     "Key with trailing slash",
			key:      "a/b/c/",
			depth:    2,
			expected: []string{"a", "b"},
		},
		{
			name:     "Key with leading slash",
			key:      "/a/b/c",
			depth:    2,
			expected: []string{"", "a"},
		},
		{
			name:     "Key with multiple slashes",
			key:      "a//b/c",
			depth:    3,
			expected: []string{"a", "", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitKeyByDepth(tt.key, tt.depth)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
