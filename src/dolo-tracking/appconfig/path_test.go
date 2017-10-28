package appconfig

import (
	"testing"
)

func TestBuildPath(t *testing.T) {
	testList := []struct {
		isFile        bool
		parts         []string
		expectedReply string
	}{
		{
			isFile:        false,
			parts:         []string{},
			expectedReply: "",
		},
		{
			isFile:        false,
			parts:         []string{""},
			expectedReply: "",
		},
		{
			isFile:        false,
			parts:         []string{"a"},
			expectedReply: "a/",
		},
		{
			isFile:        true,
			parts:         []string{"a"},
			expectedReply: "a",
		},
		{
			isFile:        false,
			parts:         []string{"a", "b"},
			expectedReply: "a/b/",
		},
		{
			isFile:        true,
			parts:         []string{"a", "b"},
			expectedReply: "a/b",
		},
		{
			isFile:        false,
			parts:         []string{"a/", "b"},
			expectedReply: "a/b/",
		},
		{
			isFile:        true,
			parts:         []string{"a/", "b/"},
			expectedReply: "a/b",
		},
		{
			isFile:        true,
			parts:         []string{"/a/", "b/", "", "", ""},
			expectedReply: "/a/b",
		},
	}

	for _, tt := range testList {
		t.Run("", func(t *testing.T) {
			result := BuildPath(tt.isFile, tt.parts...)
			if tt.expectedReply != result {
				t.Errorf("Expected: \n%s \ngot \n%s", tt.expectedReply, result)
			}
		})
	}
}
