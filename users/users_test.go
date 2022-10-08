package users

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNonPostError(t *testing.T) {
	req := httptest.NewRequest("GET", "localhost:8080", nil)
	w := httptest.NewRecorder()

	HandleUserRequest(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Error("Expected status code to indicate the status is not allowed")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Erro: %v", err)
	}
	if strings.TrimSpace(string(body)) != ErrorMethodNotSupported {
		t.Errorf("Body was %s, expected %s", string(body), ErrorMethodNotSupported)
	}
}

// Also tests the handleUserInputs function
func TestHandleUserRequestPostResponses(t *testing.T) {
	tests := []struct {
		json                 string
		expectedResponseCode int
		expectedResponseBody string
	}{
		{"", http.StatusNoContent, ""},
		{"[]", http.StatusOK, "[]"},
		{"this is not json", http.StatusBadRequest, ErrorParsingInput},
		// Treats partial objects as bad requests
		{
			`[{"date_of_birth": "1983-05-12"}]`,
			http.StatusBadRequest,
			ErrorParsingInput,
		},
		// cannot process requests with malformed dates
		{
			`[{"user_id": 1, "name": "Joe Smith", "date_of_birth": "1983-05-124", "created_on": 1642612034 }]`,
			http.StatusInternalServerError,
			ErrorProcessingInput,
		},
		// Can parse the expected values
		{
			`[{"user_id": 1, "name": "Joe Smith", "date_of_birth": "1983-05-12", "created_on": 1642612034 }]`,
			http.StatusOK,
			`[{"user_id":1,"name":"Joe Smith","weekday_of_birth":"Thursday","created_on":"2022-01-19T12:07:14-05:00"}]`,
		},
		{
			`[
				{"user_id": 1, "name": "Joe Smith", "date_of_birth": "1983-05-12", "created_on": 1642612034 },
				{"user_id": 2, "name": "Jane Smith", "date_of_birth": "1984-05-12", "created_on": 1642612035 },
				{"user_id": 3, "name": "Doe Smith", "date_of_birth": "1985-05-12", "created_on": 1642612036 }
			]`,
			http.StatusOK,
			`[{"user_id":1,"name":"Joe Smith","weekday_of_birth":"Thursday","created_on":"2022-01-19T12:07:14-05:00"},` +
				`{"user_id":2,"name":"Jane Smith","weekday_of_birth":"Saturday","created_on":"2022-01-19T12:07:15-05:00"},` +
				`{"user_id":3,"name":"Doe Smith","weekday_of_birth":"Sunday","created_on":"2022-01-19T12:07:16-05:00"}]`,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("HandleUserRequest=%d", i), func(t *testing.T) {
			req := httptest.NewRequest("POST", "localhost:8080", strings.NewReader(test.json))
			w := httptest.NewRecorder()

			HandleUserRequest(w, req)

			resp := w.Result()

			if resp.StatusCode != test.expectedResponseCode {
				t.Errorf("Expected status code to be %d, but was %d", test.expectedResponseCode, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			if strings.TrimSpace(string(body)) != test.expectedResponseBody {
				t.Errorf("Body was %s, expected %s", string(body), test.expectedResponseBody)
			}
		})
	}
}

func TestProcessUserInputs(t *testing.T) {
	tests := []struct {
		json               string
		expectsError       bool
		expectedUserInputs []UserInput
	}{
		// Throws error when json is parsed
		{"this is not json", true, nil},
		// Throws error when not an array is parsed
		{"{}", true, nil},
		// Can parse partial values
		{
			`[{"date_of_birth": "1983-05-12", "created_on": 1642612034 }]`,
			true,
			nil,
		},
		{"[{}]", true, nil},
		// Can parse the expected values
		{"", false, nil},
		{"[]", false, []UserInput{}},
		{
			`[{"user_id": 1, "name": "Joe Smith", "date_of_birth": "1983-05-12", "created_on": 1642612034 }]`,
			false,
			[]UserInput{
				baseUserInputGen(1, "Joe Smith", "1983-05-12", 1642612034),
			},
		},
		{
			`[
				{"user_id": 1, "name": "Joe Smith", "date_of_birth": "1983-05-12", "created_on": 1642612034 },
				{"user_id": 2, "name": "Jane Smith", "date_of_birth": "1984-05-12", "created_on": 1642612035 },
				{"user_id": 3, "name": "Doe Smith", "date_of_birth": "1985-05-12", "created_on": 1642612036 }
			]`,
			false,
			[]UserInput{
				baseUserInputGen(1, "Joe Smith", "1983-05-12", 1642612034),
				baseUserInputGen(2, "Jane Smith", "1984-05-12", 1642612035),
				baseUserInputGen(3, "Doe Smith", "1985-05-12", 1642612036),
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("processUserInputs=%d", i), func(t *testing.T) {
			testJsonReadCloser := io.NopCloser(strings.NewReader(test.json))
			userInputs, err := processUserInputs(&testJsonReadCloser)
			if err != nil && !test.expectsError {
				t.Errorf("Error: %v", err)
			}
			if !reflect.DeepEqual(userInputs, test.expectedUserInputs) {
				t.Errorf("Received: %v, Expected: %v", userInputs, test.expectedUserInputs)
			}
		})
	}
}

func TestTransformUserInputs(t *testing.T) {
	tests := []struct {
		userInputs          []UserInput
		expectsError        bool
		expectedUserOutputs []UserOutput
	}{
		// Throws error when it cannot parse the given struct
		// due to an invalid date of birth
		{[]UserInput{baseUserInputGen(1, "Joe Smith", "1983-05-124", 1642612034)}, true, nil},
		// Can parse valid structures,
		{[]UserInput{}, false, []UserOutput{}},
		{
			[]UserInput{baseUserInputGen(1, "Joe Smith", "1983-05-12", 1642612034)},
			false,
			[]UserOutput{{UserId: 1, Name: "Joe Smith", WeekdayOfBirth: "Thursday", CreatedOn: "2022-01-19T12:07:14-05:00"}},
		},
		{
			[]UserInput{
				baseUserInputGen(1, "Solomon Grundy", "1983-05-09", 1642612034),
				baseUserInputGen(2, "Jane Smith", "1984-05-10", 1642612035),
				baseUserInputGen(3, "Doe Smith", "1985-05-11", 1642612036),
			},
			false,
			[]UserOutput{
				{UserId: 1, Name: "Solomon Grundy", WeekdayOfBirth: "Monday", CreatedOn: "2022-01-19T12:07:14-05:00"},
				{UserId: 2, Name: "Jane Smith", WeekdayOfBirth: "Thursday", CreatedOn: "2022-01-19T12:07:15-05:00"},
				{UserId: 3, Name: "Doe Smith", WeekdayOfBirth: "Saturday", CreatedOn: "2022-01-19T12:07:16-05:00"},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("transformUserInputs=%d", i), func(t *testing.T) {
			userOutputs, err := transformUserInputs(test.userInputs)
			if err != nil && !test.expectsError {
				t.Errorf("Error: %v", err)
			}
			if !reflect.DeepEqual(userOutputs, test.expectedUserOutputs) {
				t.Errorf("Received: %v, Expected: %v", userOutputs, test.expectedUserOutputs)
			}
		})
	}
}

func baseUserInputGen(id int, name string, dateOfBirth string, createdOn int64) UserInput {
	return UserInput{
		UserId:      &id,
		Name:        &name,
		DateOfBirth: &dateOfBirth,
		CreatedOn:   &createdOn,
	}
}

func partialUserInputGen(dateOfBirth string) UserInput {
	return UserInput{
		DateOfBirth: &dateOfBirth,
	}
}
