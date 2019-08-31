package controllers

import (
	"bytes"
	"encoding/json"
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
		wantResponse     map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "New User valid.",
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
			wantResponse: map[string]interface{}{
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
			wantResponse: map[string]interface{}{
				"status":  false,
				"message": "Invalid username",
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
			wantResponse: map[string]interface{}{
				"status":  false,
				"message": "Invalid password",
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
			wantResponse: map[string]interface{}{
				"status":  false,
				"message": "Invalid username",
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
			wantResponse: map[string]interface{}{
				"status":  false,
				"message": "Invalid password",
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
			wantResponse: map[string]interface{}{
				"status":  false,
				"message": "Invalid role",
			},
			wantResponseCode: 400,
		},
		{
			name: "New User not created: this user is exist",
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
			wantResponse: map[string]interface{}{
				"status":  false,
				"message": "User exist",
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
				t.Errorf("response code got %v want %v", recorder.Code, tt.wantResponseCode)
			}

			var gotResp = map[string]interface{}{}

			err = json.Unmarshal(body, &gotResp)
			if err != nil {
				t.Errorf("Unmarshal response body error:%v", err)
			}

			for gotKey, gotValue := range gotResp {
				if !reflect.DeepEqual(gotValue, tt.wantResponse[gotKey]) {
					t.Errorf("NewUser() key %v, got %+v, wantResponse %+v", gotKey, gotValue, tt.wantResponse[gotKey])
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		dbProvider provider.IDBProvider
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUser(tt.args.dbProvider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() = %v, wantResponse %v", got, tt.want)
			}
		})
	}
}

func TestGetUserList(t *testing.T) {
	type args struct {
		dbProvider provider.IDBProvider
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUserList(tt.args.dbProvider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserList() = %v, wantResponse %v", got, tt.want)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type args struct {
		dbProvider provider.IDBProvider
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpdateUser(tt.args.dbProvider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateUser() = %v, wantResponse %v", got, tt.want)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type args struct {
		dbProvider provider.IDBProvider
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteUser(tt.args.dbProvider); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteUser() = %v, wantResponse %v", got, tt.want)
			}
		})
	}
}
