package types_test

import (
	"github.com/stretchr/testify/assert"
	"klio/expl/types"
	"testing"
)

func Test_IndexRange_SqlCondition(t *testing.T) {
	sut := types.NewIndexRange(types.HeadIndex(4102931), types.TailIndex(983142))

	sql, params := sut.SqlCondition()

	assert.Equal(t, "(head_index >= ?) AND (0-tail_index <= 0-?)", sql)
	assert.Equal(t, []any{uint(4102931), uint(983142)}, params)
}

func Test_IndexRange_String_Single(t *testing.T) {
	sut := types.NewIndexRange(types.PermanentIndex(3910293), types.PermanentIndex(3910293))

	assert.Equal(t, "p3910293", sut.String())
}

func Test_IndexRange_String_Range(t *testing.T) {
	sut := types.NewIndexRange(types.HeadIndex(892938744), types.TailIndex(1930192))

	assert.Equal(t, "892938744:-1930192", sut.String())
}
