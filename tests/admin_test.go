package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func PostReq(endpoint, body string) *http.Request {
	req := httptest.NewRequest(
		http.MethodPost,
		endpoint,
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func AdminPostReqWithToken(endpoint, body, token string) *http.Request {
	req := PostReq(endpoint, body)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func AdminPostReq(endpoint, body string) *http.Request {
	return AdminPostReqWithToken(endpoint, body, TestAdminToken)
}

func AdminDeleteReq(endpoint string) *http.Request {
	req := httptest.NewRequest(http.MethodDelete, endpoint, nil)
	req.Header.Set("Authorization", "Bearer "+TestAdminToken)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func parseJSON[T any](t *testing.T, resp *http.Response) T {
	var v T
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(body, &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func GetReq(endpoint string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestAdminCreateParser(t *testing.T) {
	t.Cleanup(truncateAll)
	app := SetupApp(testDB, t)

	req := PostReq("/admin/parser", `{"collegeName":"test","campusNames":["a", "b", "c"]}`)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	req = AdminPostReqWithToken("/admin/parser", `{"collegeName":"test","campusNames":["a", "b", "c"]}`, "")

	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	req = AdminPostReqWithToken("/admin/parser", "", "kdosgzopPOGEW39ruf3ej9kd")
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	req = AdminPostReq("/admin/parser", `{"collegeName":"test","campusNa`)
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	req = AdminPostReq("/admin/parser", `{"collegeName":"test","campusNames":["a", "b", "c"]}`)
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	req = AdminPostReq("/admin/parser", `{"collegeName":"test","campusNames":["a", "b", "c"]}`)
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)

	req = GetReq("/colleges")

	resp, err = app.Test(req)
	assert.Nil(t, err)

	colleges := parseJSON[[]dto.CollegeResponse](t, resp)

	assert.Equal(t, 1, len(colleges))
	assert.Equal(t, "test", colleges[0].Name)
	assert.Equal(t, 3, len(colleges[0].Campuses))
	campus_names := map[string]struct{}{"a": {}, "b": {}, "c": {}}
	for _, campus := range colleges[0].Campuses {
		_, ok := campus_names[campus.Name]
		assert.True(t, ok)
	}
}

func TestAdminDeleteParser(t *testing.T) {
	t.Cleanup(truncateAll)
	app := SetupApp(testDB, t)

	req := AdminPostReq("/admin/parser", `{"collegeName":"test","campusNames":["a", "b", "c"]}`)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	req = GetReq("/colleges")

	resp, err = app.Test(req)
	assert.Nil(t, err)

	colleges := parseJSON[[]dto.CollegeResponse](t, resp)

	req = AdminDeleteReq(fmt.Sprintf("/admin/parser/%d", colleges[0].ID))
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	req = GetReq(fmt.Sprintf("/colleges/%d", colleges[0].ID))
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	req = AdminDeleteReq(fmt.Sprintf("/admin/parser/%d", colleges[0].ID))
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	for _, id := range []uint{colleges[0].ID, 50} {
		req = AdminDeleteReq(fmt.Sprintf("/admin/parser/%d", id))
		resp, err = app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	}
	req = AdminDeleteReq("/admin/parser/test")
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
