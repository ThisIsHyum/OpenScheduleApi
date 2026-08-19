package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestAddLessonsBad(t *testing.T) {
	t.Cleanup(truncateAll)

	app := SetupApp(testDB, t)
	token, _ := createTestCollege(t, app)

	tests := []struct{ name, body string }{
		{"invalid json", "{"},
		{"empty list", `[]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := PostParserReq("/parser/lessons", token, tt.body)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		})
	}
}

func TestAddLessons(t *testing.T) {
	t.Cleanup(truncateAll)
	testDate := dto.Date{Time: time.Now()}
	app := SetupApp(testDB, t)
	token, college := createTestCollege(t, app)

	testGroupReq := dto.UpdateGroupsRequest{
		CampusID: college.Campuses[0].ID,
		StudentGroupNames: []string{
			"TEST_GROUP",
		},
	}
	b, err := json.Marshal(testGroupReq)
	assert.Nil(t, err)
	req := PostParserReq("/parser/groups", token, string(b))

	_, err = app.Test(req)
	assert.Nil(t, err)
	req = GetReq(fmt.Sprintf("/campuses/%d/groups", testGroupReq.CampusID))

	resp, err := app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	groups := parseJSON[[]dto.StudentGroupResponse](t, resp)
	group := groups[0]

	testLessons := []dto.Lesson{
		{LessonID: 1, Title: "TestLesson1", Cabinet: "TestCabinet1", Date: testDate, Teacher: "TestTeacher1", Order: 1, StudentGroupID: group.ID},
		{LessonID: 2, Title: "TestLesson2", Cabinet: "TestCabinet2", Date: testDate, Teacher: "TestTeacher2", Order: 2, StudentGroupID: group.ID},
		{LessonID: 3, Title: "TestLesson3", Cabinet: "TestCabinet3", Date: testDate, Teacher: "TestTeacher3", Order: 3, StudentGroupID: group.ID},
		{LessonID: 4, Title: "TestLesson4", Cabinet: "TestCabinet4", Date: testDate, Teacher: "TestTeacher4", Order: 4, StudentGroupID: group.ID},
		{LessonID: 5, Title: "TestLesson5", Cabinet: "TestCabinet5", Date: testDate, Teacher: "TestTeacher5", Order: 5, StudentGroupID: group.ID},
	}

	b, err = json.Marshal(testLessons)
	assert.Nil(t, err)
	req = PostParserReq("/parser/lessons", token, string(b))
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	req = GetReq(fmt.Sprintf("/groups/%d/schedules?date=%s", group.ID, testDate.Format("02-01-2006")))
	resp, err = app.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	type ScheduleLessonTestResponse struct {
		Title     string `json:"title"`
		Cabinet   string `json:"cabinet"`
		Teacher   string `json:"teacher"`
		Order     uint   `json:"order"`
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	}

	type ScheduleTestResponse struct {
		GroupID uint                         `json:"groupId"`
		Date    time.Time                    `json:"date"`
		Lessons []ScheduleLessonTestResponse `json:"lessons"`
	}

	schedules := parseJSON[[]ScheduleTestResponse](t, resp)
	assert.Equal(t, 1, len(schedules))
	assert.Equal(t, group.ID, schedules[0].GroupID)

	expectedTitles := []string{"TestLesson1", "TestLesson2", "TestLesson3", "TestLesson4", "TestLesson5"}
	actualTitles := make([]string, len(schedules[0].Lessons))
	for i, l := range schedules[0].Lessons {
		actualTitles[i] = l.Title
	}
	assert.ElementsMatch(t, expectedTitles, actualTitles)
}

func extractNames(groups []dto.StudentGroupResponse) []string {
	names := make([]string, len(groups))
	for i, g := range groups {
		names[i] = g.Name
	}
	return names
}
