package types_test

import (
	"github.com/stretchr/testify/assert"
	"klio/expl/types"
	"testing"
)

func Test_IndexSpecParser_ParseIndexSpec_Empty(t *testing.T) {
	sut := types.NewIndexSpecParser()

	result, err := sut.ParseIndexSpec("   ")

	assert.NoError(t, err)
	assert.Equal(t, types.NewIndexSpec(), result)
}

func Test_IndexSpecParser_ParseIndexSpec_NonEmpty(t *testing.T) {
	sut := types.NewIndexSpecParser()

	result, err := sut.ParseIndexSpec(" 4:-9  p3   p14:13 ")

	expectedIndexSpec := types.NewIndexSpec(
		types.NewIndexRange(types.HeadIndex(4), types.TailIndex(9)),
		types.NewIndexRange(types.PermanentIndex(3), types.PermanentIndex(3)),
		types.NewIndexRange(types.PermanentIndex(14), types.HeadIndex(13)),
	)

	assert.NoError(t, err)
	assert.Equal(t, expectedIndexSpec, result)
}

func Test_IndexSpecParser_ParseIndexSpec_Invalod(t *testing.T) {
	sut := types.NewIndexSpecParser()

	_, err := sut.ParseIndexSpec("x")

	assert.Error(t, err)
}

func Test_IndexSpecParser_ParseIndexRange_Single(t *testing.T) {
	sut := types.NewIndexSpecParser()

	result, err := sut.ParseIndexRange("   42  ")

	assert.NoError(t, err)
	assert.Equal(t, types.NewIndexRange(types.HeadIndex(42), types.HeadIndex(42)), result)
}

func Test_IndexSpecParser_ParseIndexRange_Range(t *testing.T) {
	sut := types.NewIndexSpecParser()

	result, err := sut.ParseIndexRange("  -17:p3 ")

	assert.NoError(t, err)
	assert.Equal(t, types.NewIndexRange(types.TailIndex(17), types.PermanentIndex(3)), result)
}

func Test_IndexSpecParser_ParseIndexRange_Invalid(t *testing.T) {
	sut := types.NewIndexSpecParser()

	_, err := sut.ParseIndexRange("1: 2")

	assert.Error(t, err)
}

func Test_IndexSpecParser_ParseIndexRange_FromInvalid(t *testing.T) {
	sut := types.NewIndexSpecParser()

	_, err := sut.ParseIndexRange("a:2")

	assert.Error(t, err)
}

func Test_IndexSpecParser_ParseIndexRange_ToInvalid(t *testing.T) {
	sut := types.NewIndexSpecParser()

	_, err := sut.ParseIndexRange(" 1:b ")

	assert.Error(t, err)
}

func Test_IndexSpecParser_ParseIndex_HeadIndex(t *testing.T) {
	sut := types.NewIndexSpecParser()

	result, err := sut.ParseIndex("  4294967295 ")

	assert.NoError(t, err)
	assert.Equal(t, types.HeadIndex(4294967295), result)
}

func Test_IndexSpecParser_ParseIndex_TailIndex(t *testing.T) {
	sut := types.NewIndexSpecParser()

	result, err := sut.ParseIndex("  -123 ")

	assert.NoError(t, err)
	assert.Equal(t, types.TailIndex(123), result)
}

func Test_IndexSpecParser_ParseIndex_PermanentIndex(t *testing.T) {
	sut := types.NewIndexSpecParser()

	result, err := sut.ParseIndex("  p2 ")

	assert.NoError(t, err)
	assert.Equal(t, types.PermanentIndex(2), result)
}

func Test_IndexSpecParser_ParseIndex_Invalid(t *testing.T) {
	sut := types.NewIndexSpecParser()

	_, err := sut.ParseIndex("x")

	assert.Error(t, err)
}

func Test_IndexSpecParser_ParseIndex_TooLarge(t *testing.T) {
	sut := types.NewIndexSpecParser()

	_, err := sut.ParseIndex("4294967296")

	assert.Error(t, err)
}
