package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"projectionist/consts"
	"projectionist/models"
	"projectionist/provider"
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	type args struct {
		user       models.User
		dbProvider provider.IDBProvider
	}
	tests := []struct {
		name             string
		args             args
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "New User valid",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test",
							Deleted:  0,
							Password: "test",
							Role:     0,
							Token:    "testToken",
						},
					},
				),
				user: models.User{
					Username: "test1",
					Password: "testPass",
					Role:     0,
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "New user created",
				"userID":  float64(2),
			},
			wantResponseCode: 200,
		},
		{
			name: "New User not valid: Invalid username - empty.",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test",
							Deleted:  0,
							Password: "test",
							Role:     0,
							Token:    "testToken",
						},
					},
				),
				user: models.User{
					Username: "",
					Password: "testPass",
					Role:     0,
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.InputDataInvalidResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "New User not valid: Invalid password - empty.",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test",
							Deleted:  0,
							Password: "test",
							Role:     0,
							Token:    "testToken",
						},
					},
				),
				user: models.User{
					Username: "test",
					Password: "",
					Role:     0,
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.InputDataInvalidResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "New User not valid: Invalid username - great than 255.",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test",
							Deleted:  0,
							Password: "test",
							Role:     0,
							Token:    "testToken",
						},
					},
				),
				user: models.User{
					Username: "testusertestusertestusertestusertestusertestusertetestusertestusertestusertestusertestuse" +
						"rtestusertetestusertestusertestusertestusertestusertestusertetestusertestusertestusertestuserte" +
						"stusertestusertetestusertestusertestusertestusertestusertestusertetestusertestusertestusertestu" +
						"sertestusertestuserte",
					Password: "",
					Role:     0,
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.InputDataInvalidResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "New User not valid: Invalid password - great than 500.",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test",
							Deleted:  0,
							Password: "test",
							Role:     0,
							Token:    "testToken",
						},
					},
				),
				user: models.User{
					Username: "test",
					Password: "testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10" +
						"testpasswordtestpasswordtestpasswordtestpassword10",
					Role: 0,
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.InputDataInvalidResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "New User not valid: Invalid role - not admin and not super admin.",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test",
							Deleted:  0,
							Password: "test",
							Role:     0,
							Token:    "testToken",
						},
					},
				),
				user: models.User{
					Username: "test",
					Password: "testPass",
					Role:     3,
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.InputDataInvalidResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "New User not created: this userID is exist",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test",
							Deleted:  0,
							Password: "test",
							Role:     0,
							Token:    "testToken",
						},
					},
				),
				user: models.User{
					Username: "test",
					Password: "testPass",
					Role:     0,
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": "A user with the same name already exists.",
			},
			wantResponseCode: 403,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			bData, err := json.Marshal(tt.args.user)
			if err != nil {
				t.Fatalf("Marshal request form:%v", err)
			}

			buffer := bytes.NewBuffer(bData)

			request, err := http.NewRequest(http.MethodPost, consts.UrlUserV1, buffer)
			if err != nil {
				t.Fatalf("New Request error: %v", err)
			}

			recorder := httptest.NewRecorder()
			handler := NewUser(tt.args.dbProvider)
			handler.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Body)
			if err != nil {
				t.Errorf("Read response error:%v", err)
			}

			if !reflect.DeepEqual(recorder.Code, tt.wantResponseCode) {
				t.Errorf("response code got %v wantResp %v", recorder.Code, tt.wantResponseCode)
			}

			var gotResp = map[string]interface{}{}

			err = json.Unmarshal(body, &gotResp)
			if err != nil {
				t.Errorf("Unmarshal response body error:%v", err)
			}

			for wantKey, wantValue := range tt.wantResponseBody {
				if !reflect.DeepEqual(wantValue, gotResp[wantKey]) {
					t.Errorf("NewUser() key %v, got %+v, want %+v", wantKey, gotResp[wantKey], wantValue)
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		dbProvider provider.IDBProvider
		urlVars    map[string]string
	}
	tests := []struct {
		name             string
		args             args
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "A user with the same name already exists.",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test",
							Role:     0,
							Password: "testPass",
							Deleted:  0,
						},
					},
				),
				urlVars: map[string]string{
					"id": "1",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "",
				"user": map[string]interface{}{
					"id":       float64(1),
					"username": "test",
					"role":     float64(0),
					"Password": "",
					"token":    "",
					"deleted":  float64(0),
				},
			},
			wantResponseCode: 200,
		},
		{
			name: consts.NotExistResp,
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test",
							Role:     0,
							Password: "testPass",
							Deleted:  0,
						},
					},
				),
				urlVars: map[string]string{
					"id": "2",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.NotExistResp,
			},
			wantResponseCode: 404,
		},
		{
			name: "user id parameter not number",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test",
							Role:     0,
							Password: "testPass",
							Deleted:  0,
						},
					},
				),
				urlVars: map[string]string{
					"id": "test",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.IdIsNotNumberResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "user id not exist in params",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test",
							Role:     0,
							Password: "testPass",
							Deleted:  0,
						},
					},
				),
				urlVars: map[string]string{},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.IdIsEmptyResp,
			},
			wantResponseCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, consts.UrlUserV1, nil)
			if err != nil {
				t.Fatalf("New Request error: %v", err)
			}

			request = mux.SetURLVars(request, tt.args.urlVars)

			recorder := httptest.NewRecorder()
			handler := GetUser(tt.args.dbProvider)
			handler.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Body)
			if err != nil {
				t.Errorf("Read response error:%v", err)
			}

			if !reflect.DeepEqual(recorder.Code, tt.wantResponseCode) {
				t.Errorf("response code got %v wantResp %v", recorder.Code, tt.wantResponseCode)
			}

			var gotResp = map[string]interface{}{}

			err = json.Unmarshal(body, &gotResp)
			if err != nil {
				t.Errorf("Unmarshal response body error:%v", err)
			}

			for wantKey, wantValue := range tt.wantResponseBody {
				if !reflect.DeepEqual(wantValue, gotResp[wantKey]) {
					t.Errorf("GetUser() key %v, got %+v, want %+v", wantKey, gotResp[wantKey], wantValue)
				}
			}
		})
	}
}

func TestGetUserList(t *testing.T) {
	type args struct {
		dbProvider provider.IDBProvider
		urlValues  map[string]string
	}
	tests := []struct {
		name             string
		args             args
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "users exists for page 1",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				urlValues: map[string]string{
					consts.PAGE_PARAM:  "1",
					consts.COUNT_PARAM: "10",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "",
				consts.KEY_USERS: []map[string]interface{}{
					{
						"id":       float64(1),
						"username": "test1",
						"role":     float64(0),
						"Password": "testPass1",
						"deleted":  float64(0),
					},
					{
						"id":       float64(2),
						"username": "test2",
						"role":     float64(0),
						"Password": "testPass2",
						"deleted":  float64(0),
					},
					{
						"id":       float64(3),
						"username": "test3",
						"role":     float64(0),
						"Password": "testPass3",
						"deleted":  float64(0),
					},
					{
						"id":       float64(4),
						"username": "test4",
						"role":     float64(0),
						"Password": "testPass4",
						"deleted":  float64(0),
					},
					{
						"id":       float64(5),
						"username": "test5",
						"role":     float64(0),
						"Password": "testPass5",
						"deleted":  float64(0),
					},
					{
						"id":       float64(6),
						"username": "test6",
						"role":     float64(0),
						"Password": "testPass6",
						"deleted":  float64(0),
					},
					{
						"id":       float64(7),
						"username": "test7",
						"role":     float64(0),
						"Password": "testPass7",
						"deleted":  float64(0),
					},
					{
						"id":       float64(8),
						"username": "test8",
						"role":     float64(0),
						"Password": "testPass8",
						"deleted":  float64(0),
					},
					{
						"id":       float64(9),
						"username": "test9",
						"role":     float64(0),
						"Password": "testPass9",
						"deleted":  float64(0),
					},
					{
						"id":       float64(10),
						"username": "test10",
						"role":     float64(0),
						"Password": "testPass10",
						"deleted":  float64(0),
					},
				},
			},
			wantResponseCode: 200,
		},
		{
			name: "users exists for page 2",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				urlValues: map[string]string{
					consts.PAGE_PARAM:  "2",
					consts.COUNT_PARAM: "10",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "",
				consts.KEY_USERS: []map[string]interface{}{
					{
						"id":       float64(11),
						"username": "test11",
						"role":     float64(0),
						"Password": "testPass11",
						"deleted":  float64(0),
					},
					{
						"id":       float64(12),
						"username": "test12",
						"role":     float64(0),
						"Password": "testPass12",
						"deleted":  float64(0),
					},
				},
			},
			wantResponseCode: 200,
		},
		{
			name: "users not exists for page 0",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				urlValues: map[string]string{
					consts.PAGE_PARAM:  "0",
					consts.COUNT_PARAM: "10",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "",
			},
			wantResponseCode: 200,
		},
		{
			name: "users not exists for page 4",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				urlValues: map[string]string{
					consts.PAGE_PARAM:  "4",
					consts.COUNT_PARAM: "10",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "",
			},
			wantResponseCode: 200,
		},
		{
			name: "page param not number",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				urlValues: map[string]string{
					consts.PAGE_PARAM:  "test",
					consts.COUNT_PARAM: "10",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.PageMustNumber,
			},
			wantResponseCode: 400,
		},
		{
			name: "count param not number",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				urlValues: map[string]string{
					consts.PAGE_PARAM:  "1",
					consts.COUNT_PARAM: "test",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.CountMustNumber,
			},
			wantResponseCode: 400,
		},
		{
			name: "count and page params is required",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				urlValues: map[string]string{},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.PageAndCountRequiredResp,
			},
			wantResponseCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request, err := http.NewRequest(http.MethodGet, consts.UrlUserV1, nil)
			if err != nil {
				t.Fatalf("New Request error: %v", err)
			}

			query := request.URL.Query()
			for paramKey, paramValue := range tt.args.urlValues {
				query.Add(paramKey, paramValue)
			}

			request.URL.RawQuery = query.Encode()

			recorder := httptest.NewRecorder()
			handler := GetUserList(tt.args.dbProvider)
			handler.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Body)
			if err != nil {
				t.Errorf("Read response error:%v", err)
			}

			if !reflect.DeepEqual(recorder.Code, tt.wantResponseCode) {
				t.Errorf("response code got %v wantResp %v", recorder.Code, tt.wantResponseCode)
			}

			var gotResp = map[string]interface{}{}

			err = json.Unmarshal(body, &gotResp)
			if err != nil {
				t.Errorf("Unmarshal response body error:%v", err)
			}

			for wantKey, wantValue := range tt.wantResponseBody {
				if wantKey == consts.KEY_USERS {
					wantUsers := wantValue.([]map[string]interface{})
					iGotUsers, ok := gotResp[consts.KEY_USERS]
					if !ok || iGotUsers == nil {
						break
					}

					gotUsers := iGotUsers.([]interface{})

					if len(wantUsers) != len(gotUsers) {
						t.Errorf("count got users %v - count want users %v", len(gotUsers), len(wantUsers))
						break
					}

					for i, wantUser := range wantUsers {
						gotUser := gotUsers[i].(map[string]interface{})
						for wantUserKey, wantUserValue := range wantUser {
							if !reflect.DeepEqual(wantUserValue, gotUser[wantUserKey]) {
								t.Errorf("i %v key %v, got %+v, want %+v", i, wantUserKey, gotUser[wantUserKey], wantUserValue)
							}
						}
					}

					continue
				}

				if !reflect.DeepEqual(wantValue, gotResp[wantKey]) {
					t.Errorf("key %v, got %+v, want %+v", wantKey, gotResp[wantKey], wantValue)
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type args struct {
		dbProvider provider.IDBProvider
		user       map[string]interface{}
		urlValues  map[string]string
	}
	tests := []struct {
		name             string
		args             args
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "user updated",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				user: map[string]interface{}{
					"id":       1,
					"username": "test1",
					"role":     1,
					"Password": "testPass11111",
					"token":    "testToken",
					"deleted":  0,
				},
				urlValues: map[string]string{
					"id": "1",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "user updated",
				"user": map[string]interface{}{
					"id":       float64(1),
					"username": "test1",
					"role":     float64(1),
					"Password": "testPass11111",
					"token":    "testToken",
					"deleted":  float64(0),
				},
			},
			wantResponseCode: 200,
		},
		{
			name: "user not updated, this user by name not exist",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				user: map[string]interface{}{
					"id":       1,
					"username": "test11111",
					"role":     1,
					"Password": "testPass11111",
					"token":    "testToken",
					"deleted":  0,
				},
				urlValues: map[string]string{
					"id": "1",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.NotExistResp,
			},
			wantResponseCode: 403,
		},
		{
			name: "user not updated, id is empty",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				user: map[string]interface{}{
					"id":       1,
					"username": "test1",
					"role":     1,
					"Password": "testPass11111",
					"token":    "testToken",
					"deleted":  0,
				},
				urlValues: map[string]string{},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.IdIsEmptyResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "user not updated, consts.IdIsNotNumberResp",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				user: map[string]interface{}{
					"id":       1,
					"username": "test1",
					"role":     1,
					"Password": "testPass11111",
					"token":    "testToken",
					"deleted":  0,
				},
				urlValues: map[string]string{
					"id": "test",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.IdIsNotNumberResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "user not updated, bad input fields",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				user: map[string]interface{}{
					"id": "test",
				},
				urlValues: map[string]string{
					"id": "1",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.BadInputDataResp,
			},
			wantResponseCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bData, err := json.Marshal(tt.args.user)
			if err != nil {
				t.Fatalf("Marshal request form:%v", err)
			}

			buffer := bytes.NewBuffer(bData)

			request, err := http.NewRequest(http.MethodPut, consts.UrlUserV1, buffer)
			if err != nil {
				t.Fatalf("New Request error: %v", err)
			}

			request = mux.SetURLVars(request, tt.args.urlValues)
			recorder := httptest.NewRecorder()

			handler := UpdateUser(tt.args.dbProvider)
			handler.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Body)
			if err != nil {
				t.Errorf("Read response error:%v", err)
			}

			if !reflect.DeepEqual(recorder.Code, tt.wantResponseCode) {
				t.Errorf("response code got %v wantResp %v", recorder.Code, tt.wantResponseCode)
			}

			var gotResp = map[string]interface{}{}

			err = json.Unmarshal(body, &gotResp)
			if err != nil {
				t.Errorf("Unmarshal response body error:%v", err)
			}

			for wantKey, wantValue := range tt.wantResponseBody {
				if !reflect.DeepEqual(wantValue, gotResp[wantKey]) {
					t.Errorf("UpdateUser() key `%v`, got `%+v`, want `%+v`", wantKey, gotResp[wantKey], wantValue)
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type args struct {
		dbProvider provider.IDBProvider
		urlValues  map[string]string
	}
	tests := []struct {
		name             string
		args             args
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "user deleted",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				urlValues: map[string]string{
					"id": "1",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "user deleted",
			},
			wantResponseCode: 200,
		},
		{
			name: "user not deleted, not exist",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				urlValues: map[string]string{
					"id": "111",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.NotDeletedResp,
			},
			wantResponseCode: 500,
		},
		{
			name: "user not deleted, id is empty",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				urlValues: map[string]string{},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.IdIsEmptyResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "user not deleted, consts.IdIsNotNumberResp",
			args: args{
				dbProvider: provider.NewMockUsersDBProvider(
					map[int]models.Model{
						1: &models.User{
							ID:       1,
							Username: "test1",
							Role:     0,
							Password: "testPass1",
							Deleted:  0,
						},
						2: &models.User{
							ID:       2,
							Username: "test2",
							Role:     0,
							Password: "testPass2",
							Deleted:  0,
						},
						3: &models.User{
							ID:       3,
							Username: "test3",
							Role:     0,
							Password: "testPass3",
							Deleted:  0,
						},
						4: &models.User{
							ID:       4,
							Username: "test4",
							Role:     0,
							Password: "testPass4",
							Deleted:  0,
						},
						5: &models.User{
							ID:       5,
							Username: "test5",
							Role:     0,
							Password: "testPass5",
							Deleted:  0,
						},
						6: &models.User{
							ID:       6,
							Username: "test6",
							Role:     0,
							Password: "testPass6",
							Deleted:  0,
						},
						7: &models.User{
							ID:       7,
							Username: "test7",
							Role:     0,
							Password: "testPass7",
							Deleted:  0,
						},
						8: &models.User{
							ID:       8,
							Username: "test8",
							Role:     0,
							Password: "testPass8",
							Deleted:  0,
						},
						9: &models.User{
							ID:       9,
							Username: "test9",
							Role:     0,
							Password: "testPass9",
							Deleted:  0,
						},
						10: &models.User{
							ID:       10,
							Username: "test10",
							Role:     0,
							Password: "testPass10",
							Deleted:  0,
						},
						11: &models.User{
							ID:       11,
							Username: "test11",
							Role:     0,
							Password: "testPass11",
							Deleted:  0,
						},
						12: &models.User{
							ID:       12,
							Username: "test12",
							Role:     0,
							Password: "testPass12",
							Deleted:  0,
						},
					},
				),
				urlValues: map[string]string{
					"id": "test",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.IdIsNotNumberResp,
			},
			wantResponseCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodDelete, consts.UrlUserV1, nil)
			if err != nil {
				t.Fatalf("New Request error: %v", err)
			}

			request = mux.SetURLVars(request, tt.args.urlValues)
			recorder := httptest.NewRecorder()

			handler := DeleteUser(tt.args.dbProvider)
			handler.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Body)
			if err != nil {
				t.Errorf("Read response error:%v", err)
			}

			if !reflect.DeepEqual(recorder.Code, tt.wantResponseCode) {
				t.Errorf("response code got %v wantResp %v", recorder.Code, tt.wantResponseCode)
			}

			var gotResp = map[string]interface{}{}

			err = json.Unmarshal(body, &gotResp)
			if err != nil {
				t.Errorf("Unmarshal response body error:%v", err)
			}

			for wantKey, wantValue := range tt.wantResponseBody {
				if !reflect.DeepEqual(wantValue, gotResp[wantKey]) {
					t.Errorf("UpdateUser() key `%v`, got `%+v`, want `%+v`", wantKey, gotResp[wantKey], wantValue)
				}
			}
		})
	}
}
