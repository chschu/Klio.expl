package types_test

import (
	"github.com/stretchr/testify/assert"
	"klio/expl/types"
	"testing"
)

func Test_IndexSpec_SqlCondition_Empty(t *testing.T) {
	sut := types.NewIndexSpec()

	sql, params := sut.SQLCondition()

	assert.Equal(t, "false", sql)
	assert.Equal(t, []any{}, params)
}

func Test_IndexSpec_SqlCondition_NonEmpty(t *testing.T) {
	sut := types.NewIndexSpec(
		types.NewIndexRange(types.NewPermanentIndex(17), types.NewTailIndex(3)),
		types.NewIndexRange(types.NewHeadIndex(44), types.NewHeadIndex(45)),
		types.NewIndexRange(types.NewTailIndex(12), types.NewTailIndex(12)),
	)

	sql, params := sut.SQLCondition()

	assert.Equal(t, "false OR ((permanent_index >= ?) AND (tail_index >= ?)) OR ((head_index >= ?) AND (head_index <= ?)) OR (tail_index = ?)", sql)
	assert.Equal(t, []any{uint(17), uint(3), uint(44), uint(45), uint(12)}, params)
}

func Test_IndexSpec_String_Empty(t *testing.T) {
	sut := types.NewIndexSpec()

	assert.Equal(t, "", sut.String())
}

func Test_IndexSpec_String_NonEmpty(t *testing.T) {
	sut := types.NewIndexSpec(
		types.NewIndexRange(types.NewHeadIndex(17), types.NewTailIndex(3)),
		types.NewIndexRange(types.NewPermanentIndex(38), types.NewPermanentIndex(38)),
	)

	assert.Equal(t, "17:-3 p38", sut.String())
}
