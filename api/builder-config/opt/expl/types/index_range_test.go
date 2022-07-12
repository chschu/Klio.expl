package types_test

import (
	"github.com/chschu/Klio.expl/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewIndexRangeFromString_Single(t *testing.T) {
	result, err := types.NewIndexRangeFromString("   42  ")

	assert.NoError(t, err)
	assert.Equal(t, types.NewIndexRange(types.NewHeadIndex(42), types.NewHeadIndex(42)), result)
}

func Test_NewIndexRangeFromString_Range(t *testing.T) {
	result, err := types.NewIndexRangeFromString("  -17:p3 ")

	assert.NoError(t, err)
	assert.Equal(t, types.NewIndexRange(types.NewTailIndex(17), types.NewPermanentIndex(3)), result)
}

func Test_NewIndexRangeFromString_Invalid(t *testing.T) {
	_, err := types.NewIndexRangeFromString("1: 2")

	assert.Error(t, err)
}

func Test_NewIndexRangeFromString_FromInvalid(t *testing.T) {
	_, err := types.NewIndexRangeFromString("a:2")

	assert.Error(t, err)
}

func Test_NewIndexRangeFromString_ToInvalid(t *testing.T) {
	_, err := types.NewIndexRangeFromString(" 1:b ")

	assert.Error(t, err)
}

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
