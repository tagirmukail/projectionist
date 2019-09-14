package controllers

import (
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

func TestDeleteCfg(t *testing.T) {
	type args struct {
		provider  provider.IDBProvider
		urlValues map[string]string
	}
	tests := []struct {
		name             string
		args             args
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "cfg successful deleted",
			args: args{
				provider: provider.NewMockCfgProvider([]models.Model{
					&models.Configuration{
						ID:   1,
						Name: "test",
						Config: map[string]interface{}{
							"test":  "test",
							"test1": "test1",
						},
						Deleted: 0,
					},
				}),
				urlValues: map[string]string{
					"id": "1",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "config deleted",
			},
			wantResponseCode: 200,
		},
		{
			name: "cfg not deleted, empty id",
			args: args{
				provider: provider.NewMockCfgProvider([]models.Model{
					&models.Configuration{
						ID:   1,
						Name: "test",
						Config: map[string]interface{}{
							"test":  "test",
							"test1": "test1",
						},
						Deleted: 0,
					},
				}),
				urlValues: map[string]string{},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.IdIsEmptyResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "cfg not deleted, id is not number",
			args: args{
				provider: provider.NewMockCfgProvider([]models.Model{
					&models.Configuration{
						ID:   1,
						Name: "test",
						Config: map[string]interface{}{
							"test":  "test",
							"test1": "test1",
						},
						Deleted: 0,
					},
				}),
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
			name: "cfg not deleted, not exist",
			args: args{
				provider: provider.NewMockCfgProvider([]models.Model{}),
				urlValues: map[string]string{
					"id": "1",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.NotExistResp,
			},
			wantResponseCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodDelete, consts.UrlCfgV1, nil)
			if err != nil {
				t.Fatalf("New Request error: %v", err)
			}

			request = mux.SetURLVars(request, tt.args.urlValues)
			recorder := httptest.NewRecorder()

			handler := DeleteCfg(tt.args.provider)
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
					t.Errorf("DeleteCfg() key `%v`, got `%+v`, want `%+v`", wantKey, gotResp[wantKey], wantValue)
				}
			}
		})
	}
}

func TestGetCfg(t *testing.T) {
	type args struct {
		provider  provider.IDBProvider
		urlValues map[string]string
	}
	tests := []struct {
		name             string
		args             args
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "get cfg - successful",
			args: args{
				provider: provider.NewMockCfgProvider([]models.Model{
					&models.Configuration{
						ID:   1,
						Name: "test",
						Config: map[string]interface{}{
							"test":  "test",
							"test1": "test1",
						},
						Deleted: 0,
					},
				}),
				urlValues: map[string]string{
					"id": "1",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "",
				"config": map[string]interface{}{
					"id":   float64(1),
					"name": "test",
					"config": map[string]interface{}{
						"test":  "test",
						"test1": "test1",
					},
					"deleted": float64(0),
				},
			},
			wantResponseCode: 200,
		},
		{
			name: "get cfg - not found",
			args: args{
				provider: provider.NewMockCfgProvider([]models.Model{
					&models.Configuration{
						ID:   1,
						Name: "test",
						Config: map[string]interface{}{
							"test":  "test",
							"test1": "test1",
						},
						Deleted: 0,
					},
				}),
				urlValues: map[string]string{
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
			name: "get cfg - id is empty",
			args: args{
				provider: provider.NewMockCfgProvider([]models.Model{
					&models.Configuration{
						ID:   1,
						Name: "test",
						Config: map[string]interface{}{
							"test":  "test",
							"test1": "test1",
						},
						Deleted: 0,
					},
				}),
				urlValues: map[string]string{},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.IdIsEmptyResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "get cfg - id is not number",
			args: args{
				provider: provider.NewMockCfgProvider([]models.Model{
					&models.Configuration{
						ID:   1,
						Name: "test",
						Config: map[string]interface{}{
							"test":  "test",
							"test1": "test1",
						},
						Deleted: 0,
					},
				}),
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
			request, err := http.NewRequest(http.MethodGet, consts.UrlCfgV1, nil)
			if err != nil {
				t.Fatalf("New Request error: %v", err)
			}

			request = mux.SetURLVars(request, tt.args.urlValues)
			recorder := httptest.NewRecorder()

			handler := GetCfg(tt.args.provider)
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
					t.Errorf("GetCfg() key `%v`, got `%+v`, want `%+v`", wantKey, gotResp[wantKey], wantValue)
				}
			}
		})
	}
}

func TestGetCfgList(t *testing.T) {
	type args struct {
		provider  provider.IDBProvider
		urlValues map[string]string
	}
	tests := []struct {
		name             string
		args             args
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "get configs list - successful, list empty",
			args: args{
				provider: provider.NewMockCfgProvider([]models.Model{
					&models.Configuration{
						ID:   1,
						Name: "test",
						Config: map[string]interface{}{
							"test":  "test",
							"test1": "test1",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   2,
						Name: "test2",
						Config: map[string]interface{}{
							"test2":  "test2",
							"test11": "test11",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   3,
						Name: "test3",
						Config: map[string]interface{}{
							"test3":   "test3",
							"test111": "test111",
						},
						Deleted: 1,
					},
					&models.Configuration{
						ID:   4,
						Name: "test4",
						Config: map[string]interface{}{
							"test4":  "test4",
							"test14": "test14",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   5,
						Name: "test5",
						Config: map[string]interface{}{
							"test5":   "test5",
							"test115": "test115",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   6,
						Name: "test6",
						Config: map[string]interface{}{
							"test6":    "test6",
							"test1116": "test1116",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   7,
						Name: "test7",
						Config: map[string]interface{}{
							"test7":  "test7",
							"test17": "test17",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   8,
						Name: "test8",
						Config: map[string]interface{}{
							"test8":   "test8",
							"test118": "test118",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   9,
						Name: "test9",
						Config: map[string]interface{}{
							"test9":    "test9",
							"test1119": "test1119",
						},
						Deleted: 1,
					},
				}),
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
			name: "get configs list - successful",
			args: args{
				provider: provider.NewMockCfgProvider([]models.Model{
					&models.Configuration{
						ID:   1,
						Name: "test",
						Config: map[string]interface{}{
							"test":  "test",
							"test1": "test1",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   2,
						Name: "test2",
						Config: map[string]interface{}{
							"test2":  "test2",
							"test11": "test11",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   3,
						Name: "test3",
						Config: map[string]interface{}{
							"test3":   "test3",
							"test111": "test111",
						},
						Deleted: 1,
					},
					&models.Configuration{
						ID:   4,
						Name: "test4",
						Config: map[string]interface{}{
							"test4":  "test4",
							"test14": "test14",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   5,
						Name: "test5",
						Config: map[string]interface{}{
							"test5":   "test5",
							"test115": "test115",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   6,
						Name: "test6",
						Config: map[string]interface{}{
							"test6":    "test6",
							"test1116": "test1116",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   7,
						Name: "test7",
						Config: map[string]interface{}{
							"test7":  "test7",
							"test17": "test17",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   8,
						Name: "test8",
						Config: map[string]interface{}{
							"test8":   "test8",
							"test118": "test118",
						},
						Deleted: 0,
					},
					&models.Configuration{
						ID:   9,
						Name: "test9",
						Config: map[string]interface{}{
							"test9":    "test9",
							"test1119": "test1119",
						},
						Deleted: 1,
					},
				}),
				urlValues: map[string]string{
					consts.PAGE_PARAM:  "1",
					consts.COUNT_PARAM: "10",
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "",
				consts.KEY_CONFIGS: []interface{}{
					map[string]interface{}{
						"id":   float64(1),
						"name": "test",
						"config": map[string]interface{}{
							"test":  "test",
							"test1": "test1",
						},
						"deleted": float64(0),
					},
					map[string]interface{}{
						"id":   float64(2),
						"name": "test2",
						"config": map[string]interface{}{
							"test2":  "test2",
							"test11": "test11",
						},
						"deleted": float64(0),
					},
					map[string]interface{}{
						"id":   float64(4),
						"name": "test4",
						"config": map[string]interface{}{
							"test4":  "test4",
							"test14": "test14",
						},
						"deleted": float64(0),
					},
					map[string]interface{}{
						"id":   float64(5),
						"name": "test5",
						"config": map[string]interface{}{
							"test5":   "test5",
							"test115": "test115",
						},
						"deleted": float64(0),
					},
					map[string]interface{}{
						"id":   float64(6),
						"name": "test6",
						"config": map[string]interface{}{
							"test6":    "test6",
							"test1116": "test1116",
						},
						"deleted": float64(0),
					},
					map[string]interface{}{
						"id":   float64(7),
						"name": "test7",
						"config": map[string]interface{}{
							"test7":  "test7",
							"test17": "test17",
						},
						"deleted": float64(0),
					},
					map[string]interface{}{
						"id":   float64(8),
						"name": "test8",
						"config": map[string]interface{}{
							"test8":   "test8",
							"test118": "test118",
						},
						"deleted": float64(0),
					},
				},
			},
			wantResponseCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, consts.UrlCfgV1, nil)
			if err != nil {
				t.Fatalf("New Request error: %v", err)
			}

			query := request.URL.Query()
			for paramKey, paramValue := range tt.args.urlValues {
				query.Add(paramKey, paramValue)
			}

			request.URL.RawQuery = query.Encode()

			recorder := httptest.NewRecorder()

			handler := GetCfgList(tt.args.provider)
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
				wantConfigs, ok := wantValue.([]interface{})
				if ok {
					gotConfigs, ok := gotResp[wantKey].([]interface{})
					if !ok {
						t.Errorf("GetCfgList() key `%v`, got `%v`, want `%v`", wantKey, gotConfigs, wantValue)
					}
					for id, iWantConfig := range wantConfigs {
						wantConfig, ok := iWantConfig.(map[string]interface{})
						if !ok {
							t.Errorf("GetCfgList wantConfig is not map[string]interface{}")
						}

						gotConfig, ok := gotConfigs[id].(map[string]interface{})
						if !ok {
							t.Errorf("GetCfgList gotConfig is not map[string]interface{}")
						}

						for wantCfgKey, wantCfgValue := range wantConfig {
							if !reflect.DeepEqual(wantCfgValue, gotConfig[wantCfgKey]) {
								t.Errorf(
									"GetCfgList() key `%v`, got `%+v`, want `%+v`",
									wantCfgKey,
									gotConfig[wantCfgKey],
									wantCfgValue,
								)
							}
						}

					}
					continue
				}
				if !reflect.DeepEqual(wantValue, gotResp[wantKey]) {
					t.Errorf("GetCfgList() key `%v`, got `%+v`, want `%+v`", wantKey, gotResp[wantKey], wantValue)
				}
			}
		})
	}
}

func TestNewCfg(t *testing.T) {
	type args struct {
		provider  provider.IDBProvider
		urlValues map[string]string
	}
	tests := []struct {
		name             string
		args             args
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, consts.UrlCfgV1, nil)
			if err != nil {
				t.Fatalf("New Request error: %v", err)
			}

			request = mux.SetURLVars(request, tt.args.urlValues)
			recorder := httptest.NewRecorder()

			handler := NewCfg(tt.args.provider)
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
					t.Errorf("NewCfg() key `%v`, got `%+v`, want `%+v`", wantKey, gotResp[wantKey], wantValue)
				}
			}
		})
	}
}

func TestUpdateCfg(t *testing.T) {
	type args struct {
		provider  provider.IDBProvider
		urlValues map[string]string
	}
	tests := []struct {
		name             string
		args             args
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, consts.UrlCfgV1, nil)
			if err != nil {
				t.Fatalf("New Request error: %v", err)
			}

			request = mux.SetURLVars(request, tt.args.urlValues)
			recorder := httptest.NewRecorder()

			handler := UpdateCfg(tt.args.provider)
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
					t.Errorf("UpdateCfg() key `%v`, got `%+v`, want `%+v`", wantKey, gotResp[wantKey], wantValue)
				}
			}
		})
	}
}
