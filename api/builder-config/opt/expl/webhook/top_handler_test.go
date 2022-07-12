package webhook_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/chschu/Klio.expl/types"
	"github.com/chschu/Klio.expl/webhook"
	"github.com/chschu/Klio.expl/webhook/generated/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_TopHandler_Top_Success_NoResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	topperMock := mocks.NewMockTopper(ctrl)

	maxTopResults := 1431

	sut := webhook.NewTopHandler(topperMock, maxTopResults)

	in := &webhook.Request{
		UserName: "nouser",
		Text:     "!top ",
	}

	ctx := context.WithValue(context.Background(), "dummy", "dummy")
	req := httptest.NewRequest("DUMMY", "/dummy", nil).WithContext(ctx)

	now := time.Now()

	var entries []types.Entry

	topperMock.EXPECT().Top(ctx, maxTopResults).Return(entries, nil)

	out, err := sut.Handle(in, req, now)

	assert.NoError(t, err)
	assert.Equal(t, "Ich habe leider keinen Eintrag gefunden.", out.Text)
}

func Test_TopHandler_Top_Success_SingleResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	topperMock := mocks.NewMockTopper(ctrl)

	maxTopResults := 8172

	sut := webhook.NewTopHandler(topperMock, maxTopResults)

	in := &webhook.Request{
		UserName: "singleuser",
		Text:     "  !top ",
	}

	ctx := context.WithValue(context.Background(), "dummy", "dummy")
	req := httptest.NewRequest("DUMMY", "/dummy", nil).WithContext(ctx)

	now := time.Now()

	entryKey1 := "k1"
	entryCount1 := uint(147124)
	entries := []types.Entry{{Key: entryKey1, HeadIndex: entryCount1}}

	topperMock.EXPECT().Top(ctx, maxTopResults).Return(entries, nil)

	out, err := sut.Handle(in, req, now)

	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Ich habe den folgenden wichtigsten Eintrag gefunden:\n```\n%s(%d)\n```",
		entryKey1, entryCount1), out.Text)
}

func Test_TopHandler_Top_Success_MultipleResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	topperMock := mocks.NewMockTopper(ctrl)

	maxTopResults := 3131

	sut := webhook.NewTopHandler(topperMock, maxTopResults)

	in := &webhook.Request{
		UserName: "multiuser",
		Text:     "!top   ",
	}

	ctx := context.WithValue(context.Background(), "dummy", "dummy")
	req := httptest.NewRequest("DUMMY", "/dummy", nil).WithContext(ctx)

	now := time.Now()

	entryKey1 := "k1"
	entryCount1 := uint(147124)
	entryKey2 := "K2"
	entryCount2 := uint(82736)
	entries := []types.Entry{{Key: entryKey1, HeadIndex: entryCount1}, {Key: entryKey2, HeadIndex: entryCount2}}

	topperMock.EXPECT().Top(ctx, maxTopResults).Return(entries, nil)

	out, err := sut.Handle(in, req, now)

	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("Ich habe die folgenden %d wichtigsten Einträge gefunden:\n```\n%s(%d), %s(%d)\n```",
		len(entries), entryKey1, entryCount1, entryKey2, entryCount2), out.Text)
}

func Test_TopHandler_Top_SoftFail_InvalidSyntax(t *testing.T) {
	ctrl := gomock.NewController(t)
	topperMock := mocks.NewMockTopper(ctrl)

	sut := webhook.NewTopHandler(topperMock, 10)

	in := &webhook.Request{
		UserName:    "pebkac",
		Text:        "!top redundant",
		TriggerWord: "!triggy",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	out, err := sut.Handle(in, req, time.Now())

	assert.NoError(t, err)
	assert.Equal(t, "Syntax: !triggy", out.Text)
}

func Test_TopHandler_Top_Fail_TopReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	topperMock := mocks.NewMockTopper(ctrl)

	sut := webhook.NewTopHandler(topperMock, 100)

	in := &webhook.Request{
		UserName: "unlucky",
		Text:     "!top",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	now := time.Now()

	expectedError := errors.New("expected error")

	topperMock.EXPECT().Top(gomock.Any(), gomock.Any()).Return(nil, expectedError)

	out, err := sut.Handle(in, req, now)

	assert.Nil(t, out)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}
