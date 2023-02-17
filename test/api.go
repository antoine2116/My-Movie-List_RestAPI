package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type APITestCase struct {
	Path           string
	Method         string
	Body           string
	ExpectedStatus int
}

func Endpoint(asserts *assert.Assertions, r *gin.Engine, testCase APITestCase) {
	req, err := http.NewRequest(testCase.Method, testCase.Path, bytes.NewBufferString(testCase.Body))
	asserts.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()

	asserts.Equal(testCase.ExpectedStatus, resp.StatusCode)
}
