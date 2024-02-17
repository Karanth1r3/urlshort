package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Karanth1r3/url-short-learn/internal/httpapi/server/handlers/url/save"
	"github.com/Karanth1r3/url-short-learn/internal/util/logger/handlers/slogdiscard"
	"github.com/Karanth1r3/url-short-learn/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	testCase := []struct {
		name      string
		inputBody string
		inputURL  string
		alias     string
		respError string
		mockError error
	}{
		{
			name:     "OK",
			inputURL: "https://google.com",
			alias:    "test_alias",
		},
		{
			name:     "Empty alias",
			alias:    "",
			inputURL: "https://google.com",
		},
		{
			name:     "Empty URL",
			alias:    "some_alias",
			inputURL: "",
		},
		{
			name:      "Invalid URL",
			inputURL:  "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			alias:     "test_alias",
			inputURL:  "https://google.com",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, testC := range testCase {
		tc := testC

		urlSaverMock := mocks.NewURLSaver(t)

		if tc.respError == "" || tc.mockError != nil {
			urlSaverMock.On("SaveURL", tc.inputURL, mock.AnythingOfType("string")).
				Return(errors.New("unexpected")).
				Once()
		}

		handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)

		input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.inputURL, tc.alias)

		req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusOK)

		body := rr.Body.String()

		var resp save.Response

		require.NoError(t, json.Unmarshal([]byte(body), &resp))

		require.Equal(t, tc.respError, resp.Error)

		//TODO: Add more cheecks
	}
}
