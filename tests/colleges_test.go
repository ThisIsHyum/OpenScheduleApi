package tests

import (
	"fmt"
	"testing"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/dto"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func createTestCollege(t *testing.T, app *fiber.App) (string, dto.CollegeResponse) {
	req := AdminPostReq("/admin/parser", `{"collegeName":"test","campusNames":["a","b","c"]}`)
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	tokenResp := parseJSON[dto.NewParserResponse](t, resp)

	req = GetReq("/colleges")
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	colleges := parseJSON[[]dto.CollegeResponse](t, resp)

	return tokenResp.Token, colleges[0]
}

func TestGetColleges(t *testing.T) {
	t.Cleanup(truncateAll)
	app := SetupApp(testDB, t)

	createTestCollege(t, app)

	req := GetReq("/colleges")
	resp, err := app.Test(req)
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

func TestGetCollegesByID(t *testing.T) {
	t.Cleanup(truncateAll)
	app := SetupApp(testDB, t)

	_, college := createTestCollege(t, app)

	req := GetReq(fmt.Sprintf("/colleges/%d", college.ID))
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	college = parseJSON[dto.CollegeResponse](t, resp)

	assert.Equal(t, uint(1), college.ID)
	assert.Equal(t, "test", college.Name)
	assert.Equal(t, 3, len(college.Campuses))
	campus_names := map[string]struct{}{"a": {}, "b": {}, "c": {}}
	for _, campus := range college.Campuses {
		_, ok := campus_names[campus.Name]
		assert.True(t, ok)
	}

	req = GetReq("/colleges/999999999")
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	req = GetReq("/colleges/test")
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestGetCollegesByName(t *testing.T) {
	t.Cleanup(truncateAll)
	app := SetupApp(testDB, t)
	createTestCollege(t, app)

	req := GetReq("/colleges?name=test")
	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	colleges := parseJSON[[]dto.CollegeResponse](t, resp)

	assert.Equal(t, 1, len(colleges))
	assert.Equal(t, "test", colleges[0].Name)
	assert.Equal(t, 3, len(colleges[0].Campuses))
	campus_names := map[string]struct{}{"a": {}, "b": {}, "c": {}}
	for _, campus := range colleges[0].Campuses {
		_, ok := campus_names[campus.Name]
		assert.True(t, ok)
	}

	req = GetReq("/colleges?name=t")
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	colleges = parseJSON[[]dto.CollegeResponse](t, resp)
	assert.Equal(t, 0, len(colleges))
}
