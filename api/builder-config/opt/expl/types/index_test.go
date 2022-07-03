package types_test

import (
	"github.com/stretchr/testify/assert"
	"klio/expl/types"
	"testing"
)

func Test_HeadIndex_SqlCondition_StartingWith(t *testing.T) {
	sut := types.HeadIndex(1923781)

	sql, params := sut.SqlCondition(types.IndexStartingWith)

	assert.Equal(t, "head_index >= ?", sql)
	assert.Equal(t, []any{uint(1923781)}, params)
}

func Test_HeadIndex_SqlCondition_EndingWith(t *testing.T) {
	sut := types.HeadIndex(83817281)

	sql, params := sut.SqlCondition(types.IndexEndingWith)

	assert.Equal(t, "head_index <= ?", sql)
	assert.Equal(t, []any{uint(83817281)}, params)
}

func Test_HeadIndex_String(t *testing.T) {
	sut := types.HeadIndex(481931)

	assert.Equal(t, "481931", sut.String())
}

func Test_TailIndex_SqlCondition_StartingWith(t *testing.T) {
	sut := types.TailIndex(8919283)

	sql, params := sut.SqlCondition(types.IndexStartingWith)

	assert.Equal(t, "0-tail_index >= 0-?", sql)
	assert.Equal(t, []any{uint(8919283)}, params)
}

func Test_TailIndex_SqlCondition_EndingWith(t *testing.T) {
	sut := types.TailIndex(38389912)

	sql, params := sut.SqlCondition(types.IndexEndingWith)

	assert.Equal(t, "0-tail_index <= 0-?", sql)
	assert.Equal(t, []any{uint(38389912)}, params)
}

func Test_TailIndex_String(t *testing.T) {
	sut := types.TailIndex(1939124)

	assert.Equal(t, "-1939124", sut.String())
}

func Test_PermanentIndex_SqlCondition_StartingWith(t *testing.T) {
	sut := types.PermanentIndex(2799171571)

	sql, params := sut.SqlCondition(types.IndexStartingWith)

	assert.Equal(t, "permanent_index >= ?", sql)
	assert.Equal(t, []any{uint(2799171571)}, params)
}

func Test_PermanentIndex_SqlCondition_EndingWith(t *testing.T) {
	sut := types.PermanentIndex(2284892642)

	sql, params := sut.SqlCondition(types.IndexEndingWith)

	assert.Equal(t, "permanent_index <= ?", sql)
	assert.Equal(t, []any{uint(2284892642)}, params)
}

func Test_PermanentIndex_String(t *testing.T) {
	sut := types.PermanentIndex(7189482)

	assert.Equal(t, "p7189482", sut.String())
}
