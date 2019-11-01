package readingtag

import "testing"

func TestIncludesString(t *testing.T) {
	var cases = []struct {
		intention string
		array     []string
		lookup    string
		want      bool
	}{
		{
			"nil",
			nil,
			"",
			false,
		},
		{
			"found",
			[]string{"hello", "world"},
			"WORLD",
			true,
		},
		{
			"not found",
			[]string{"hello", "world"},
			"bob",
			false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			if result := IncludesString(testCase.array, testCase.lookup); result != testCase.want {
				t.Errorf("Includes() = %#v, want %#v", result, testCase.want)
			}
		})
	}
}
