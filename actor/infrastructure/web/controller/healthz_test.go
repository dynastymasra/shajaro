package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"shajaro/actor/infrastructure/web/controller"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingController_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.HEAD("/v1/healthz", controller.HealthzController)

	req, _ := http.NewRequest(http.MethodHead, "/v1/healthz", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
