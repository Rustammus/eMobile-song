package v1

import (
	"database/sql"
	"eMobile/internal/config"
	"eMobile/internal/crud"
	"eMobile/internal/dto"
	"eMobile/internal/service"
	mockservice "eMobile/internal/service/mocks"
	"eMobile/pkg/logging"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func parseTime(l, s string) time.Time {
	t, err := time.Parse(l, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestHandler_audioCreate(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockIAudioService, user *dto.AudioCreate)

	testTable := []struct {
		name            string
		inputBody       string
		inputDTO        *dto.AudioCreate
		mockBehaviour   mockBehaviour
		expectedCode    int
		expectedBody    string
		bodyMustContain string
	}{
		{
			name: "201_valid_input",
			inputBody: `{
							"group": "classic",
							"song": "Some song"
						}`,
			inputDTO: &dto.AudioCreate{
				Group: "classic",
				Song:  "Some song",
			},
			mockBehaviour: func(s *mockservice.MockIAudioService, audio *dto.AudioCreate) {
				s.EXPECT().Create(audio).Return(
					pgtype.UUID{Valid: true},
					nil,
				)
			},
			expectedCode: 201,
			expectedBody: `{"data":"00000000-0000-0000-0000-000000000000", "message":"audio created correctly"}`,
		},
		{
			name: "400_invalid_values",
			inputBody: `{
							"group": "",
							"song": ""
						}`,
			inputDTO: &dto.AudioCreate{
				Group: "",
				Song:  "",
			},
			mockBehaviour: func(s *mockservice.MockIAudioService, audio *dto.AudioCreate) {
			},
			expectedCode: 400,
			expectedBody: `{"error":"'group' is required and cannot be empty;'song' is required and cannot be empty;", "message":"validation err"}`,
		},
		{
			name: "400_invalid_struct",
			inputBody: `{
							"group": "",
						}`,
			inputDTO: &dto.AudioCreate{
				Group: "",
				Song:  "",
			},
			mockBehaviour: func(s *mockservice.MockIAudioService, audio *dto.AudioCreate) {
			},
			expectedCode: 400,
			expectedBody: `{"error":"invalid character '}' looking for beginning of object key string", "message":"read body err"}`,
		},
		{
			name: "500_unknown_err",
			inputBody: `{
							"group": "classic",
							"song": "Some song"
						}`,
			inputDTO: &dto.AudioCreate{
				Group: "classic",
				Song:  "Some song",
			},
			mockBehaviour: func(s *mockservice.MockIAudioService, audio *dto.AudioCreate) {
				s.EXPECT().Create(audio).Return(
					pgtype.UUID{},
					errors.New("unknown error"),
				)
			},
			expectedCode: 500,
			expectedBody: `{"error":"unknown error", "message":"create audio err"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()
			audioService := mockservice.NewMockIAudioService(c)
			testCase.mockBehaviour(audioService, testCase.inputDTO)

			services := service.Service{Audio: audioService}
			handler := NewHandler(Deps{
				Service: services,
				Logger:  logging.GetLoggerTest(),
			})

			//Test server
			r := httprouter.New()
			r.POST("/audios", handler.audioCreate)

			//http test
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/audios", strings.NewReader(testCase.inputBody))

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != "" {
				if testCase.expectedBody == "nil" {
					testCase.expectedBody = "{}"
				}
				assert.JSONEq(t, testCase.expectedBody, w.Body.String())
			}
			if testCase.bodyMustContain != "" {
				assert.Contains(t, w.Body.String(), testCase.bodyMustContain)
			}
		})
	}
}

func TestHandler_audioList(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockIAudioService, filter *dto.AudioFilter, pag crud.Pagination)

	testTable := []struct {
		name            string
		inputQuery      string
		inputPag        crud.Pagination
		inputDTO        *dto.AudioFilter
		mockBehaviour   mockBehaviour
		expectedCode    int
		expectedBody    string
		bodyMustContain string
	}{
		{
			name:       "200_empty_query",
			inputQuery: "",
			inputPag:   crud.Pagination{0, 50},
			inputDTO:   nil,
			mockBehaviour: func(s *mockservice.MockIAudioService, filter *dto.AudioFilter, pag crud.Pagination) {
				s.EXPECT().ListPag(pag).Return([]dto.AudioRead{{
					UUID:        pgtype.UUID{Valid: true},
					Group:       "group1",
					Song:        "song1",
					ReleaseDate: pgtype.Date{Valid: true},
					Link:        "link1",
					CreatedAt:   pgtype.Timestamptz{Valid: true},
					UpdatedAt:   pgtype.Timestamptz{Valid: true},
				}}, nil)
			},
			expectedCode:    200,
			expectedBody:    `{"data": [{"created_at":"0001-01-01T00:00:00Z", "group":"group1", "link":"link1", "release_date":"0001-01-01", "song":"song1", "updated_at":"0001-01-01T00:00:00Z", "uuid":"00000000-0000-0000-0000-000000000000"}], "message":"audios got correctly", "next_pagination":{"limit":50, "offset":50}}`,
			bodyMustContain: "",
		},
		{
			name:       "200_full_query",
			inputQuery: "?group=group1&song=some%20text&after=2012-09-23&before=2013-09-23&link=link1&lyric=some%20lyric&limit=35&offset=10",
			inputPag:   crud.Pagination{10, 35},
			inputDTO: &dto.AudioFilter{
				Group:             sql.NullString{String: "group1", Valid: true},
				Song:              sql.NullString{String: "some text", Valid: true},
				ReleaseDateAfter:  pgtype.Date{Time: parseTime("2006-01-02", "2012-09-23"), Valid: true},
				ReleaseDateBefore: pgtype.Date{Time: parseTime("2006-01-02", "2013-09-23"), Valid: true},
				Link:              sql.NullString{String: "link1", Valid: true},
				Lyric:             sql.NullString{String: "some lyric", Valid: true},
			},
			mockBehaviour: func(s *mockservice.MockIAudioService, filter *dto.AudioFilter, pag crud.Pagination) {
				s.EXPECT().ListByFilter(filter, pag).Return([]dto.AudioRead{{
					UUID:        pgtype.UUID{Valid: true},
					Group:       "group1",
					Song:        "song1",
					ReleaseDate: pgtype.Date{Valid: true},
					Link:        "link1",
					CreatedAt:   pgtype.Timestamptz{Valid: true},
					UpdatedAt:   pgtype.Timestamptz{Valid: true},
				}}, nil)
			},
			expectedCode:    200,
			expectedBody:    `{"data": [{"created_at":"0001-01-01T00:00:00Z", "group":"group1", "link":"link1", "release_date":"0001-01-01", "song":"song1", "updated_at":"0001-01-01T00:00:00Z", "uuid":"00000000-0000-0000-0000-000000000000"}], "message":"audios got correctly", "next_pagination":{"limit":35, "offset":45}}`,
			bodyMustContain: "",
		},
		{
			name:       "400_invalid_query_date",
			inputQuery: "?after=2012-19-23&before=2013-19-23",
			inputPag:   crud.Pagination{0, 50},
			inputDTO:   nil,
			mockBehaviour: func(s *mockservice.MockIAudioService, filter *dto.AudioFilter, pag crud.Pagination) {
			},
			expectedCode:    400,
			expectedBody:    `{"error":"invalid after format, example: 2006-09-25", "message":"validation err"}`,
			bodyMustContain: "",
		},
		{
			name:       "200_no_rows",
			inputQuery: "",
			inputPag:   crud.Pagination{0, 50},
			inputDTO:   nil,
			mockBehaviour: func(s *mockservice.MockIAudioService, filter *dto.AudioFilter, pag crud.Pagination) {
				s.EXPECT().ListPag(pag).Return(nil, pgx.ErrNoRows)
			},
			expectedCode:    200,
			expectedBody:    `{"data": [], "message":"no rows find", "next_pagination": {"limit":50, "offset":50}}`,
			bodyMustContain: "",
		},
		{
			name:       "500_unknown_error",
			inputQuery: "",
			inputPag:   crud.Pagination{0, 50},
			inputDTO:   nil,
			mockBehaviour: func(s *mockservice.MockIAudioService, filter *dto.AudioFilter, pag crud.Pagination) {
				s.EXPECT().ListPag(pag).Return(nil, errors.New("unknown error"))
			},
			expectedCode:    500,
			expectedBody:    `{"error":"unknown error", "message":"error on list audios"}`,
			bodyMustContain: "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()
			audioService := mockservice.NewMockIAudioService(c)
			testCase.mockBehaviour(audioService, testCase.inputDTO, testCase.inputPag)

			services := service.Service{Audio: audioService}
			handler := NewHandler(Deps{
				Service: services,
				Logger:  logging.GetLoggerTest(),
				Config: &config.Config{
					Server: config.Server{
						PagLimit: 50,
					},
				},
			})

			//Test server
			r := httprouter.New()
			r.GET("/audios", handler.audioList)

			//http test
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/audios"+testCase.inputQuery, nil)

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != "" {
				if testCase.expectedBody == "nil" {
					testCase.expectedBody = "{}"
				}
				assert.JSONEq(t, testCase.expectedBody, w.Body.String())
			}
			if testCase.bodyMustContain != "" {
				assert.Contains(t, w.Body.String(), testCase.bodyMustContain)
			}
		})
	}
}

func TestHandler_audioFindByUUID(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockIAudioService, uuid pgtype.UUID)

	testTable := []struct {
		name            string
		inputQueryRaw   string
		inputPathUUID   string
		inputUUID       pgtype.UUID
		mockBehaviour   mockBehaviour
		expectedCode    int
		expectedBody    string
		bodyMustContain string
	}{
		{
			name:          "200_valid_uuid",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputUUID:     pgtype.UUID{Valid: true},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID) {
				s.EXPECT().Find(uuid).Return(&dto.AudioRead{
					UUID:        pgtype.UUID{Valid: true},
					Group:       "group1",
					Song:        "song1",
					ReleaseDate: pgtype.Date{Valid: true},
					Link:        "link1",
					CreatedAt:   pgtype.Timestamptz{Valid: true},
					UpdatedAt:   pgtype.Timestamptz{Valid: true},
				}, nil)
			},
			expectedCode:    200,
			expectedBody:    `{"data": {"created_at":"0001-01-01T00:00:00Z", "group":"group1", "link":"link1", "release_date":"0001-01-01", "song":"song1", "updated_at":"0001-01-01T00:00:00Z", "uuid":"00000000-0000-0000-0000-000000000000"}, "message":"audio got correctly"}`,
			bodyMustContain: "",
		},
		{
			name:          "200_valid_uuid_2",
			inputPathUUID: "00000000000000000000000000000000",
			inputUUID:     pgtype.UUID{Valid: true},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID) {
				s.EXPECT().Find(uuid).Return(&dto.AudioRead{
					UUID:        pgtype.UUID{Valid: true},
					Group:       "group1",
					Song:        "song1",
					ReleaseDate: pgtype.Date{Valid: true},
					Link:        "link1",
					CreatedAt:   pgtype.Timestamptz{Valid: true},
					UpdatedAt:   pgtype.Timestamptz{Valid: true},
				}, nil)
			},
			expectedCode:    200,
			expectedBody:    `{"data": {"created_at":"0001-01-01T00:00:00Z", "group":"group1", "link":"link1", "release_date":"0001-01-01", "song":"song1", "updated_at":"0001-01-01T00:00:00Z", "uuid":"00000000-0000-0000-0000-000000000000"}, "message":"audio got correctly"}`,
			bodyMustContain: "",
		},
		{
			name:          "200_valid_uuid_with_lyrics",
			inputQueryRaw: "?full=true",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputUUID:     pgtype.UUID{Valid: true},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID) {
				s.EXPECT().FindWithLyric(uuid).Return(&dto.AudioReadFull{
					UUID:        pgtype.UUID{Valid: true},
					Group:       "group1",
					Song:        "song1",
					ReleaseDate: pgtype.Date{Valid: true},
					Link:        "link1",
					Lyrics: []dto.LyricRead{
						{
							UUID:      pgtype.UUID{Valid: true},
							AudioUUID: pgtype.UUID{Valid: true},
							Order:     0,
							Text:      "lyric1",
							CreatedAt: pgtype.Timestamptz{Valid: true},
							UpdatedAt: pgtype.Timestamptz{Valid: true},
						},
						{
							UUID:      pgtype.UUID{Valid: true},
							AudioUUID: pgtype.UUID{Valid: true},
							Order:     1,
							Text:      "lyric2",
							CreatedAt: pgtype.Timestamptz{Valid: true},
							UpdatedAt: pgtype.Timestamptz{Valid: true},
						},
					},
					CreatedAt: pgtype.Timestamptz{Valid: true},
					UpdatedAt: pgtype.Timestamptz{Valid: true},
				}, nil)
			},
			expectedCode:    200,
			expectedBody:    `{"data": {"created_at":"0001-01-01T00:00:00Z", "group":"group1", "link":"link1", "release_date":"0001-01-01", "lyrics": [{"audio_uuid":"00000000-0000-0000-0000-000000000000", "created_at":"0001-01-01T00:00:00Z", "order":0, "text":"lyric1", "updated_at":"0001-01-01T00:00:00Z", "uuid":"00000000-0000-0000-0000-000000000000"}, {"audio_uuid":"00000000-0000-0000-0000-000000000000", "created_at":"0001-01-01T00:00:00Z", "order":1, "text":"lyric2", "updated_at":"0001-01-01T00:00:00Z", "uuid":"00000000-0000-0000-0000-000000000000"}], "song":"song1", "updated_at":"0001-01-01T00:00:00Z", "uuid":"00000000-0000-0000-0000-000000000000"}, "message":"audio got correctly"}`,
			bodyMustContain: "",
		},
		{
			name:          "400_invalid_uuid",
			inputPathUUID: "00000000-0000-0000-0000-00000000000k",
			inputUUID:     pgtype.UUID{},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID) {
			},
			expectedCode:    400,
			expectedBody:    `{"error":"encoding/hex: invalid byte: U+006B 'k'", "message":"invalid uuid in path param"}`,
			bodyMustContain: "",
		},
		{
			name:          "200_no_rows",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputUUID:     pgtype.UUID{Valid: true},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID) {
				s.EXPECT().Find(uuid).Return(nil, pgx.ErrNoRows)
			},
			expectedCode:    200,
			expectedBody:    `{"data": {}, "message":"no rows find"}`,
			bodyMustContain: "",
		},
		{
			name:          "500_unknown_error",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputUUID:     pgtype.UUID{Valid: true},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID) {
				s.EXPECT().Find(uuid).Return(nil, errors.New("unknown error"))
			},
			expectedCode:    500,
			expectedBody:    `{"error":"unknown error", "message":"error on find audio by uuid"}`,
			bodyMustContain: "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()
			audioService := mockservice.NewMockIAudioService(c)
			testCase.mockBehaviour(audioService, testCase.inputUUID)

			services := service.Service{Audio: audioService}
			handler := NewHandler(Deps{
				Service: services,
				Logger:  logging.GetLoggerTest(),
				Config: &config.Config{
					Server: config.Server{
						PagLimit: 50,
					},
				},
			})

			//Test server
			r := httprouter.New()
			r.GET("/audios/:uuid", handler.audioFindByUUID)

			//http test
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/audios/"+testCase.inputPathUUID+testCase.inputQueryRaw, nil)

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != "" {
				if testCase.expectedBody == "nil" {
					testCase.expectedBody = "{}"
				}
				assert.JSONEq(t, testCase.expectedBody, w.Body.String())
			}
			if testCase.bodyMustContain != "" {
				assert.Contains(t, w.Body.String(), testCase.bodyMustContain)
			}
		})
	}
}

func TestHandler_audioUpdateByUUID(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockIAudioService, uuid pgtype.UUID, audio *dto.AudioUpdate)
	testTable := []struct {
		name            string
		inputPathUUID   string
		inputBodyRaw    string
		inputUUID       pgtype.UUID
		inputAudio      *dto.AudioUpdate
		mockBehaviour   mockBehaviour
		expectedCode    int
		expectedBody    string
		bodyMustContain string
	}{
		{
			name:          "200_valid_input_full",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputBodyRaw: `{
								"group": "group22",
								"song": "song22",
								"release_date": "2010-11-12",
								"link": "link22",
								"lyrics": "new lyrics\nsame\n\nsecond_lyric"
							}`,
			inputUUID: pgtype.UUID{Valid: true},
			inputAudio: &dto.AudioUpdate{
				Group:       sql.NullString{String: "group22", Valid: true},
				Song:        sql.NullString{String: "song22", Valid: true},
				ReleaseDate: pgtype.Date{Time: parseTime("2006-01-02", "2010-11-12"), Valid: true},
				Link:        sql.NullString{String: "link22", Valid: true},
				LyricsRaw:   sql.NullString{String: "new lyrics\nsame\n\nsecond_lyric", Valid: true},
				Lyrics:      nil,
			},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID, audio *dto.AudioUpdate) {
				s.EXPECT().Update(uuid, audio).Return(&dto.AudioRead{
					UUID:        pgtype.UUID{Valid: true},
					Group:       "group22",
					Song:        "song22",
					ReleaseDate: pgtype.Date{Time: parseTime("2006-01-02", "2010-11-12"), Valid: true},
					Link:        "link22",
					CreatedAt:   pgtype.Timestamptz{Valid: true},
					UpdatedAt:   pgtype.Timestamptz{Valid: true},
				}, nil)
			},
			expectedCode:    200,
			expectedBody:    `{"data": {"created_at":"0001-01-01T00:00:00Z", "group":"group22", "link":"link22", "release_date":"2010-11-12", "song":"song22", "updated_at":"0001-01-01T00:00:00Z", "uuid":"00000000-0000-0000-0000-000000000000"}, "message":"audio updated correctly"}`,
			bodyMustContain: "",
		},
		{
			name:          "200_valid_input_lyrics_only",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputBodyRaw: `{
								"lyrics": "new lyrics\nsame\n\nsecond_lyric"
							}`,
			inputUUID: pgtype.UUID{Valid: true},
			inputAudio: &dto.AudioUpdate{
				LyricsRaw: sql.NullString{String: "new lyrics\nsame\n\nsecond_lyric", Valid: true},
			},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID, audio *dto.AudioUpdate) {
				s.EXPECT().Update(uuid, audio).Return(&dto.AudioRead{
					UUID:        pgtype.UUID{Valid: true},
					Group:       "group22",
					Song:        "song22",
					ReleaseDate: pgtype.Date{Time: parseTime("2006-01-02", "2010-11-12"), Valid: true},
					Link:        "link22",
					CreatedAt:   pgtype.Timestamptz{Valid: true},
					UpdatedAt:   pgtype.Timestamptz{Valid: true},
				}, nil)
			},
			expectedCode:    200,
			expectedBody:    `{"data": {"created_at":"0001-01-01T00:00:00Z", "group":"group22", "link":"link22", "release_date":"2010-11-12", "song":"song22", "updated_at":"0001-01-01T00:00:00Z", "uuid":"00000000-0000-0000-0000-000000000000"}, "message":"audio updated correctly"}`,
			bodyMustContain: "",
		},
		{
			name:          "400_invalid_date_input",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputBodyRaw: `{
								"group": "group22",
								"song": "song22",
								"release_date": "2010-13-12",
								"link": "link22",
								"lyrics": "new lyrics\nsame\n\nsecond_lyric"
							}`,
			inputUUID: pgtype.UUID{Valid: true},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID, audio *dto.AudioUpdate) {
			},
			expectedCode:    400,
			expectedBody:    `{"error":"invalid date format, example: 2006-09-25;", "message":"validation error"}`,
			bodyMustContain: "",
		},
		{
			name:          "400_invalid_uuid_path",
			inputPathUUID: "00000000-0000-0000-0000-00000000000k",
			inputBodyRaw: `{
								"group": "group22",
								"song": "song22",
								"release_date": "2010-03-12",
								"link": "link22",
								"lyrics": "new lyrics\nsame\n\nsecond_lyric"
							}`,
			inputUUID: pgtype.UUID{Valid: true},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID, audio *dto.AudioUpdate) {
			},
			expectedCode:    400,
			expectedBody:    `{"error":"encoding/hex: invalid byte: U+006B 'k'", "message":"invalid uuid in path param"}`,
			bodyMustContain: "",
		},
		{
			name:          "200_no_rows",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputBodyRaw: `{
								"group": "group22",
								"song": "song22",
								"release_date": "2010-11-12",
								"link": "link22",
								"lyrics": "new lyrics\nsame\n\nsecond_lyric"
							}`,
			inputUUID: pgtype.UUID{Valid: true},
			inputAudio: &dto.AudioUpdate{
				Group:       sql.NullString{String: "group22", Valid: true},
				Song:        sql.NullString{String: "song22", Valid: true},
				ReleaseDate: pgtype.Date{Time: parseTime("2006-01-02", "2010-11-12"), Valid: true},
				Link:        sql.NullString{String: "link22", Valid: true},
				LyricsRaw:   sql.NullString{String: "new lyrics\nsame\n\nsecond_lyric", Valid: true},
				Lyrics:      nil,
			},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID, audio *dto.AudioUpdate) {
				s.EXPECT().Update(uuid, audio).Return(nil, pgx.ErrNoRows)
			},
			expectedCode:    200,
			expectedBody:    `{"data": {}, "message":"no rows updated"}`,
			bodyMustContain: "",
		},
		{
			name:          "500_unknown_error",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputBodyRaw: `{
								"group": "group22",
								"song": "song22",
								"release_date": "2010-11-12",
								"link": "link22",
								"lyrics": "new lyrics\nsame\n\nsecond_lyric"
							}`,
			inputUUID: pgtype.UUID{Valid: true},
			inputAudio: &dto.AudioUpdate{
				Group:       sql.NullString{String: "group22", Valid: true},
				Song:        sql.NullString{String: "song22", Valid: true},
				ReleaseDate: pgtype.Date{Time: parseTime("2006-01-02", "2010-11-12"), Valid: true},
				Link:        sql.NullString{String: "link22", Valid: true},
				LyricsRaw:   sql.NullString{String: "new lyrics\nsame\n\nsecond_lyric", Valid: true},
				Lyrics:      nil,
			},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID, audio *dto.AudioUpdate) {
				s.EXPECT().Update(uuid, audio).Return(nil, errors.New("unknown error"))
			},
			expectedCode:    500,
			expectedBody:    `{"error":"unknown error", "message":"error on update audio"}`,
			bodyMustContain: "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()
			audioService := mockservice.NewMockIAudioService(c)
			testCase.mockBehaviour(audioService, testCase.inputUUID, testCase.inputAudio)

			services := service.Service{Audio: audioService}
			handler := NewHandler(Deps{
				Service: services,
				Logger:  logging.GetLoggerTest(),
				Config: &config.Config{
					Server: config.Server{
						PagLimit: 50,
					},
				},
			})

			//Test server
			r := httprouter.New()
			r.PATCH("/audios/:uuid", handler.audioUpdateByUUID)

			//http test
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/audios/"+testCase.inputPathUUID, strings.NewReader(testCase.inputBodyRaw))

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != "" {
				if testCase.expectedBody == "nil" {
					testCase.expectedBody = "{}"
				}
				assert.JSONEq(t, testCase.expectedBody, w.Body.String())
			}
			if testCase.bodyMustContain != "" {
				assert.Contains(t, w.Body.String(), testCase.bodyMustContain)
			}
		})
	}
}

func TestHandler_audioDeleteByUUID(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockIAudioService, uuid pgtype.UUID)
	testTable := []struct {
		name            string
		inputPathUUID   string
		inputUUID       pgtype.UUID
		mockBehaviour   mockBehaviour
		expectedCode    int
		expectedBody    string
		bodyMustContain string
	}{
		{
			name:          "200_valid_path",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputUUID:     pgtype.UUID{Valid: true},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID) {
				s.EXPECT().Delete(uuid).Return(nil)
			},
			expectedCode:    200,
			expectedBody:    `{"data":"00000000-0000-0000-0000-000000000000", "message":"audio deleted correctly"}`,
			bodyMustContain: "",
		},
		{
			name:          "200_no_rows",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputUUID:     pgtype.UUID{Valid: true},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID) {
				s.EXPECT().Delete(uuid).Return(pgx.ErrNoRows)
			},
			expectedCode:    200,
			expectedBody:    `{"data": {}, "message":"no rows deleted"}`,
			bodyMustContain: "",
		},
		{
			name:          "400_invalid_path_uuid",
			inputPathUUID: "00000000-0000-0000-0000-00000000000y",
			inputUUID:     pgtype.UUID{Valid: true},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID) {
			},
			expectedCode:    400,
			expectedBody:    `{"error":"encoding/hex: invalid byte: U+0079 'y'", "message":"invalid uuid in path param"}`,
			bodyMustContain: "",
		},
		{
			name:          "500_unknown_error",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputUUID:     pgtype.UUID{Valid: true},
			mockBehaviour: func(s *mockservice.MockIAudioService, uuid pgtype.UUID) {
				s.EXPECT().Delete(uuid).Return(errors.New("unknown error"))
			},
			expectedCode:    500,
			expectedBody:    `{"error":"unknown error", "message":"error on delete audio by uuid"}`,
			bodyMustContain: "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()
			audioService := mockservice.NewMockIAudioService(c)
			testCase.mockBehaviour(audioService, testCase.inputUUID)

			services := service.Service{Audio: audioService}
			handler := NewHandler(Deps{
				Service: services,
				Logger:  logging.GetLoggerTest(),
				Config: &config.Config{
					Server: config.Server{
						PagLimit: 50,
					},
				},
			})

			//Test server
			r := httprouter.New()
			r.DELETE("/audios/:uuid", handler.audioDeleteByUUID)

			//http test
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/audios/"+testCase.inputPathUUID, nil)

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != "" {
				if testCase.expectedBody == "nil" {
					testCase.expectedBody = "{}"
				}
				assert.JSONEq(t, testCase.expectedBody, w.Body.String())
			}
			if testCase.bodyMustContain != "" {
				assert.Contains(t, w.Body.String(), testCase.bodyMustContain)
			}
		})
	}
}

func TestHandler_audioLyricsList(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockILyricService, uuid pgtype.UUID, pag crud.Pagination)

	testTable := []struct {
		name            string
		inputPathUUID   string
		inputQuery      string
		inputUUID       pgtype.UUID
		inputPag        crud.Pagination
		mockBehaviour   mockBehaviour
		expectedCode    int
		expectedBody    string
		bodyMustContain string
	}{
		{
			name:          "200_valid_input_all",
			inputPathUUID: "00000000-0000-0000-0000-000000000001",
			inputQuery:    "?limit=35&offset=10",
			inputUUID:     pgtype.UUID{Valid: true, Bytes: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
			inputPag: crud.Pagination{
				Limit:  35,
				Offset: 10,
			},
			mockBehaviour: func(s *mockservice.MockILyricService, uuid pgtype.UUID, pag crud.Pagination) {
				s.EXPECT().ListByAudioPag(uuid, pag).Return([]dto.LyricRead{{
					UUID:      pgtype.UUID{Valid: true, Bytes: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}},
					AudioUUID: pgtype.UUID{Valid: true, Bytes: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
					Order:     0,
					Text:      "text1",
					CreatedAt: pgtype.Timestamptz{Valid: true},
					UpdatedAt: pgtype.Timestamptz{Valid: true},
				},
					{
						UUID:      pgtype.UUID{Valid: true, Bytes: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3}},
						AudioUUID: pgtype.UUID{Valid: true, Bytes: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
						Order:     1,
						Text:      "text2",
						CreatedAt: pgtype.Timestamptz{Valid: true},
						UpdatedAt: pgtype.Timestamptz{Valid: true},
					},
				}, nil)
			},
			expectedCode: 200,
			expectedBody: `{
									"data": [
										{"audio_uuid":"00000000-0000-0000-0000-000000000001", "created_at":"0001-01-01T00:00:00Z", "order":0, "text":"text1", "updated_at":"0001-01-01T00:00:00Z", "uuid":"00000000-0000-0000-0000-000000000002"},
										{"audio_uuid":"00000000-0000-0000-0000-000000000001", "created_at":"0001-01-01T00:00:00Z", "order":1, "text":"text2", "updated_at":"0001-01-01T00:00:00Z", "uuid":"00000000-0000-0000-0000-000000000003"}
									],
									"message":"lyrics got correctly",
									"next_pagination": {"limit":35, "offset":45}
								 }`,
			bodyMustContain: "",
		},
		{
			name:          "400_invalid_input_uuid",
			inputPathUUID: "00000000-0000-0000-0000-00000000000y",
			inputQuery:    "",
			inputUUID:     pgtype.UUID{Valid: true},
			inputPag: crud.Pagination{
				Limit:  0,
				Offset: 50,
			},
			mockBehaviour: func(s *mockservice.MockILyricService, uuid pgtype.UUID, pag crud.Pagination) {
			},
			expectedCode:    400,
			expectedBody:    `{"error":"encoding/hex: invalid byte: U+0079 'y'", "message":"invalid uuid in path param"}`,
			bodyMustContain: "",
		},
		{
			name:          "200_no_rows",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputQuery:    "",
			inputUUID:     pgtype.UUID{Valid: true},
			inputPag: crud.Pagination{
				Limit:  50,
				Offset: 0,
			},
			mockBehaviour: func(s *mockservice.MockILyricService, uuid pgtype.UUID, pag crud.Pagination) {
				s.EXPECT().ListByAudioPag(uuid, pag).Return(nil, pgx.ErrNoRows)
			},
			expectedCode:    200,
			expectedBody:    `{"data": [], "message":"no rows find", "next_pagination": {"limit":50, "offset":50}}`,
			bodyMustContain: "",
		},
		{
			name:          "500_unknown_error",
			inputPathUUID: "00000000-0000-0000-0000-000000000000",
			inputQuery:    "",
			inputUUID:     pgtype.UUID{Valid: true},
			inputPag: crud.Pagination{
				Limit:  50,
				Offset: 0,
			},
			mockBehaviour: func(s *mockservice.MockILyricService, uuid pgtype.UUID, pag crud.Pagination) {
				s.EXPECT().ListByAudioPag(uuid, pag).Return(nil, errors.New("unknown error"))
			},
			expectedCode:    500,
			expectedBody:    `{"error":"unknown error", "message":"error on list audio lyrics"}`,
			bodyMustContain: "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()
			lyricService := mockservice.NewMockILyricService(c)
			testCase.mockBehaviour(lyricService, testCase.inputUUID, testCase.inputPag)

			services := service.Service{Lyric: lyricService}
			handler := NewHandler(Deps{
				Service: services,
				Logger:  logging.GetLoggerTest(),
				Config: &config.Config{
					Server: config.Server{
						PagLimit: 50,
					},
				},
			})

			//Test server
			r := httprouter.New()
			r.GET("/audios/:uuid/lyrics", handler.audioLyricsList)

			//http test
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/audios/"+testCase.inputPathUUID+"/lyrics"+testCase.inputQuery, nil)

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != "" {
				if testCase.expectedBody == "nil" {
					testCase.expectedBody = "{}"
				}
				assert.JSONEq(t, testCase.expectedBody, w.Body.String())
			}
			if testCase.bodyMustContain != "" {
				assert.Contains(t, w.Body.String(), testCase.bodyMustContain)
			}
		})
	}
}
