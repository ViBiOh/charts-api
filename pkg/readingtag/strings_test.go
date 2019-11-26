package readingtag

import "testing"

func TestIncludesUint64(t *testing.T) {
	var cases = []struct {
		intention string
		array     []uint64
		lookup    uint64
		want      bool
	}{
		{
			"nil",
			nil,
			0,
			false,
		},
		{
			"found",
			[]uint64{8000, 6000},
			8000,
			true,
		},
		{
			"not found",
			[]uint64{8000, 6000},
			7000,
			false,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			if result := IncludesUint64(testCase.array, testCase.lookup); result != testCase.want {
				t.Errorf("Includes() = %#v, want %#v", result, testCase.want)
			}
		})
	}
}
