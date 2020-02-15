package provider

import (
	"github.com/dgraph-io/badger/v2"
	"projectionist/models"
	"projectionist/utils/errors"
	"reflect"
	"strconv"
	"testing"
)

func TestNewCfgProvider(t *testing.T) {
	db := NewTestDB(t, false, false)
	defer db.Close()

	type args struct {
		db *badger.DB
	}
	tests := []struct {
		name    string
		args    args
		want    *CfgProvider
		wantErr bool
	}{
		{
			name:    "ok",
			args:    args{db: db},
			want:    &CfgProvider{db: db},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCfgProvider(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCfgProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCfgProvider() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMaxID(t *testing.T) {
	type data struct {
		entries []*badger.Entry
	}

	type args struct {
		db *badger.DB
	}
	tests := []struct {
		name    string
		args    args
		data    data
		want    int
		wantErr bool
	}{
		{
			name: "error key not found",
			args: args{
				db: NewTestDB(t, false, false),
			},
			data:    data{[]*badger.Entry{}},
			want:    0,
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				db: NewTestDB(t, false, false),
			},
			data: data{
				entries: []*badger.Entry{
					badger.NewEntry([]byte(MaxID), []byte(strconv.Itoa(5))),
				},
			},
			want:    5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.args.db.Close()

			err := prepareData(tt.args.db, tt.data.entries)
			if err != nil {
				t.Error(err)
			}

			tt.args.db.View(func(txn *badger.Txn) error {
				got, err := getMaxID(txn)
				if (err != nil) != tt.wantErr {
					t.Errorf("getMaxID() error = %v, wantErr %v", err, tt.wantErr)
				}
				if got != tt.want {
					t.Errorf("getMaxID() got = %v, want %v", got, tt.want)
				}

				return nil
			})
		})
	}
}

func TestCfgProvider_Save(t *testing.T) {
	type data struct {
		entries []*badger.Entry
	}

	type fields struct {
		db    *badger.DB
		maxID int
	}
	type args struct {
		m models.Model
	}
	tests := []struct {
		name      string
		fields    fields
		data      data
		args      args
		wantMaxID int
		wantErr   bool
	}{
		{
			name: "ok",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{
				[]*badger.Entry{
					{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
						ID:      1,
						Name:    "test",
						Config:  map[string]interface{}{"host": "localhost"},
						Deleted: 0,
					})},
				},
			},
			args: args{
				m: &models.Configuration{
					ID:      2,
					Name:    "test2",
					Config:  map[string]interface{}{"host": "localhost1"},
					Deleted: 0,
				},
			},
			wantMaxID: 1,
			wantErr:   false,
		},
		{
			name: "cfg is exist",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{
				[]*badger.Entry{
					{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
						ID:      1,
						Name:    "test",
						Config:  map[string]interface{}{"host": "localhost"},
						Deleted: 0,
					})},
				},
			},
			args: args{
				m: &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				},
			},
			wantMaxID: 1,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.fields.db.Close()

			err := prepareData(tt.fields.db, tt.data.entries)
			if err != nil {
				t.Errorf("prepareData() error: %v", err)
			}

			c := &CfgProvider{
				db:    tt.fields.db,
				maxID: tt.fields.maxID,
			}
			if err := c.Save(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			if c.maxID != tt.wantMaxID {
				t.Errorf("Save() maxID = %v, wantMaxID %v", c.maxID, tt.wantMaxID)
			}
		})
	}
}

func TestCfgProvider_GetByName(t *testing.T) {
	type data struct {
		entries []*badger.Entry
	}

	type fields struct {
		db    *badger.DB
		maxID int
	}
	type args struct {
		in0  models.Model
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		data    data
		args    args
		want    models.Model
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{
				[]*badger.Entry{
					{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
						ID:      1,
						Name:    "test",
						Config:  map[string]interface{}{"host": "localhost"},
						Deleted: 0,
					})},
				},
			},
			args: args{
				in0:  &models.Configuration{},
				name: "test",
			},
			want: &models.Configuration{
				ID:      1,
				Name:    "test",
				Config:  map[string]interface{}{"host": "localhost"},
				Deleted: 0,
			},
			wantErr: false,
		},
		{
			name: "bad key",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{
				[]*badger.Entry{
					{Key: []byte("1test|0"), Value: marshalModel(t, &models.Configuration{
						ID:      1,
						Name:    "test",
						Config:  map[string]interface{}{"host": "localhost"},
						Deleted: 0,
					})},
				},
			},
			args: args{
				in0:  &models.Configuration{},
				name: "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ok",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{
				[]*badger.Entry{
					{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
						ID:      1,
						Name:    "test",
						Config:  map[string]interface{}{"host": "localhost"},
						Deleted: 0,
					})},
				},
			},
			args: args{
				in0:  &models.Configuration{},
				name: "test",
			},
			want: &models.Configuration{
				ID:      1,
				Name:    "test",
				Config:  map[string]interface{}{"host": "localhost"},
				Deleted: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.fields.db.Close()

			err := prepareData(tt.fields.db, tt.data.entries)
			if err != nil {
				t.Errorf("prepareData error: %v", err)
			}

			c := &CfgProvider{
				db:    tt.fields.db,
				maxID: tt.fields.maxID,
			}
			got, err := c.GetByName(tt.args.in0, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByName() error = %v, wantErr %v", err.Error(), tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCfgProvider_GetByID(t *testing.T) {
	type data struct {
		entries []*badger.Entry
	}

	type fields struct {
		db    *badger.DB
		maxID int
	}
	type args struct {
		in0 models.Model
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		data    data
		args    args
		want    models.Model
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{entries: []*badger.Entry{
				{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
			}},
			args: args{
				in0: &models.Configuration{},
				id:  1,
			},
			want: &models.Configuration{
				ID:      1,
				Name:    "test",
				Config:  map[string]interface{}{"host": "localhost"},
				Deleted: 0,
			},
			wantErr: false,
		},
		{
			name: "not exist",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{entries: []*badger.Entry{
				{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
			}},
			args: args{
				in0: &models.Configuration{},
				id:  2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := prepareData(tt.fields.db, tt.data.entries)
			if err != nil {
				t.Errorf("prepareData error: %v", err)
			}
			defer tt.fields.db.Close()

			c := &CfgProvider{
				db:    tt.fields.db,
				maxID: tt.fields.maxID,
			}
			got, err := c.GetByID(tt.args.in0, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCfgProvider_IsExistByName(t *testing.T) {
	type data struct {
		entries []*badger.Entry
	}
	type fields struct {
		db    *badger.DB
		maxID int
	}
	type args struct {
		m models.Model
	}
	tests := []struct {
		name   string
		fields fields
		data   data
		args   args
		want   error
		want1  bool
	}{
		{
			name: "ok",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{entries: []*badger.Entry{
				{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
			}},
			args: args{
				m: &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				},
			},
			want:  nil,
			want1: true,
		},
		{
			name: "not exist",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{entries: []*badger.Entry{
				{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
			}},
			args: args{
				m: &models.Configuration{
					ID:      2,
					Name:    "test2",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				},
			},
			want:  errors.ErrNotExist,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := prepareData(tt.fields.db, tt.data.entries)
			if err != nil {
				t.Errorf("prepareData error: %v", err)
			}
			defer tt.fields.db.Close()

			c := &CfgProvider{
				db:    tt.fields.db,
				maxID: tt.fields.maxID,
			}
			got, got1 := c.IsExistByName(tt.args.m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsExistByName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("IsExistByName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCfgProvider_Count(t *testing.T) {
	type data struct {
		entries []*badger.Entry
	}
	type fields struct {
		db    *badger.DB
		maxID int
	}
	type args struct {
		in0 models.Model
	}
	tests := []struct {
		name    string
		fields  fields
		data    data
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{entries: []*badger.Entry{
				{Key: []byte(MaxID), Value: []byte(strconv.Itoa(1))},
				{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
			}},
			args:    args{in0: nil},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := prepareData(tt.fields.db, tt.data.entries)
			if err != nil {
				t.Errorf("prepareData error: %v", err)
			}
			defer tt.fields.db.Close()

			c := &CfgProvider{
				db:    tt.fields.db,
				maxID: tt.fields.maxID,
			}
			got, err := c.Count(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Count() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCfgProvider_Pagination(t *testing.T) {
	type data struct {
		entries []*badger.Entry
	}
	type fields struct {
		db    *badger.DB
		maxID int
	}
	type args struct {
		m     models.Model
		start int
		stop  int
	}
	tests := []struct {
		name    string
		fields  fields
		data    data
		args    args
		want    []models.Model
		wantErr bool
	}{
		{
			name: "normal pagination",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{entries: []*badger.Entry{
				{Key: []byte(MaxID), Value: []byte(strconv.Itoa(1))},
				{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
				{Key: []byte("2|test2|0"), Value: marshalModel(t, &models.Configuration{
					ID:      2,
					Name:    "test2",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
				{Key: []byte("3|test3|0"), Value: marshalModel(t, &models.Configuration{
					ID:      3,
					Name:    "test3",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
				{Key: []byte("4|test4|0"), Value: marshalModel(t, &models.Configuration{
					ID:      4,
					Name:    "test4",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
				{Key: []byte("5|test5|0"), Value: marshalModel(t, &models.Configuration{
					ID:      5,
					Name:    "test5",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
			}},
			args: args{
				m:     nil,
				start: 0,
				stop:  4,
			},
			want: []models.Model{
				&models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				},
				&models.Configuration{
					ID:      2,
					Name:    "test2",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				},
				&models.Configuration{
					ID:      3,
					Name:    "test3",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				},
				&models.Configuration{
					ID:      4,
					Name:    "test4",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				},
			},
			wantErr: false,
		},
		{
			name: "start > stop",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{entries: []*badger.Entry{
				{Key: []byte(MaxID), Value: []byte(strconv.Itoa(1))},
				{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
				{Key: []byte("2|test2|0"), Value: marshalModel(t, &models.Configuration{
					ID:      2,
					Name:    "test2",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
				{Key: []byte("3|test3|0"), Value: marshalModel(t, &models.Configuration{
					ID:      3,
					Name:    "test3",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
				{Key: []byte("4|test4|0"), Value: marshalModel(t, &models.Configuration{
					ID:      4,
					Name:    "test4",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
				{Key: []byte("5|test5|0"), Value: marshalModel(t, &models.Configuration{
					ID:      5,
					Name:    "test5",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
			}},
			args: args{
				m:     nil,
				start: 5,
				stop:  4,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "start < 0",
			fields: fields{
				db:    NewTestDB(t, false, false),
				maxID: 0,
			},
			data: data{entries: []*badger.Entry{
				{Key: []byte(MaxID), Value: []byte(strconv.Itoa(1))},
				{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
				{Key: []byte("2|test2|0"), Value: marshalModel(t, &models.Configuration{
					ID:      2,
					Name:    "test2",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
				{Key: []byte("3|test3|0"), Value: marshalModel(t, &models.Configuration{
					ID:      3,
					Name:    "test3",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
				{Key: []byte("4|test4|0"), Value: marshalModel(t, &models.Configuration{
					ID:      4,
					Name:    "test4",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
				{Key: []byte("5|test5|0"), Value: marshalModel(t, &models.Configuration{
					ID:      5,
					Name:    "test5",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
			}},
			args: args{
				m:     nil,
				start: -10,
				stop:  4,
			},
			want: []models.Model{
				&models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				},
				&models.Configuration{
					ID:      2,
					Name:    "test2",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				},
				&models.Configuration{
					ID:      3,
					Name:    "test3",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				},
				&models.Configuration{
					ID:      4,
					Name:    "test4",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := prepareData(tt.fields.db, tt.data.entries)
			if err != nil {
				t.Errorf("prepareData error: %v", err)
			}
			defer tt.fields.db.Close()

			c := &CfgProvider{
				db:    tt.fields.db,
				maxID: tt.fields.maxID,
			}
			got, err := c.Pagination(tt.args.m, tt.args.start, tt.args.stop)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pagination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pagination() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCfgProvider_Update(t *testing.T) {
	type data struct {
		entries []*badger.Entry
	}
	type fields struct {
		db    *badger.DB
		maxID int
	}
	type args struct {
		m  models.Model
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		data    data
		args    args
		wantErr bool
	}{
		{
			name:   "update ok",
			fields: fields{db: NewTestDB(t, false, false)},
			data: data{entries: []*badger.Entry{
				{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
			}},
			args: args{
				m: &models.Configuration{
					ID:      1,
					Name:    "testtest",
					Config:  map[string]interface{}{"host": "localhost", "port": 2222},
					Deleted: 0,
				},
				id: 1,
			},
			wantErr: false,
		},
		{
			name:   "not exist",
			fields: fields{db: NewTestDB(t, false, false)},
			data: data{entries: []*badger.Entry{
				{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
			}},
			args: args{
				m: &models.Configuration{
					ID:      2,
					Name:    "testtest",
					Config:  map[string]interface{}{"host": "localhost", "port": 2222},
					Deleted: 0,
				},
				id: 2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := prepareData(tt.fields.db, tt.data.entries)
			if err != nil {
				t.Errorf("prepareData error: %v", err)
			}
			defer tt.fields.db.Close()

			c := &CfgProvider{
				db:    tt.fields.db,
				maxID: tt.fields.maxID,
			}
			if err := c.Update(tt.args.m, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCfgProvider_Delete(t *testing.T) {
	type data struct {
		entries []*badger.Entry
	}
	type fields struct {
		db    *badger.DB
		maxID int
	}
	type args struct {
		m  models.Model
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		data    data
		args    args
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{db: NewTestDB(t, false, false)},
			data: data{entries: []*badger.Entry{
				{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
			}},
			args: args{
				m:  nil,
				id: 1,
			},
			wantErr: false,
		},
		{
			name:   "not exist",
			fields: fields{db: NewTestDB(t, false, false)},
			data: data{entries: []*badger.Entry{
				{Key: []byte("1|test|0"), Value: marshalModel(t, &models.Configuration{
					ID:      1,
					Name:    "test",
					Config:  map[string]interface{}{"host": "localhost"},
					Deleted: 0,
				})},
			}},
			args: args{
				m:  nil,
				id: 2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := prepareData(tt.fields.db, tt.data.entries)
			if err != nil {
				t.Errorf("prepareData error: %v", err)
			}
			defer tt.fields.db.Close()

			c := &CfgProvider{
				db:    tt.fields.db,
				maxID: tt.fields.maxID,
			}
			if err := c.Delete(tt.args.m, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
