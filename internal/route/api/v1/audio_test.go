package v1

import (
	"eMobile/internal/config"
	"eMobile/internal/crud"
	"eMobile/internal/dto"
	"eMobile/internal/service"
	mockservice "eMobile/internal/service/mocks"
	"eMobile/pkg/logging"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"strings"
	"testing"
)

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
			expectedBody:    "nil",
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
			req := httptest.NewRequest("GET", "/audios/"+testCase.inputPathUUID, nil)

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
