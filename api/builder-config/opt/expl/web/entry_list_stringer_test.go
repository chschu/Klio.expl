package web_test

import (
	"github.com/chschu/Klio.expl/types"
	"github.com/chschu/Klio.expl/web"
	"github.com/chschu/Klio.expl/web/generated/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func Test_EntryListStringer_String_NoResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	entryStringerMock := mocks.NewMockEntryStringer(ctrl)

	sut := web.NewEntryListStringer(entryStringerMock)

	text := sut.String([]types.Entry{})

	assert.Equal(t, "Ich habe leider keinen Eintrag gefunden.\n", text)
}

func Test_EntryListStringer_String_OneResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	entryStringerMock := mocks.NewMockEntryStringer(ctrl)

	sut := web.NewEntryListStringer(entryStringerMock)

	entry1 := &types.Entry{Id: rand.Int31()}
	entry1String := uuid.Must(uuid.NewUUID()).String()
	entries := []types.Entry{*entry1}

	entryStringerMock.EXPECT().String(entry1).Return(entry1String)

	text := sut.String(entries)

	assert.Equal(t, "Ich habe den folgenden Eintrag gefunden:\n"+entry1String+"\n", text)
}

func Test_EntryListStringer_String_MultipleResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	entryStringerMock := mocks.NewMockEntryStringer(ctrl)

	sut := web.NewEntryListStringer(entryStringerMock)

	entry1 := &types.Entry{Id: rand.Int31()}
	entry1String := uuid.Must(uuid.NewUUID()).String()
	entry2 := &types.Entry{Id: entry1.Id + 1}
	entry2String := uuid.Must(uuid.NewUUID()).String()
	entries := []types.Entry{*entry1, *entry2}

	entryStringerMock.EXPECT().String(entry1).Return(entry1String)
	entryStringerMock.EXPECT().String(entry2).Return(entry2String)

	text := sut.String(entries)

	assert.Equal(t, "Ich habe die folgenden 2 Einträge gefunden:\n"+entry1String+"\n"+entry2String+"\n", text)
}
