package webhook_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"klio/expl/generated/expldb_mocks"
	"klio/expl/generated/security_mocks"
	"klio/expl/generated/types_mocks"
	"klio/expl/generated/webhook_mocks"
	"klio/expl/types"
	"klio/expl/webhook"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_FindHandler_Find_ParamsPassedToFinder(t *testing.T) {
	ctrl := gomock.NewController(t)
	finderMock := expldb_mocks.NewMockFinder(ctrl)
	jwtGeneratorMock := security_mocks.NewMockJwtGenerator(ctrl)
	entryStringerMock := types_mocks.NewMockEntryStringer(ctrl)
	settingsMock := webhook_mocks.NewMockFindHandlerSettings(ctrl)
	sut := webhook.NewFindHandler(finderMock, "/prefix/", jwtGeneratorMock, entryStringerMock, settingsMock)

	in := &webhook.Request{
		Text: "!find test-regex  ",
	}

	ctx := context.WithValue(context.Background(), "dummy", "dummy")
	req := httptest.NewRequest("DUMMY", "/dummy", nil).WithContext(ctx)

	limit := 278372

	settingsMock.EXPECT().MaxFindCount().Return(limit)
	finderMock.EXPECT().FindWithLimit(ctx, "test-regex", limit).Return(nil, 0, nil)

	_, err := sut.Handle(in, req, time.Now())
	assert.NoError(t, err)
}

func Test_FindHandler_Find_FinderResultReturned_OneResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	finderMock := expldb_mocks.NewMockFinder(ctrl)
	jwtGeneratorMock := security_mocks.NewMockJwtGenerator(ctrl)
	entryStringerMock := types_mocks.NewMockEntryStringer(ctrl)
	settingsMock := webhook_mocks.NewMockFindHandlerSettings(ctrl)
	sut := webhook.NewFindHandler(finderMock, "/prefix/", jwtGeneratorMock, entryStringerMock, settingsMock)

	in := &webhook.Request{
		Text: "!find foo",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	entry1 := &types.Entry{Id: rand.Int31()}
	entry1String := uuid.Must(uuid.NewUUID()).String()
	entries := []types.Entry{*entry1}

	settingsMock.EXPECT().MaxFindCount()
	finderMock.EXPECT().FindWithLimit(gomock.Any(), gomock.Any(), gomock.Any()).Return(entries, 1, nil)
	entryStringerMock.EXPECT().String(entry1).Return(entry1String)

	out, err := sut.Handle(in, req, time.Now())

	assert.NoError(t, err)
	assert.Equal(t, "Ich habe den folgenden Eintrag gefunden:\n```\n"+entry1String+"\n```", out.Text)
}

func Test_FindHandler_Find_FinderResultReturned_MultipleResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	finderMock := expldb_mocks.NewMockFinder(ctrl)
	jwtGeneratorMock := security_mocks.NewMockJwtGenerator(ctrl)
	entryStringerMock := types_mocks.NewMockEntryStringer(ctrl)
	settingsMock := webhook_mocks.NewMockFindHandlerSettings(ctrl)
	sut := webhook.NewFindHandler(finderMock, "/prefix/", jwtGeneratorMock, entryStringerMock, settingsMock)

	in := &webhook.Request{
		Text: "!find foo",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	entry1 := &types.Entry{Id: rand.Int31()}
	entry1String := uuid.Must(uuid.NewUUID()).String()
	entry2 := &types.Entry{Id: entry1.Id + 1}
	entry2String := uuid.Must(uuid.NewUUID()).String()
	entries := []types.Entry{*entry1, *entry2}

	settingsMock.EXPECT().MaxFindCount()
	finderMock.EXPECT().FindWithLimit(gomock.Any(), gomock.Any(), gomock.Any()).Return(entries, 2, nil)
	entryStringerMock.EXPECT().String(entry1).Return(entry1String)
	entryStringerMock.EXPECT().String(entry2).Return(entry2String)

	out, err := sut.Handle(in, req, time.Now())

	assert.NoError(t, err)
	assert.Equal(t, "Ich habe die folgenden 2 Einträge gefunden:\n```\n"+
		entry1String+"\n"+entry2String+"\n```", out.Text)
}

func Test_FindHandler_Find_FinderResultReturned_OverLimitResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	finderMock := expldb_mocks.NewMockFinder(ctrl)
	jwtGeneratorMock := security_mocks.NewMockJwtGenerator(ctrl)
	entryStringerMock := types_mocks.NewMockEntryStringer(ctrl)
	settingsMock := webhook_mocks.NewMockFindHandlerSettings(ctrl)
	sut := webhook.NewFindHandler(finderMock, "/find-prefix/", jwtGeneratorMock, entryStringerMock, settingsMock)

	in := &webhook.Request{
		Text: "!find bar",
	}

	req := httptest.NewRequest("DUMMY", "https://example.com:8443/dummy", nil)

	now := time.Now()
	jwtValidity := 97837 * time.Second

	entry1 := &types.Entry{Id: rand.Int31()}
	entry1String := uuid.Must(uuid.NewUUID()).String()
	entry2 := &types.Entry{Id: entry1.Id + 1}
	entry2String := uuid.Must(uuid.NewUUID()).String()
	entries := []types.Entry{*entry1, *entry2}

	settingsMock.EXPECT().MaxFindCount()
	finderMock.EXPECT().FindWithLimit(gomock.Any(), gomock.Any(), gomock.Any()).Return(entries, 81738, nil)
	settingsMock.EXPECT().FindTokenValidity().Return(jwtValidity)
	jwtGeneratorMock.EXPECT().Generate("bar", now.Add(jwtValidity)).Return("eyJWT", nil)
	entryStringerMock.EXPECT().String(entry1).Return(entry1String)
	entryStringerMock.EXPECT().String(entry2).Return(entry2String)

	out, err := sut.Handle(in, req, now)

	assert.NoError(t, err)
	assert.Equal(t, "Ich habe [81738 Einträge](https://example.com:8443/find-prefix/eyJWT) gefunden, das sind die letzten 2:\n```\n"+
		entry1String+"\n"+entry2String+"\n```", out.Text)
}
