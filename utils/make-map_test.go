package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/witwoywhy/go-cores/utils"
)

func TestMakeMap(t *testing.T) {
	type Data struct {
		Id   int64
		Name string
	}

	type testCase struct {
		given  []Data
		expect any
	}

	given := []Data{
		{
			Id:   1,
			Name: "A B",
		},
		{
			Id:   2,
			Name: "C D",
		},
	}

	t.Run("key=Id, val=Name", func(t *testing.T) {
		tc := testCase{
			given: given,
			expect: map[int64]string{
				1: "A B",
				2: "C D",
			},
		}

		got := utils.MakeMap[int64, string](tc.given, "Id", "Name")
		assert.Equal(t, tc.expect, got)
	})

	t.Run("key=Name, val=Id", func(t *testing.T) {
		tc := testCase{
			given: given,
			expect: map[string]int64{
				"A B": 1,
				"C D": 2,
			},
		}

		got := utils.MakeMap[string, int64](tc.given, "Name", "Id")
		assert.Equal(t, tc.expect, got)
	})

	t.Run("key=Id, val=true", func(t *testing.T) {
		tc := testCase{
			given: given,
			expect: map[int64]bool{
				1: true,
				2: true,
			},
		}

		got := utils.MakeMap[int64, bool](tc.given, "Id", "true")
		assert.Equal(t, tc.expect, got)
	})

	t.Run("key=Id, val=struct", func(t *testing.T) {
		tc := testCase{
			given: given,
			expect: map[int64]Data{
				1: given[0],
				2: given[1],
			},
		}

		got := utils.MakeMap[int64, Data](tc.given, "Id", "struct")
		assert.Equal(t, tc.expect, got)
	})
}
