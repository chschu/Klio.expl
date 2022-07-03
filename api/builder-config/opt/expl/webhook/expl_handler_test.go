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

func Test_ExplHandler_Expl_SoftFail_InvalidSyntax(t *testing.T) {
	ctrl := gomock.NewController(t)
	explainerMock := expldb_mocks.NewMockExplainer(ctrl)
	indexSpecParserMock := types_mocks.NewMockIndexSpecParser(ctrl)
	jwtGeneratorMock := security_mocks.NewMockJwtGenerator(ctrl)
	entryStringerMock := types_mocks.NewMockEntryStringer(ctrl)
	settingsMock := webhook_mocks.NewMockExplHandlerSettings(ctrl)
	sut := webhook.NewExplHandler(explainerMock, "/prefix/", indexSpecParserMock, jwtGeneratorMock, entryStringerMock, settingsMock)

	in := &webhook.Request{
		Text:        "!expl ",
		TriggerWord: "!triggy",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	out, err := sut.Handle(in, req, time.Now())

	assert.NoError(t, err)
	assert.Equal(t, "Syntax: !triggy <Begriff> ( <Index> | <VonIndex>:<BisIndex> )*", out.Text)

}

func Test_ExplHandler_Expl_ParamsPassedToExplainer_WithIndexSpec(t *testing.T) {
	ctrl := gomock.NewController(t)
	explainerMock := expldb_mocks.NewMockExplainer(ctrl)
	indexSpecParserMock := types_mocks.NewMockIndexSpecParser(ctrl)
	jwtGeneratorMock := security_mocks.NewMockJwtGenerator(ctrl)
	entryStringerMock := types_mocks.NewMockEntryStringer(ctrl)
	settingsMock := webhook_mocks.NewMockExplHandlerSettings(ctrl)
	sut := webhook.NewExplHandler(explainerMock, "/prefix/", indexSpecParserMock, jwtGeneratorMock, entryStringerMock, settingsMock)

	in := &webhook.Request{
		Text: "!expl test-key index-spec-string  ",
	}

	ctx := context.WithValue(context.Background(), "dummy", "dummy")
	req := httptest.NewRequest("DUMMY", "/dummy", nil).WithContext(ctx)

	limit := 278372

	indexSpecMock := types_mocks.NewMockIndexSpec(ctrl)

	indexSpecParserMock.EXPECT().ParseIndexSpec("index-spec-string").Return(indexSpecMock, nil)
	settingsMock.EXPECT().MaxExplCount().Return(limit)
	explainerMock.EXPECT().ExplainWithLimit(ctx, "test-key", indexSpecMock, limit).Return(nil, 0, nil)

	_, err := sut.Handle(in, req, time.Now())
	assert.NoError(t, err)
}

func Test_ExplHandler_Expl_ParamsPassedToExplainer_WithoutIndexSpec(t *testing.T) {
	ctrl := gomock.NewController(t)
	explainerMock := expldb_mocks.NewMockExplainer(ctrl)
	indexSpecParserMock := types_mocks.NewMockIndexSpecParser(ctrl)
	jwtGeneratorMock := security_mocks.NewMockJwtGenerator(ctrl)
	entryStringerMock := types_mocks.NewMockEntryStringer(ctrl)
	settingsMock := webhook_mocks.NewMockExplHandlerSettings(ctrl)
	sut := webhook.NewExplHandler(explainerMock, "/prefix/", indexSpecParserMock, jwtGeneratorMock, entryStringerMock, settingsMock)

	in := &webhook.Request{
		Text: "!expl test-key ",
	}

	ctx := context.WithValue(context.Background(), "dummy", "dummy")
	req := httptest.NewRequest("DUMMY", "/dummy", nil).WithContext(ctx)

	limit := 3918293

	settingsMock.EXPECT().MaxExplCount().Return(limit)
	explainerMock.EXPECT().ExplainWithLimit(ctx, "test-key", types.IndexSpecAll(), limit).Return(nil, 0, nil)

	_, err := sut.Handle(in, req, time.Now())
	assert.NoError(t, err)
}

func Test_ExplHandler_Expl_ExplainerResultReturned_OneResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	explainerMock := expldb_mocks.NewMockExplainer(ctrl)
	indexSpecParserMock := types_mocks.NewMockIndexSpecParser(ctrl)
	jwtGeneratorMock := security_mocks.NewMockJwtGenerator(ctrl)
	entryStringerMock := types_mocks.NewMockEntryStringer(ctrl)
	settingsMock := webhook_mocks.NewMockExplHandlerSettings(ctrl)
	sut := webhook.NewExplHandler(explainerMock, "/prefix/", indexSpecParserMock, jwtGeneratorMock, entryStringerMock, settingsMock)

	in := &webhook.Request{
		Text: "!expl foo",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	entry1 := &types.Entry{Id: rand.Int31()}
	entry1String := uuid.Must(uuid.NewUUID()).String()
	entries := []types.Entry{*entry1}

	settingsMock.EXPECT().MaxExplCount()
	explainerMock.EXPECT().ExplainWithLimit(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(entries, 1, nil)
	entryStringerMock.EXPECT().String(entry1).Return(entry1String)

	out, err := sut.Handle(in, req, time.Now())

	assert.NoError(t, err)
	assert.Equal(t, "Ich habe den folgenden Eintrag gefunden:\n```\n"+entry1String+"\n```", out.Text)
}

func Test_ExplHandler_Expl_ExplainerResultReturned_MultipleResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	explainerMock := expldb_mocks.NewMockExplainer(ctrl)
	indexSpecParserMock := types_mocks.NewMockIndexSpecParser(ctrl)
	jwtGeneratorMock := security_mocks.NewMockJwtGenerator(ctrl)
	entryStringerMock := types_mocks.NewMockEntryStringer(ctrl)
	settingsMock := webhook_mocks.NewMockExplHandlerSettings(ctrl)
	sut := webhook.NewExplHandler(explainerMock, "/prefix/", indexSpecParserMock, jwtGeneratorMock, entryStringerMock, settingsMock)

	in := &webhook.Request{
		Text: "!expl foo",
	}

	req := httptest.NewRequest("DUMMY", "/dummy", nil)

	entry1 := &types.Entry{Id: rand.Int31()}
	entry1String := uuid.Must(uuid.NewUUID()).String()
	entry2 := &types.Entry{Id: entry1.Id + 1}
	entry2String := uuid.Must(uuid.NewUUID()).String()
	entries := []types.Entry{*entry1, *entry2}

	settingsMock.EXPECT().MaxExplCount()
	explainerMock.EXPECT().ExplainWithLimit(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(entries, 2, nil)
	entryStringerMock.EXPECT().String(entry1).Return(entry1String)
	entryStringerMock.EXPECT().String(entry2).Return(entry2String)

	out, err := sut.Handle(in, req, time.Now())

	assert.NoError(t, err)
	assert.Equal(t, "Ich habe die folgenden 2 Einträge gefunden:\n```\n"+
		entry1String+"\n"+entry2String+"\n```", out.Text)
}

func Test_ExplHandler_Expl_ExplainerResultReturned_OverLimitResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	explainerMock := expldb_mocks.NewMockExplainer(ctrl)
	indexSpecParserMock := types_mocks.NewMockIndexSpecParser(ctrl)
	jwtGeneratorMock := security_mocks.NewMockJwtGenerator(ctrl)
	entryStringerMock := types_mocks.NewMockEntryStringer(ctrl)
	settingsMock := webhook_mocks.NewMockExplHandlerSettings(ctrl)
	sut := webhook.NewExplHandler(explainerMock, "/web-expl-prefix/", indexSpecParserMock, jwtGeneratorMock, entryStringerMock, settingsMock)

	in := &webhook.Request{
		Text: "!expl foo",
	}

	req := httptest.NewRequest("DUMMY", "https://example.org:1234/dummy", nil)

	now := time.Now()
	jwtValidity := 13921 * time.Second

	entry1 := &types.Entry{Id: rand.Int31()}
	entry1String := uuid.Must(uuid.NewUUID()).String()
	entry2 := &types.Entry{Id: entry1.Id + 1}
	entry2String := uuid.Must(uuid.NewUUID()).String()
	entries := []types.Entry{*entry1, *entry2}

	settingsMock.EXPECT().MaxExplCount()
	explainerMock.EXPECT().ExplainWithLimit(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(entries, 198319, nil)
	settingsMock.EXPECT().ExplTokenValidity().Return(jwtValidity)
	jwtGeneratorMock.EXPECT().Generate("foo", now.Add(jwtValidity)).Return("eyJWT", nil)
	entryStringerMock.EXPECT().String(entry1).Return(entry1String)
	entryStringerMock.EXPECT().String(entry2).Return(entry2String)

	out, err := sut.Handle(in, req, now)

	assert.NoError(t, err)
	assert.Equal(t, "Ich habe [198319 Einträge](https://example.org:1234/web-expl-prefix/eyJWT) gefunden, das sind die letzten 2:\n```\n"+
		entry1String+"\n"+entry2String+"\n```", out.Text)
}
