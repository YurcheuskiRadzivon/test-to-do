package admin_test

import (
	"net/http"
	"testing"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/admin"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/admin/mock"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go.uber.org/mock/gomock"
)

func TestGetUsersIvnalid_Unit(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	am := mock.NewMockAuthManager(cntrl)
	em := mock.NewMockEncryptManager(cntrl)
	us := mock.NewMockUserService(cntrl)
	am.EXPECT().GetUserID(gomock.Any()).Return(1, nil)

	ac := admin.NewAdminControl(us, am, em)

	testApp := fiber.New()

	testCtx := &fasthttp.RequestCtx{}

	testReq := &fasthttp.Request{}
	testReq.Header.SetMethod("GET")
	testReq.SetRequestURI("/test")

	testFiberCtx := testApp.AcquireCtx(testCtx)
	defer testApp.ReleaseCtx(testFiberCtx)

	_ = ac.GetUsers(testFiberCtx)
	assert.Equal(t, http.StatusBadRequest, testFiberCtx.Response().StatusCode())
}

func TestGetUsersValid_Unit(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	am := mock.NewMockAuthManager(cntrl)
	em := mock.NewMockEncryptManager(cntrl)
	us := mock.NewMockUserService(cntrl)
	am.EXPECT().GetUserID(gomock.Any()).Return(0, nil)
	us.EXPECT().GetUsers(gomock.Any()).Return([]entity.User{}, nil)

	ac := admin.NewAdminControl(us, am, em)

	testApp := fiber.New()

	testCtx := &fasthttp.RequestCtx{}

	testReq := &fasthttp.Request{}
	testReq.Header.SetMethod("GET")
	testReq.SetRequestURI("/test")

	testFiberCtx := testApp.AcquireCtx(testCtx)
	defer testApp.ReleaseCtx(testFiberCtx)

	_ = ac.GetUsers(testFiberCtx)
	assert.Equal(t, http.StatusOK, testFiberCtx.Response().StatusCode())
}
