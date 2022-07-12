package webhook_test

import (
	"context"
	"errors"
	"github.com/chschu/Klio.expl/service"
	"github.com/chschu/Klio.expl/types"
	"github.com/chschu/Klio.expl/webhook"
	"github.com/chschu/Klio.expl/webhook/generated/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_FindHandler_Find_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	finderMock := mocks.NewMockLimitedFinder(ctrl)
	entryListStringerMock := mocks.NewMockEntryListStringer(ctrl)

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
	finderMock := mocks.NewMockLimitedFinder(ctrl)
	entryListStringerMock := mocks.NewMockEntryListStringer(ctrl)

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

func Test_FindHandler_Find_SoftFail_InvalidRegex(t *testing.T) {
	ctrl := gomock.NewController(t)
	finderMock := mocks.NewMockLimitedFinder(ctrl)
	entryListStringerMock := mocks.NewMockEntryListStringer(ctrl)

	sut := webhook.NewFindHandler(finderMock, entryListStringerMock, 50, "/dummy/", time.Minute)

	in := &webhook.Request{
		UserName:    "regexnoob",
		Text:        "!find this-is-not-a-regex",
		TriggerWord: "!triggy",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	now := time.Now()

	finderMock.EXPECT().FindWithLimit(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, 0, service.ErrFindRegexInvalid)

	out, err := sut.Handle(in, req, now)

	assert.NoError(t, err)
	assert.Equal(t, "Syntax: !triggy <POSIX-Regex>", out.Text)
}

func Test_FindHandler_Find_Fail_FindReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	finderMock := mocks.NewMockLimitedFinder(ctrl)
	entryListStringerMock := mocks.NewMockEntryListStringer(ctrl)

	sut := webhook.NewFindHandler(finderMock, entryListStringerMock, 50, "/dummy/", time.Minute)

	in := &webhook.Request{
		UserName: "unlucky",
		Text:     "!find this-is-unfortunate",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	now := time.Now()

	expectedError := errors.New("expected error")

	finderMock.EXPECT().FindWithLimit(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, 0, expectedError)

	out, err := sut.Handle(in, req, now)

	assert.Nil(t, out)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}
