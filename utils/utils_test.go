package utils

import (
	"os"
	"testing"
	"time"
)

func TestCreateDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create dir",
			args: args{
				path: "./test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateDir(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("CreateDir() error = %v, wantErr %v", err, tt.wantErr)
			}

			time.Sleep(10 * time.Second)

			err := os.Remove(tt.args.path)
			if err != nil {
				t.Errorf("CreateDir dir: %v delete error: %v", &tt.args.path, err)
			}
		})
	}
}

func TestGetFileName(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "full file path",
			args: args{
				path: "/first/path/test.json",
			},
			want: "test.json",
		},
		{
			name: "empty file path",
			args: args{
				path: "",
			},
			want: "",
		},
		{
			name: "file path without file extension",
			args: args{
				path: "/path/json",
			},
			want: "json",
		},
		{
			name: "file path without absolute path",
			args: args{
				path: ".file.json",
			},
			want: ".file.json",
		},
		{
			name: "file path relative path",
			args: args{
				path: "~/path/file.new.json",
			},
			want: "file.new.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileName(tt.args.path); got != tt.want {
				t.Errorf("GetFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination(t *testing.T) {
	type args struct {
		page  int
		count int
	}
	tests := []struct {
		name      string
		args      args
		wantStart int
		wantEnd   int
	}{
		{
			name: "pagination - page 1, count 10",
			args: args{
				page:  1,
				count: 10,
			},
			wantStart: 0,
			wantEnd:   10,
		},
		{
			name: "pagination - page 0, count 10",
			args: args{
				page:  0,
				count: 10,
			},
			wantStart: 0,
			wantEnd:   0,
		},
		{
			name: "pagination - page 2, count 10",
			args: args{
				page:  2,
				count: 10,
			},
			wantStart: 10,
			wantEnd:   20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := Pagination(tt.args.page, tt.args.count)
			if gotStart != tt.wantStart {
				t.Errorf("Pagination() gotStart = %v, want %v", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("Pagination() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}
