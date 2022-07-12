package webhook_test

import (
	"github.com/chschu/Klio.expl/types"
	"github.com/chschu/Klio.expl/webhook"
	"github.com/chschu/Klio.expl/webhook/generated/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_EntryListStringer_String_NoResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	jwtGeneratorMock := mocks.NewMockJWTGenerator(ctrl)
	entryStringerMock := mocks.NewMockEntryStringer(ctrl)

	sut := webhook.NewEntryListStringer(jwtGeneratorMock, entryStringerMock)

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	text := sut.String([]types.Entry{}, 0, "dummy", "/dummy/", time.Now(), req)

	assert.Equal(t, "Ich habe leider keinen Eintrag gefunden.", text)
}

func Test_EntryListStringer_String_OneResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	jwtGeneratorMock := mocks.NewMockJWTGenerator(ctrl)
	entryStringerMock := mocks.NewMockEntryStringer(ctrl)

	sut := webhook.NewEntryListStringer(jwtGeneratorMock, entryStringerMock)

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	entry1 := &types.Entry{Id: rand.Int31()}
	entry1String := uuid.Must(uuid.NewUUID()).String()
	entries := []types.Entry{*entry1}

	entryStringerMock.EXPECT().String(entry1).Return(entry1String)

	text := sut.String(entries, 1, "dummy", "/dummy/", time.Now(), req)

	assert.Equal(t, "Ich habe den folgenden Eintrag gefunden:\n```\n"+entry1String+"\n```", text)
}

func Test_EntryListStringer_String_MultipleResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	jwtGeneratorMock := mocks.NewMockJWTGenerator(ctrl)
	entryStringerMock := mocks.NewMockEntryStringer(ctrl)

	sut := webhook.NewEntryListStringer(jwtGeneratorMock, entryStringerMock)

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	entry1 := &types.Entry{Id: rand.Int31()}
	entry1String := uuid.Must(uuid.NewUUID()).String()
	entry2 := &types.Entry{Id: entry1.Id + 1}
	entry2String := uuid.Must(uuid.NewUUID()).String()
	entries := []types.Entry{*entry1, *entry2}

	entryStringerMock.EXPECT().String(entry1).Return(entry1String)
	entryStringerMock.EXPECT().String(entry2).Return(entry2String)

	text := sut.String(entries, 2, "dummy", "/dummy/", time.Now(), req)

	assert.Equal(t, "Ich habe die folgenden 2 Einträge gefunden:\n```\n"+
		entry1String+"\n"+entry2String+"\n```", text)
}

func Test_EntryListStringer_String_IncompleteResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	jwtGeneratorMock := mocks.NewMockJWTGenerator(ctrl)
	entryStringerMock := mocks.NewMockEntryStringer(ctrl)

	sut := webhook.NewEntryListStringer(jwtGeneratorMock, entryStringerMock)

	req := httptest.NewRequest("DUMMY", "https://example.com:8443/anything/here", nil)

	entry1 := &types.Entry{Id: rand.Int31()}
	entry1String := uuid.Must(uuid.NewUUID()).String()
	entry2 := &types.Entry{Id: entry1.Id + 1}
	entry2String := uuid.Must(uuid.NewUUID()).String()
	entries := []types.Entry{*entry1, *entry2}

	webUrlSubject := "test-subject"
	webUrlExpiresAt := time.Now()

	jwtGeneratorMock.EXPECT().GenerateJWT(webUrlSubject, webUrlExpiresAt).Return("eyJWT", nil)
	entryStringerMock.EXPECT().String(entry1).Return(entry1String)
	entryStringerMock.EXPECT().String(entry2).Return(entry2String)

	text := sut.String(entries, 91238, webUrlSubject, "/web-prefix/", webUrlExpiresAt, req)

	assert.Equal(t, "Ich habe [91238 Einträge](https://example.com:8443/web-prefix/eyJWT) gefunden, das sind die letzten 2:\n```\n"+
		entry1String+"\n"+entry2String+"\n```", text)
}
