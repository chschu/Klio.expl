package types_test

import (
	"github.com/stretchr/testify/assert"
	"klio/expl/types"
	"testing"
)

func Test_IndexRange_SQLCondition_Single(t *testing.T) {
	sut := types.NewIndexRange(types.NewPermanentIndex(39766283), types.NewPermanentIndex(39766283))

	sql, params := sut.SQLCondition()

	assert.Equal(t, "permanent_index = ?", sql)
	assert.Equal(t, []any{uint(39766283)}, params)
}

func Test_IndexRange_SQLCondition_Range(t *testing.T) {
	sut := types.NewIndexRange(types.NewHeadIndex(4102931), types.NewTailIndex(983142))

	sql, params := sut.SQLCondition()

	assert.Equal(t, "(head_index >= ?) AND (tail_index >= ?)", sql)
	assert.Equal(t, []any{uint(4102931), uint(983142)}, params)
}

func Test_IndexRange_String_Single(t *testing.T) {
	sut := types.NewIndexRange(types.NewPermanentIndex(3910293), types.NewPermanentIndex(3910293))

	assert.Equal(t, "p3910293", sut.String())
}

func Test_IndexRange_String_Range(t *testing.T) {
	sut := types.NewIndexRange(types.NewHeadIndex(892938744), types.NewTailIndex(1930192))

	assert.Equal(t, "892938744:-1930192", sut.String())
}
