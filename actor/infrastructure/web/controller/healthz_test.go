package controller_test

import (
	"net/http"
	"net/http/httptest"
	"sirius/actor/infrastructure/web/controller"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingController_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.HEAD("/v1/ping", controller.HealthzController)

	req, _ := http.NewRequest(http.MethodHead, "/v1/ping", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
