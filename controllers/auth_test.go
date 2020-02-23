package controllers

//
//import (
//	"bytes"
//	"encoding/json"
//	"io/ioutil"
//	"net/http"
//	"net/http/httptest"
//	"projectionist/consts"
//	"projectionist/models"
//	"projectionist/provider"
//	"reflect"
//	"testing"
//)
//
//func TestLoginApi(t *testing.T) {
//	type args struct {
//		provider       provider.IDBProvider
//		tokenSecretKey string
//		form           map[string]interface{}
//		urlValues      map[string]string
//	}
//	tests := []struct {
//		name         string
//		args         args
//		wantRespBody map[string]interface{}
//		wantRespCode int
//	}{
//		{
//			name: "login - successful",
//			args: args{
//				provider: provider.NewMockUsersDBProvider(
//					map[int]models.Model{
//						1: &models.User{
//							ID:       1,
//							Username: "test",
//							Deleted:  0,
//							Password: "nogiruki",
//							Role:     0,
//							Token:    "testToken",
//						},
//					},
//				),
//				tokenSecretKey: "",
//				form: map[string]interface{}{
//					"username": "test",
//					"password": "nogiruki",
//				},
//				urlValues: nil,
//			},
//			wantRespBody: map[string]interface{}{
//				"status":  true,
//				"message": "Login successful",
//				"user": map[string]interface{}{
//					"id":       float64(1),
//					"username": "test",
//					"role":     float64(0),
//					"Password": "",
//					"token":    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjF9.Vcp2grZ53t_OG3jwSXsRwfc_UUjboNgZarkAGiX0jgM",
//				},
//			},
//			wantRespCode: 200,
//		},
//		{
//			name: "login - user not found",
//			args: args{
//				provider: provider.NewMockUsersDBProvider(
//					map[int]models.Model{
//						1: &models.User{
//							ID:       1,
//							Username: "test",
//							Deleted:  0,
//							Password: "testPass",
//							Role:     0,
//							Token:    "testToken",
//						},
//					},
//				),
//				tokenSecretKey: "",
//				form: map[string]interface{}{
//					"username": "test1",
//					"password": "testPass",
//				},
//				urlValues: nil,
//			},
//			wantRespBody: map[string]interface{}{
//				"status":  false,
//				"message": consts.NotExistResp,
//			},
//			wantRespCode: 404,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			bData, err := json.Marshal(tt.args.form)
//			if err != nil {
//				t.Fatalf("Marshal request form:%v", err)
//			}
//
//			buffer := bytes.NewBuffer(bData)
//
//			request, err := http.NewRequest(http.MethodGet, consts.UrlCfgV1, buffer)
//			if err != nil {
//				t.Fatalf("New Request error: %v", err)
//			}
//
//			query := request.URL.Query()
//			for paramKey, paramValue := range tt.args.urlValues {
//				query.Add(paramKey, paramValue)
//			}
//
//			request.URL.RawQuery = query.Encode()
//
//			recorder := httptest.NewRecorder()
//
//			handler := LoginApi(tt.args.provider, tt.args.tokenSecretKey)
//			handler.ServeHTTP(recorder, request)
//
//			body, err := ioutil.ReadAll(recorder.Body)
//			if err != nil {
//				t.Errorf("Read response error:%v", err)
//			}
//
//			if !reflect.DeepEqual(recorder.Code, tt.wantRespCode) {
//				t.Errorf("response code got %v wantResp %v", recorder.Code, tt.wantRespCode)
//			}
//
//			var gotResp = map[string]interface{}{}
//
//			err = json.Unmarshal(body, &gotResp)
//			if err != nil {
//				t.Errorf("Unmarshal response body error:%v", err)
//			}
//
//			for wantKey, wantValue := range tt.wantRespBody {
//				wantUser, ok := wantValue.(map[string]interface{})
//				if ok {
//					gotUser, ok := gotResp[wantKey].(map[string]interface{})
//					if !ok {
//						t.Fatalf("LoginApi() gotUser not map[string]itnerface{}")
//					}
//
//					for wantUserKey, wantUserValue := range wantUser {
//						if !reflect.DeepEqual(wantUserValue, gotUser[wantUserKey]) {
//							t.Errorf("LoginApi() key `%v`, got `%+v`, want `%+v`", wantUserKey, gotUser[wantUserKey], wantUserValue)
//						}
//					}
//
//					continue
//				}
//				if !reflect.DeepEqual(wantValue, gotResp[wantKey]) {
//					t.Errorf("LoginApi() key `%v`, got `%+v`, want `%+v`", wantKey, gotResp[wantKey], wantValue)
//				}
//			}
//		})
//	}
//}
