package test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type APITestCase struct {
	Path             string
	Method           string
	Body             string
	ExpectedStatus   int
	ExpectedResponse string
	Message          string
}

func Endpoint(asserts *assert.Assertions, r *gin.Engine, testCase APITestCase) {
	req, err := http.NewRequest(testCase.Method, testCase.Path, bytes.NewBufferString(testCase.Body))
	asserts.NoError(err, testCase.Message)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	// Assert Response Status
	asserts.Equal(testCase.ExpectedStatus, resp.StatusCode, testCase.Message)

	// Assert Response Body
	if testCase.ExpectedResponse != "" {
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		asserts.Nil(err)
		asserts.Regexp(testCase.ExpectedResponse, string(body), testCase.Message)
	}
}
