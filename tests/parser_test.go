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

func ParserReq(endpoint, token, method string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, endpoint, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}
func GetParserReq(endpoint, token string) *http.Request {
	return ParserReq(endpoint, token, http.MethodGet, nil)
}
func PostParserReq(endpoint, token, body string) *http.Request {
	return ParserReq(endpoint, token, http.MethodPost, strings.NewReader(body))
}
func DeleteParserReq(endpoint, token string) *http.Request {
	return ParserReq(endpoint, token, http.MethodDelete, nil)
}

func TestGetParser(t *testing.T) {
	t.Cleanup(truncateAll)

	app := SetupApp(testDB, t)
	token, _ := createTestCollege(t, app)

	req := GetParserReq("/parser", token)

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	collegeId := parseJSON[dto.GetParserResponse](t, resp).CollegeID
	req = GetReq(fmt.Sprintf("/colleges/%d", collegeId))
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	college := parseJSON[dto.CollegeResponse](t, resp)
	assert.Equal(t, collegeId, college.ID)
	assert.Equal(t, "test", college.Name)

	req = GetParserReq("/parser", "test_wrong_token_key")
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	req = GetParserReq("/parser", "")
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	req = GetReq("/parser")
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestUpdateGroupsBad(t *testing.T) {
	t.Cleanup(truncateAll)

	app := SetupApp(testDB, t)

	token, _ := createTestCollege(t, app)

	tests := []struct{ name, body string }{
		{"invalid json", "{"},
		{"empty struct", `{}`},
		{"missing groups", `{"campusId":1}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := PostParserReq("/parser/groups", token, tt.body)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		})
	}
}
func TestUpdateGroups(t *testing.T) {
	t.Cleanup(truncateAll)

	app := SetupApp(testDB, t)
	token, college := createTestCollege(t, app)

	groupsReqTests := []struct {
		req      dto.UpdateGroupsRequest
		expected []string
	}{
		{
			req: dto.UpdateGroupsRequest{
				CampusID: college.Campuses[0].ID,
				StudentGroupNames: []string{
					"CAMPUS0_TEST1", "CAMPUS0_TesT2", "CAMPUS0_TESt3",
				},
			},
			expected: []string{
				"CAMPUS0_TEST1", "CAMPUS0_TesT2", "CAMPUS0_TESt3",
			},
		},
		{
			req: dto.UpdateGroupsRequest{
				CampusID: college.Campuses[1].ID,
				StudentGroupNames: []string{
					"CAMPUS2_TEST1", "CAMPUS2_TesT2", "CAMPUS2_TESt3",
				},
			},
			expected: []string{
				"CAMPUS2_TEST1", "CAMPUS2_TesT2", "CAMPUS2_TESt3",
			},
		},
		{
			req: dto.UpdateGroupsRequest{
				CampusID: college.Campuses[2].ID,
				StudentGroupNames: []string{
					"CAMPUS3_TEST1", "CAMPUS3_TesT2", "CAMPUS3_TESt3",
				},
			},
			expected: []string{
				"CAMPUS3_TEST1", "CAMPUS3_TesT2", "CAMPUS3_TESt3"},
		},
		{
			req: dto.UpdateGroupsRequest{
				CampusID: college.Campuses[1].ID,
				StudentGroupNames: []string{
					"CAMPUS2_TEST1",
					"CAMPUS2_TESt4",
				},
			},
			expected: []string{
				"CAMPUS2_TEST1", "CAMPUS2_TesT2",
				"CAMPUS2_TESt3", "CAMPUS2_TESt4",
			},
		},
	}

	for _, groupReqTest := range groupsReqTests {
		b, err := json.Marshal(groupReqTest.req)
		assert.Nil(t, err)
		req := PostParserReq("/parser/groups", token, string(b))
		resp, err := app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

		req = GetReq(fmt.Sprintf("/campuses/%d/groups", groupReqTest.req.CampusID))
		resp, err = app.Test(req)
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		groups := parseJSON[[]dto.StudentGroupResponse](t, resp)
		assert.ElementsMatch(t, groupReqTest.expected, extractNames(groups))
	}
}

func extractNames(groups []dto.StudentGroupResponse) []string {
	names := make([]string, len(groups))
	for i, g := range groups {
		names[i] = g.Name
	}
	return names
}
