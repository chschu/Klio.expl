package webhook_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"klio/expl/generated/types_mocks"
	"klio/expl/generated/webhook_mocks"
	"klio/expl/types"
	"klio/expl/webhook"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_ExplHandler_Expl_Success_WithIndexSpec(t *testing.T) {
	ctrl := gomock.NewController(t)
	explainerMock := webhook_mocks.NewMockLimitedExplainer(ctrl)
	indexSpecParserMock := types_mocks.NewMockIndexSpecParser(ctrl)
	entryListStringerMock := webhook_mocks.NewMockEntryListStringer(ctrl)

	maxImmediateResults := 3919283
	webUrlPathPrefix := "/web-expl-prefix/"
	webUrlValidity := 291293 * time.Second

	sut := webhook.NewExplHandler(explainerMock, indexSpecParserMock, entryListStringerMock, maxImmediateResults, webUrlPathPrefix, webUrlValidity)

	in := &webhook.Request{
		Text: "!expl test-key index-spec-string  ",
	}

	ctx := context.WithValue(context.Background(), "dummy", "dummy")
	req := httptest.NewRequest("DUMMY", "/dummy", nil).WithContext(ctx)

	now := time.Now()

	entries := []types.Entry{{Id: rand.Int31()}, {Id: rand.Int31()}}
	total := 8492193
	resultText := "result text"

	indexSpecMock := types_mocks.NewMockIndexSpec(ctrl)

	indexSpecParserMock.EXPECT().ParseIndexSpec("index-spec-string").Return(indexSpecMock, nil)
	explainerMock.EXPECT().ExplainWithLimit(ctx, "test-key", indexSpecMock, maxImmediateResults).Return(entries, total, nil)
	entryListStringerMock.EXPECT().String(entries, total, "test-key", webUrlPathPrefix, now.Add(webUrlValidity), req).Return(resultText)

	out, err := sut.Handle(in, req, now)
	assert.NoError(t, err)
	assert.Equal(t, resultText, out.Text)
}

func Test_ExplHandler_Expl_Success_WithoutIndexSpec(t *testing.T) {
	ctrl := gomock.NewController(t)
	explainerMock := webhook_mocks.NewMockLimitedExplainer(ctrl)
	indexSpecParserMock := types_mocks.NewMockIndexSpecParser(ctrl)
	entryListStringerMock := webhook_mocks.NewMockEntryListStringer(ctrl)

	maxImmediateResults := 3919283
	webUrlPathPrefix := "/web-expl-prefix/"
	webUrlValidity := 291293 * time.Second

	sut := webhook.NewExplHandler(explainerMock, indexSpecParserMock, entryListStringerMock, maxImmediateResults, webUrlPathPrefix, webUrlValidity)

	in := &webhook.Request{
		Text: "!expl test-key ",
	}

	ctx := context.WithValue(context.Background(), "dummy", "dummy")
	req := httptest.NewRequest("DUMMY", "/dummy", nil).WithContext(ctx)

	now := time.Now()

	entries := []types.Entry{{Id: rand.Int31()}, {Id: rand.Int31()}}
	total := 312934
	resultText := "another result text"

	explainerMock.EXPECT().ExplainWithLimit(ctx, "test-key", types.IndexSpecAll(), maxImmediateResults).Return(entries, total, nil)
	entryListStringerMock.EXPECT().String(entries, total, "test-key", webUrlPathPrefix, now.Add(webUrlValidity), req).Return(resultText)

	out, err := sut.Handle(in, req, now)
	assert.NoError(t, err)
	assert.Equal(t, resultText, out.Text)
}

func Test_ExplHandler_Expl_SoftFail_InvalidSyntax(t *testing.T) {
	ctrl := gomock.NewController(t)
	explainerMock := webhook_mocks.NewMockLimitedExplainer(ctrl)
	indexSpecParserMock := types_mocks.NewMockIndexSpecParser(ctrl)
	entryListStringerMock := webhook_mocks.NewMockEntryListStringer(ctrl)

	sut := webhook.NewExplHandler(explainerMock, indexSpecParserMock, entryListStringerMock, 50, "/dummy/", time.Minute)

	in := &webhook.Request{
		Text:        "!expl ",
		TriggerWord: "!triggy",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	out, err := sut.Handle(in, req, time.Now())

	assert.NoError(t, err)
	assert.Equal(t, "Syntax: !triggy <Begriff> ( <Index> | <VonIndex>:<BisIndex> )*", out.Text)
}
