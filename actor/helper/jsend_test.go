package helper_test

import (
	"sirius/actor/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Test string
}

func TestSuccessResponse(t *testing.T) {
	resp := helper.SuccessResponse()

	assert.Equal(t, "success", resp.Status)
}

func TestObjectResponse(t *testing.T) {
	data := testStruct{
		Test: "test",
	}

	resp := helper.ObjectResponse(data)

	assert.Equal(t, "success", resp.Status)
	assert.ObjectsAreEqualValues(data, resp.Data)
}

func TestFailResponse(t *testing.T) {
	resp := helper.FailResponse("failed")

	assert.Equal(t, "failed", resp.Status)
	assert.Equal(t, "failed", resp.Message)
}
