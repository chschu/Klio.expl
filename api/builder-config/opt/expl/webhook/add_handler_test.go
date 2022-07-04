package webhook_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"klio/expl/types"
	"klio/expl/webhook"
	"klio/expl/webhook/generated/mocks"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_AddHandler_Add_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	adderMock := mocks.NewMockAdder(ctrl)
	entryStringerMock := mocks.NewMockEntryStringer(ctrl)

	sut := webhook.NewAddHandler(adderMock, entryStringerMock, 50, 50)

	in := &webhook.Request{
		UserName: "user",
		Text:     "!add this is a test",
	}

	ctx := context.WithValue(context.Background(), "dummy", "dummy")
	req := httptest.NewRequest("DUMMY", "/dummy", nil).WithContext(ctx)

	now := time.Now()

	entry := &types.Entry{Id: rand.Int31()}
	entryString := uuid.Must(uuid.NewUUID()).String()

	adderMock.EXPECT().Add(ctx, "this", "is a test", "user", now).Return(entry, nil)
	entryStringerMock.EXPECT().String(entry).Return(entryString)

	out, err := sut.Handle(in, req, now)

	assert.NoError(t, err)
	assert.Equal(t, "Ich habe den folgenden neuen Eintrag hinzugefügt:\n```\n"+entryString+"\n```", out.Text)
}

func Test_AddHandler_Add_SoftFail_InvalidSyntax(t *testing.T) {
	ctrl := gomock.NewController(t)
	adderMock := mocks.NewMockAdder(ctrl)
	entryStringerMock := mocks.NewMockEntryStringer(ctrl)

	sut := webhook.NewAddHandler(adderMock, entryStringerMock, 50, 50)

	in := &webhook.Request{
		UserName:    "pebkac",
		Text:        "!add not-explained",
		TriggerWord: "!trigger",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	out, err := sut.Handle(in, req, time.Now())

	assert.NoError(t, err)
	assert.Equal(t, "Syntax: !trigger <Begriff> <Erklärung>", out.Text)
}

func Test_AddHandler_Add_SoftFail_KeyTooLong(t *testing.T) {
	ctrl := gomock.NewController(t)
	adderMock := mocks.NewMockAdder(ctrl)
	entryStringerMock := mocks.NewMockEntryStringer(ctrl)

	sut := webhook.NewAddHandler(adderMock, entryStringerMock, 9, 50)

	in := &webhook.Request{
		UserName: "emojifan",
		Text:     "!add 😇👍😘😋😱 those are great!",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	out, err := sut.Handle(in, req, time.Now())

	assert.NoError(t, err)
	assert.Equal(t, "Tut mir leid, der Begriff ist leider zu lang.", out.Text)
}

func Test_AddHandler_Add_SoftFail_ValueTooLong(t *testing.T) {
	ctrl := gomock.NewController(t)
	adderMock := mocks.NewMockAdder(ctrl)
	entryStringerMock := mocks.NewMockEntryStringer(ctrl)

	sut := webhook.NewAddHandler(adderMock, entryStringerMock, 50, 15)

	in := &webhook.Request{
		UserName: "verbosedude",
		Text:     "!add key this is too long",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	out, err := sut.Handle(in, req, time.Now())

	assert.NoError(t, err)
	assert.Equal(t, "Tut mir leid, die Erklärung ist leider zu lang.", out.Text)
}

func Test_AddHandler_Add_Fail_AddReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	adderMock := mocks.NewMockAdder(ctrl)
	entryStringerMock := mocks.NewMockEntryStringer(ctrl)

	sut := webhook.NewAddHandler(adderMock, entryStringerMock, 50, 50)

	in := &webhook.Request{
		UserName: "unlucky",
		Text:     "!add this is unfortunate",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	now := time.Now()

	expectedError := errors.New("expected error")

	adderMock.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, expectedError)

	out, err := sut.Handle(in, req, now)

	assert.Nil(t, out)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}
