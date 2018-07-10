package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dynastymasra/shajaro/actor/infrastructure/web/controller"

	"context"

	"github.com/dynastymasra/shajaro/actor/config"

	"github.com/dynastymasra/shajaro/actor/infrastructure/provider"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HealthzControllerSuite struct {
	suite.Suite
	context.Context
}

func Test_HealthzController(t *testing.T) {
	suite.Run(t, new(HealthzControllerSuite))
}

func (s *HealthzControllerSuite) SetupSuite() {
	config.InitConfig()
	provider.ConnectSQL()
}

func (s *HealthzControllerSuite) TearDownSuite() {
	db, _ := provider.ConnectSQL()
	provider.CloseDB(db)
}

func (s *HealthzControllerSuite) Test_PingController_Success() {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/v1/healthz", nil)

	controller.HealthzController(w, r)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}
