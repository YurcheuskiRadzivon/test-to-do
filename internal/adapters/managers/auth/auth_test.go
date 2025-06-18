package authmanage_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	authmanage "github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/managers/auth"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

const (
	testSecretKey = "secret_key"
)

func TestGetUserIdValid_Integration(t *testing.T) {
	testJWTS := jwtservice.New(testSecretKey)
	am := authmanage.NewAuthManage(testJWTS)

	cases := []struct {
		name     string
		inHeader string
		out      struct {
			id int
		}
	}{
		{
			name:     "test_1",
			inHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowfQ.d26rEQQyVkR9ByhfRukCFrKDILtNGQB1qxOxohXsNtY",
			out: struct {
				id int
			}{
				id: 0,
			},
		},
		{
			name:     "test_2",
			inHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyfQ.2hcCaw1GXwpCuVBxvWOsuZeDrKz9eBI5mMR8WbjUzXg",
			out: struct {
				id int
			}{
				id: 2,
			},
		},
	}

	for _, c := range cases {
		testApp := fiber.New()

		testReq := &fasthttp.Request{}
		testReq.Header.SetMethod("GET")
		testReq.SetRequestURI("/test")
		testReq.Header.Set(jwtservice.HeaderAuthorization, c.inHeader)

		testCtx := &fasthttp.RequestCtx{}
		testCtx.Init(testReq, nil, nil)

		testFiberCtx := testApp.AcquireCtx(testCtx)
		defer testApp.ReleaseCtx(testFiberCtx)

		userID, err := am.GetUserID(testFiberCtx)

		assert.NoError(t, err, c.name)
		assert.Equal(t, c.out.id, userID)
	}

}

func TestGetUserIdInvalid_Integration(t *testing.T) {
	testJWTS := jwtservice.New(testSecretKey)
	am := authmanage.NewAuthManage(testJWTS)

	cases := []struct {
		name     string
		inHeader string
		out      struct {
			err string
		}
	}{
		{
			name:     "test_1",
			inHeader: "QB1qxOxohXsNtY",
			out: struct {
				err string
			}{
				err: jwtservice.ErrInvalidOrExpiredToken,
			},
		},
	}

	for _, c := range cases {
		testApp := fiber.New()

		testReq := &fasthttp.Request{}
		testReq.Header.SetMethod("GET")
		testReq.SetRequestURI("/test")
		testReq.Header.Set(jwtservice.HeaderAuthorization, c.inHeader)

		testCtx := &fasthttp.RequestCtx{}
		testCtx.Init(testReq, nil, nil)

		testFiberCtx := testApp.AcquireCtx(testCtx)
		defer testApp.ReleaseCtx(testFiberCtx)

		_, err := am.GetUserID(testFiberCtx)

		assert.EqualError(t, err, jwtservice.ErrInvalidOrExpiredToken)
	}

}

func TestValidateValid_Integration(t *testing.T) {
	testJWTS := jwtservice.New(testSecretKey)
	am := authmanage.NewAuthManage(testJWTS)

	cases := []struct {
		name     string
		inHeader string
	}{
		{
			name:     "test_1",
			inHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowfQ.d26rEQQyVkR9ByhfRukCFrKDILtNGQB1qxOxohXsNtY",
		},
		{
			name:     "test_2",
			inHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyfQ.2hcCaw1GXwpCuVBxvWOsuZeDrKz9eBI5mMR8WbjUzXg",
		},
	}

	for _, c := range cases {
		testApp := fiber.New()

		testReq := &fasthttp.Request{}
		testReq.Header.SetMethod("GET")
		testReq.SetRequestURI("/test")
		testReq.Header.Set(jwtservice.HeaderAuthorization, c.inHeader)

		testCtx := &fasthttp.RequestCtx{}
		testCtx.Init(testReq, nil, nil)

		testFiberCtx := testApp.AcquireCtx(testCtx)
		defer testApp.ReleaseCtx(testFiberCtx)

		err := am.Validate(testFiberCtx)

		assert.NoError(t, err, c.name)
	}
}

func TestValidateInvalid_Integration(t *testing.T) {
	testJWTS := jwtservice.New(testSecretKey)
	am := authmanage.NewAuthManage(testJWTS)

	cases := []struct {
		name     string
		inHeader string
		out      struct {
			err string
		}
	}{
		{
			name:     "test_1",
			inHeader: "QB1qxOxohXsNtY",
			out: struct {
				err string
			}{
				err: "INVALID_OR_EXPIRED_TOKEN",
			},
		},
	}

	for _, c := range cases {
		testApp := fiber.New()

		testReq := &fasthttp.Request{}
		testReq.Header.SetMethod("GET")
		testReq.SetRequestURI("/test")
		testReq.Header.Set(jwtservice.HeaderAuthorization, c.inHeader)

		testCtx := &fasthttp.RequestCtx{}
		testCtx.Init(testReq, nil, nil)

		testFiberCtx := testApp.AcquireCtx(testCtx)
		defer testApp.ReleaseCtx(testFiberCtx)

		err := am.Validate(testFiberCtx)
		assert.EqualError(t, err, c.out.err)
	}
}

func TestCreateToken_Integration(t *testing.T) {
	testJWTS := jwtservice.New(testSecretKey)
	am := authmanage.NewAuthManage(testJWTS)

	cases := []struct {
		name string
		in   int
		out  string
	}{
		{
			name: "test_1",
			in:   0,
			out:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowfQ.d26rEQQyVkR9ByhfRukCFrKDILtNGQB1qxOxohXsNtY",
		},
	}

	for _, c := range cases {
		testApp := fiber.New()

		testReq := &fasthttp.Request{}
		testReq.Header.SetMethod("GET")
		testReq.SetRequestURI("/test")

		testCtx := &fasthttp.RequestCtx{}
		testCtx.Init(testReq, nil, nil)

		testFiberCtx := testApp.AcquireCtx(testCtx)
		defer testApp.ReleaseCtx(testFiberCtx)

		err := am.CreateAuthResponse(testFiberCtx, c.in)
		assert.NoError(t, err, c.name)

		assert.Equal(t, http.StatusOK, testFiberCtx.Response().StatusCode())

		var resp response.LoginResponse
		if err := json.Unmarshal(testFiberCtx.Response().Body(), &resp); err != nil {
			t.Fatalf("Cannot parse json: %s", err)
		}

		assert.NotEmpty(t, resp.Token)

		assert.Equal(t, c.out, resp.Token)
	}
}
