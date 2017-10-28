package format

import "testing"

func TestNewJSONString(t *testing.T) {
	testList := []struct {
		o             interface{}
		expectedReply string
	}{
		{
			o:             32,
			expectedReply: "32",
		},
		{
			o:             make(chan bool),
			expectedReply: "",
		},
		{
			o:             func() {},
			expectedReply: "",
		},
		{
			o:             3 + 4i,
			expectedReply: "",
		},
		{
			o: struct {
				A string `json:"a"`
				B string `json:"-"`
				C int    `json:"cde"`
			}{
				A: "A",
				B: "HIDDEN",
				C: 32,
			},
			expectedReply: `{
	"a": "A",
	"cde": 32
}`,
		},
		{
			o: struct {
				A string `json:"a"`
				B struct {
					Simon string
				} `json:"BLetter"`
				C int `json:"cde"`
			}{
				A: "A",
				B: struct {
					Simon string
				}{Simon: "Y$"},
				C: 32,
			},
			expectedReply: `{
	"a": "A",
	"BLetter": {
		"Simon": "Y$"
	},
	"cde": 32
}`,
		},
	}

	for _, tt := range testList {
		t.Run("", func(t *testing.T) {
			result := NewJSONString(tt.o)
			if tt.expectedReply != result {
				t.Errorf("Expected: \n%s \ngot \n%s", tt.expectedReply, result)
			}
		})
	}
}
