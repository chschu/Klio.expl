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
		types.NewIndexRange(types.NewHeadIndex(4), types.NewTailIndex(9)),
		types.NewIndexRange(types.NewPermanentIndex(3), types.NewPermanentIndex(3)),
		types.NewIndexRange(types.NewPermanentIndex(14), types.NewHeadIndex(13)),
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
	assert.Equal(t, types.NewIndexRange(types.NewHeadIndex(42), types.NewHeadIndex(42)), result)
}

func Test_IndexSpecParser_ParseIndexRange_Range(t *testing.T) {
	sut := types.NewIndexSpecParser()

	result, err := sut.ParseIndexRange("  -17:p3 ")

	assert.NoError(t, err)
	assert.Equal(t, types.NewIndexRange(types.NewTailIndex(17), types.NewPermanentIndex(3)), result)
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
	assert.Equal(t, types.NewHeadIndex(4294967295), result)
}

func Test_IndexSpecParser_ParseIndex_TailIndex(t *testing.T) {
	sut := types.NewIndexSpecParser()

	result, err := sut.ParseIndex("  -123 ")

	assert.NoError(t, err)
	assert.Equal(t, types.NewTailIndex(123), result)
}

func Test_IndexSpecParser_ParseIndex_PermanentIndex(t *testing.T) {
	sut := types.NewIndexSpecParser()

	result, err := sut.ParseIndex("  p2 ")

	assert.NoError(t, err)
	assert.Equal(t, types.NewPermanentIndex(2), result)
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
