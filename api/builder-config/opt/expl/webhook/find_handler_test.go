package webhook_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"klio/expl/generated/webhook_mocks"
	"klio/expl/types"
	"klio/expl/webhook"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_FindHandler_Find_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	finderMock := webhook_mocks.NewMockLimitedFinder(ctrl)
	entryListStringerMock := webhook_mocks.NewMockEntryListStringer(ctrl)

	maxImmediateResults := 278372
	webUrlPathPrefix := "/web-find-prefix/"
	webUrlValidity := 97837 * time.Second

	sut := webhook.NewFindHandler(finderMock, entryListStringerMock, maxImmediateResults, webUrlPathPrefix, webUrlValidity)

	in := &webhook.Request{
		Text: "!find test-regex  ",
	}

	ctx := context.WithValue(context.Background(), "dummy", "dummy")
	req := httptest.NewRequest("DUMMY", "/dummy", nil).WithContext(ctx)

	now := time.Now()

	entries := []types.Entry{{Id: rand.Int31()}, {Id: rand.Int31()}}
	total := 391247
	resultText := "result text!"

	finderMock.EXPECT().FindWithLimit(ctx, "test-regex", maxImmediateResults).Return(entries, total, nil)
	entryListStringerMock.EXPECT().String(entries, total, "test-regex", webUrlPathPrefix, now.Add(webUrlValidity), req).Return(resultText)

	out, err := sut.Handle(in, req, now)
	assert.NoError(t, err)
	assert.Equal(t, resultText, out.Text)
}

func Test_FindHandler_Find_SoftFail_InvalidSyntax(t *testing.T) {
	ctrl := gomock.NewController(t)
	finderMock := webhook_mocks.NewMockLimitedFinder(ctrl)
	entryListStringerMock := webhook_mocks.NewMockEntryListStringer(ctrl)

	sut := webhook.NewFindHandler(finderMock, entryListStringerMock, 50, "/dummy/", time.Minute)

	in := &webhook.Request{
		Text:        "!find this is invalid ",
		TriggerWord: "!trig",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	out, err := sut.Handle(in, req, time.Now())

	assert.NoError(t, err)
	assert.Equal(t, "Syntax: !trig <POSIX-Regex>", out.Text)
}
