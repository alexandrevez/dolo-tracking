package hash

import "testing"

func TestMD5String(t *testing.T) {
	testList := []struct {
		input          []string
		expectedOutput string
	}{
		{
			input:          []string{"", "", "", "", ""},
			expectedOutput: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			input:          []string{"d41d8cd98f00b204e9800998ecf8427e"},
			expectedOutput: "74be16979710d4c4e7c6647856088456",
		},
		{
			input:          []string{},
			expectedOutput: "d41d8cd98f00b204e9800998ecf8427e",
		},
	}

	for _, tt := range testList {
		t.Run("", func(t *testing.T) {
			result := MD5String(tt.input...)
			if result != tt.expectedOutput {
				t.Errorf("Expected hash %s but got %s", tt.expectedOutput, result)
			}
		})
	}
}

func TestSHA256String(t *testing.T) {
	testList := []struct {
		input          []string
		expectedOutput string
	}{
		{
			input:          []string{"", "", "", "", ""},
			expectedOutput: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			input:          []string{"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
			expectedOutput: "cd372fb85148700fa88095e3492d3f9f5beb43e555e5ff26d95f5a6adc36f8e6",
		},
		{
			input:          []string{},
			expectedOutput: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}

	for _, tt := range testList {
		t.Run("", func(t *testing.T) {
			result := Sha256String(tt.input...)
			if result != tt.expectedOutput {
				t.Errorf("Expected hash %s but got %s", tt.expectedOutput, result)
			}
		})
	}
}
