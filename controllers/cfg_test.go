package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
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
	helper := NewHelper(t)
	defer helper.ctrl.Finish()

	type args struct {
		queryArgs map[string]string
	}

	tests := []struct {
		name     string
		args     args
		mocks    []func(mockProvider *provider.MockIDBProviderMockRecorder)
		wantCode int
		wantResp map[string]interface{}
	}{
		{
			name: "ok",
			args: args{queryArgs: map[string]string{"id": "1"}},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Delete(&models.Configuration{}, 1).Return(nil)
				},
			},
			wantCode: http.StatusOK,
			wantResp: map[string]interface{}{
				"status":  true,
				"message": "config deleted",
			},
		},
		{
			name:     "id is empty",
			args:     args{queryArgs: map[string]string{}},
			mocks:    []func(mockProvider *provider.MockIDBProviderMockRecorder){},
			wantCode: http.StatusBadRequest,
			wantResp: map[string]interface{}{
				"status":  false,
				"message": consts.IdIsEmptyResp,
			},
		},
		{
			name:     "id not number",
			args:     args{queryArgs: map[string]string{"id": "test"}},
			mocks:    []func(mockProvider *provider.MockIDBProviderMockRecorder){},
			wantCode: http.StatusBadRequest,
			wantResp: map[string]interface{}{
				"status":  false,
				"message": consts.IdIsNotNumberResp,
			},
		},
		{
			name: "cfg not found",
			args: args{queryArgs: map[string]string{"id": "1"}},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Delete(&models.Configuration{}, 1).Return(consts.ErrNotFound)
				},
			},
			wantCode: http.StatusNotFound,
			wantResp: map[string]interface{}{
				"status":  false,
				"message": consts.NotExistResp,
			},
		},
		{
			name: "delete error",
			args: args{queryArgs: map[string]string{"id": "1"}},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Delete(&models.Configuration{}, 1).Return(errors.New("delete error"))
				},
			},
			wantCode: http.StatusInternalServerError,
			wantResp: map[string]interface{}{
				"status":  false,
				"message": consts.NotDeletedResp,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, mock := range tt.mocks {
				mock(helper.mockProvider)
			}

			request, err := http.NewRequest(http.MethodDelete, consts.UrlCfgV1, nil)
			if err != nil {
				t.Fatalf("New Request error: %v", err)
			}

			request = mux.SetURLVars(request, tt.args.queryArgs)
			recorder := httptest.NewRecorder()

			handler := DeleteCfg(helper.provider)
			handler.ServeHTTP(recorder, request)

			body, err := ioutil.ReadAll(recorder.Body)
			if err != nil {
				t.Errorf("Read response error:%v", err)
			}

			if recorder.Code != tt.wantCode {
				t.Errorf("response code got %v wantResp %v", recorder.Code, tt.wantCode)
			}

			var gotResp = map[string]interface{}{}

			err = json.Unmarshal(body, &gotResp)
			if err != nil {
				t.Errorf("Unmarshal response body error:%v", err)
			}

			for wantKey, wantValue := range tt.wantResp {
				if !reflect.DeepEqual(wantValue, gotResp[wantKey]) {
					t.Errorf("DeleteCfg() key `%v`, got `%+v`, want `%+v`", wantKey, gotResp[wantKey], wantValue)
				}
			}
		})
	}
}

func TestGetCfg(t *testing.T) {
	helper := NewHelper(t)
	defer helper.ctrl.Finish()

	type args struct {
		provider  provider.IDBProvider
		urlValues map[string]string
	}
	tests := []struct {
		name             string
		args             args
		mocks            []func(mockProvider *provider.MockIDBProviderMockRecorder)
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "get cfg - successful",
			args: args{
				provider: helper.provider,
				urlValues: map[string]string{
					"id": "1",
				},
			},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.GetByID(models.Model(nil), int64(1)).
						Return(&models.Configuration{
							ID:   1,
							Name: "test",
							Config: map[string]interface{}{
								"test":  "test",
								"test1": "test1",
							},
							Deleted: 0,
						}, nil)
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
			wantResponseCode: http.StatusOK,
		},
		{
			name: "get cfg - not found",
			args: args{
				provider: helper.provider,
				urlValues: map[string]string{
					"id": "2",
				},
			},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.GetByID(models.Model(nil), int64(2)).
						Return(nil, consts.ErrNotFound)
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.NotExistResp,
			},
			wantResponseCode: http.StatusNotFound,
		},
		{
			name: "get cfg - error",
			args: args{
				provider: helper.provider,
				urlValues: map[string]string{
					"id": "2",
				},
			},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.GetByID(models.Model(nil), int64(2)).
						Return(nil, errors.New("get error"))
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.SmtWhenWrongResp,
			},
			wantResponseCode: http.StatusInternalServerError,
		},
		{
			name: "get cfg - id is empty",
			args: args{
				provider:  helper.provider,
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
				provider: helper.provider,
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
			for _, mock := range tt.mocks {
				mock(helper.mockProvider)
			}

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
	helper := NewHelper(t)
	defer helper.ctrl.Finish()

	type args struct {
		provider  provider.IDBProvider
		urlValues map[string]string
	}
	tests := []struct {
		name             string
		args             args
		mocks            []func(mockProvider *provider.MockIDBProviderMockRecorder)
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "get configs list - successful, list empty",
			args: args{
				provider: helper.provider,
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
			name: "get configs list - successful for first page",
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Count(&models.Configuration{}).Return(10, nil)
				},
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Pagination(&models.Configuration{}, 0, 10).Return(
						[]models.Model{
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
								Deleted: 0,
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
								Deleted: 0,
							},
						}, nil)
				},
			},
			args: args{
				provider: helper.provider,
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
						"id":   float64(3),
						"name": "test3",
						"config": map[string]interface{}{
							"test3":   "test3",
							"test111": "test111",
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
					map[string]interface{}{
						"id":   float64(9),
						"name": "test9",
						"config": map[string]interface{}{
							"test9":    "test9",
							"test1119": "test1119",
						},
						"deleted": float64(0),
					},
				},
			},
			wantResponseCode: 200,
		},
		{
			name: "get configs list - successful for second page",
			args: args{
				provider: helper.provider,
				urlValues: map[string]string{
					consts.PAGE_PARAM:  "2",
					consts.COUNT_PARAM: "3",
				},
			},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Count(&models.Configuration{}).Return(10, nil)
				},
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Pagination(&models.Configuration{}, 3, 6).Return(
						[]models.Model{
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
						}, nil)
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "",
				consts.KEY_CONFIGS: []interface{}{
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
				},
			},
			wantResponseCode: 200,
		},
		{
			name: "get configs list - error get count",
			args: args{
				provider: helper.provider,
				urlValues: map[string]string{
					consts.PAGE_PARAM:  "2",
					consts.COUNT_PARAM: "3",
				},
			},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Count(&models.Configuration{}).Return(0, errors.New("count error"))
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.SmtWhenWrongResp,
			},
			wantResponseCode: http.StatusInternalServerError,
		},
		{
			name: "get configs list - error pagination",
			args: args{
				provider: helper.provider,
				urlValues: map[string]string{
					consts.PAGE_PARAM:  "2",
					consts.COUNT_PARAM: "3",
				},
			},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Count(&models.Configuration{}).Return(10, nil)
				},
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Pagination(&models.Configuration{}, 3, 6).Return(
						nil, errors.New("error"))
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.SmtWhenWrongResp,
			},
			wantResponseCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		for _, mock := range tt.mocks {
			mock(helper.mockProvider)
		}

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

						if len(wantConfigs) != len(gotConfigs) {
							t.Fatalf("GetCfgList len configs got %d want %d", len(gotConfigs), len(wantConfigs))
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
	helper := NewHelper(t)
	defer helper.ctrl.Finish()

	type args struct {
		urlValues map[string]string
		provider  provider.IDBProvider
		config    map[string]interface{}
	}
	tests := []struct {
		name             string
		args             args
		mocks            []func(mockProvider *provider.MockIDBProviderMockRecorder)
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "new config successful created",
			args: args{
				provider: helper.provider,
				config: map[string]interface{}{
					"name":   "test",
					"config": map[string]string{"test333": "test333"},
				},
				urlValues: map[string]string{},
			},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Save(&models.Configuration{
						ID:      0,
						Name:    "test",
						Config:  map[string]interface{}{"test333": "test333"},
						Deleted: 0,
					}).Return(nil)
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "File test saved",
			},
			wantResponseCode: 200,
		},
		{
			name: "new config create error",
			args: args{
				provider: helper.provider,
				config: map[string]interface{}{
					"name":   "test",
					"config": map[string]string{"test333": "test333"},
				},
				urlValues: map[string]string{},
			},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Save(&models.Configuration{
						ID:      0,
						Name:    "test",
						Config:  map[string]interface{}{"test333": "test333"},
						Deleted: 0,
					}).Return(errors.New("save error"))
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.SmtWhenWrongResp,
			},
			wantResponseCode: http.StatusInternalServerError,
		},
		{
			name: "new config - parameter name empty",
			args: args{
				provider: helper.provider,
				config: map[string]interface{}{
					"config": map[string]string{"test333": "test333"},
				},
				urlValues: map[string]string{},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.InputDataInvalidResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "new config - bad input data",
			args: args{
				provider: helper.provider,
				config: map[string]interface{}{
					"name":   "test",
					"config": []string{},
				},
				urlValues: map[string]string{},
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
			for _, mock := range tt.mocks {
				mock(helper.mockProvider)
			}

			bData, err := json.Marshal(tt.args.config)
			if err != nil {
				t.Fatalf("Marshal request form:%v", err)
			}

			buffer := bytes.NewBuffer(bData)

			request, err := http.NewRequest(http.MethodGet, consts.UrlCfgV1, buffer)
			if err != nil {
				t.Fatalf("New Request error: %v", err)
			}

			query := request.URL.Query()
			for paramKey, paramValue := range tt.args.urlValues {
				query.Add(paramKey, paramValue)
			}

			request.URL.RawQuery = query.Encode()

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
	helper := NewHelper(t)
	defer helper.ctrl.Finish()

	type args struct {
		config    map[string]interface{}
		provider  provider.IDBProvider
		urlValues map[string]string
	}
	tests := []struct {
		name             string
		args             args
		mocks            []func(mockProvider *provider.MockIDBProviderMockRecorder)
		wantResponseBody map[string]interface{}
		wantResponseCode int
	}{
		{
			name: "update config - successful",
			args: args{
				config: map[string]interface{}{
					"id":   1,
					"name": "test1-1",
					"config": map[string]interface{}{
						"test":  "test22222",
						"test1": "test232323",
						"test3": "test333323",
					},
					"deleted": 0,
				},
				provider: helper.provider,
				urlValues: map[string]string{
					"id": "1",
				},
			},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Update(&models.Configuration{
						ID:   1,
						Name: "test1-1",
						Config: map[string]interface{}{
							"test":  "test22222",
							"test1": "test232323",
							"test3": "test333323",
						},
						Deleted: 0,
					}, 1).Return(nil)
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  true,
				"message": "Config updated",
				"config": map[string]interface{}{
					"id":   float64(1),
					"name": "test1-1",
					"config": map[string]interface{}{
						"test":  "test22222",
						"test1": "test232323",
						"test3": "test333323",
					},
					"deleted": float64(0),
				},
			},
			wantResponseCode: 200,
		},
		{
			name: "update config - error",
			args: args{
				config: map[string]interface{}{
					"id":   1,
					"name": "test1-1",
					"config": map[string]interface{}{
						"test":  "test22222",
						"test1": "test232323",
						"test3": "test333323",
					},
					"deleted": 0,
				},
				provider: helper.provider,
				urlValues: map[string]string{
					"id": "1",
				},
			},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Update(&models.Configuration{
						ID:   1,
						Name: "test1-1",
						Config: map[string]interface{}{
							"test":  "test22222",
							"test1": "test232323",
							"test3": "test333323",
						},
						Deleted: 0,
					}, 1).Return(errors.New("update error"))
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.SmtWhenWrongResp,
			},
			wantResponseCode: http.StatusInternalServerError,
		},
		{
			name: "update config - not exist",
			args: args{
				config: map[string]interface{}{
					"id":   3,
					"name": "test1-1",
					"config": map[string]interface{}{
						"test":  "test22222",
						"test1": "test232323",
						"test3": "test333323",
					},
					"deleted": 0,
				},
				provider: helper.provider,
				urlValues: map[string]string{
					"id": "3",
				},
			},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){
				func(mockProvider *provider.MockIDBProviderMockRecorder) {
					mockProvider.Update(&models.Configuration{
						ID:   3,
						Name: "test1-1",
						Config: map[string]interface{}{
							"test":  "test22222",
							"test1": "test232323",
							"test3": "test333323",
						},
						Deleted: 0,
					}, 3).Return(consts.ErrNotFound)
				},
			},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.NotExistResp,
			},
			wantResponseCode: 404,
		},
		{
			name: "update config - id is empty",
			args: args{
				config: map[string]interface{}{
					"id":   3,
					"name": "test1-1",
					"config": map[string]interface{}{
						"test":  "test22222",
						"test1": "test232323",
						"test3": "test333323",
					},
					"deleted": 0,
				},
				provider:  helper.provider,
				urlValues: map[string]string{},
			},
			mocks: []func(mockProvider *provider.MockIDBProviderMockRecorder){},
			wantResponseBody: map[string]interface{}{
				"status":  false,
				"message": consts.IdIsEmptyResp,
			},
			wantResponseCode: 400,
		},
		{
			name: "update config - id is not number",
			args: args{
				config: map[string]interface{}{
					"id":   3,
					"name": "test1-1",
					"config": map[string]interface{}{
						"test":  "test22222",
						"test1": "test232323",
						"test3": "test333323",
					},
					"deleted": 0,
				},
				provider: helper.provider,
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
			name: "update config - bad input data",
			args: args{
				config:   map[string]interface{}{},
				provider: helper.provider,
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
			for _, mock := range tt.mocks {
				mock(helper.mockProvider)
			}

			bData, err := json.Marshal(tt.args.config)
			if err != nil {
				t.Fatalf("Marshal request form:%v", err)
			}

			buffer := bytes.NewBuffer(bData)

			request, err := http.NewRequest(http.MethodGet, consts.UrlCfgV1, buffer)
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
