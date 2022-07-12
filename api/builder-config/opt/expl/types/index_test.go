package types_test

import (
	"github.com/chschu/Klio.expl/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewIndexFromString_HeadIndex(t *testing.T) {
	result, err := types.NewIndexFromString("  4294967295 ")

	assert.NoError(t, err)
	assert.Equal(t, types.NewHeadIndex(4294967295), result)
}

func Test_NewIndexFromString_TailIndex(t *testing.T) {
	result, err := types.NewIndexFromString("  -123 ")

	assert.NoError(t, err)
	assert.Equal(t, types.NewTailIndex(123), result)
}

func Test_NewIndexFromString_PermanentIndex(t *testing.T) {
	result, err := types.NewIndexFromString("  p2 ")

	assert.NoError(t, err)
	assert.Equal(t, types.NewPermanentIndex(2), result)
}

func Test_NewIndexFromString_Invalid(t *testing.T) {
	_, err := types.NewIndexFromString("x")

	assert.Error(t, err)
}

func Test_NewIndexFromString_TooLarge(t *testing.T) {
	_, err := types.NewIndexFromString("4294967296")

	assert.Error(t, err)
}

func Test_HeadIndex_SQLConditionStartingWith(t *testing.T) {
	sut := types.NewHeadIndex(1923781)

	sql, params := sut.SQLConditionStartingWith()

	assert.Equal(t, "head_index >= ?", sql)
	assert.Equal(t, []any{uint(1923781)}, params)
}

func Test_HeadIndex_SQLConditionMatching(t *testing.T) {
	sut := types.NewHeadIndex(818664812)

	sql, params := sut.SQLConditionMatching()

	assert.Equal(t, "head_index = ?", sql)
	assert.Equal(t, []any{uint(818664812)}, params)
}

func Test_HeadIndex_SQLConditionEndingWith(t *testing.T) {
	sut := types.NewHeadIndex(83817281)

	sql, params := sut.SQLConditionEndingWith()

	assert.Equal(t, "head_index <= ?", sql)
	assert.Equal(t, []any{uint(83817281)}, params)
}

func Test_HeadIndex_String(t *testing.T) {
	sut := types.NewHeadIndex(481931)

	assert.Equal(t, "481931", sut.String())
}

func Test_TailIndex_SQLConditionStartingWith(t *testing.T) {
	sut := types.NewTailIndex(8919283)

	sql, params := sut.SQLConditionStartingWith()

	assert.Equal(t, "tail_index <= ?", sql)
	assert.Equal(t, []any{uint(8919283)}, params)
}

func Test_TailIndex_SQLConditionMatching(t *testing.T) {
	sut := types.NewTailIndex(1938782663)

	sql, params := sut.SQLConditionMatching()

	assert.Equal(t, "tail_index = ?", sql)
	assert.Equal(t, []any{uint(1938782663)}, params)
}

func Test_TailIndex_SQLConditionEndingWith(t *testing.T) {
	sut := types.NewTailIndex(38389912)

	sql, params := sut.SQLConditionEndingWith()

	assert.Equal(t, "tail_index >= ?", sql)
	assert.Equal(t, []any{uint(38389912)}, params)
}

func Test_TailIndex_String(t *testing.T) {
	sut := types.NewTailIndex(1939124)

	assert.Equal(t, "-1939124", sut.String())
}

func Test_PermanentIndex_SQLConditionStartingWith(t *testing.T) {
	sut := types.NewPermanentIndex(2799171571)

	sql, params := sut.SQLConditionStartingWith()

	assert.Equal(t, "permanent_index >= ?", sql)
	assert.Equal(t, []any{uint(2799171571)}, params)
}

func Test_PermanentIndex_SQLConditionMatching(t *testing.T) {
	sut := types.NewPermanentIndex(168947826)

	sql, params := sut.SQLConditionMatching()

	assert.Equal(t, "permanent_index = ?", sql)
	assert.Equal(t, []any{uint(168947826)}, params)
}

func Test_PermanentIndex_SQLCondition_EndingWith(t *testing.T) {
	sut := types.NewPermanentIndex(2284892642)

	sql, params := sut.SQLConditionEndingWith()

	assert.Equal(t, "permanent_index <= ?", sql)
	assert.Equal(t, []any{uint(2284892642)}, params)
}

func Test_PermanentIndex_String(t *testing.T) {
	sut := types.NewPermanentIndex(7189482)

	assert.Equal(t, "p7189482", sut.String())
}
