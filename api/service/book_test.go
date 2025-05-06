package service

import (
	"Library-Management-System/api/dto/request"
	"Library-Management-System/api/util"
	"Library-Management-System/api/util/xerror"
	"testing"
)

func TestBookService_Add(t *testing.T) {
	tests := []struct {
		name  string
		req   *request.AddBookOption
		want1 xerror.OpenError
	}{
		{
			name: "add book1",
			req: &request.AddBookOption{
				Name:   "test-book",
				Author: "test-author",
			},
			want1: nil,
		},
	}
	s := &BookService{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := s.Add(tt.req)
			if tt.want1 != nil {
				if got1 == nil || got1.Error() != tt.want1.Error() {
					t.Fatalf("Create() oe = %v, want %v", got1, tt.want1)
				}
			} else {
				if got1 != nil {
					t.Fatalf("Create() got1 = [%d, %s, %s, %s], want %v",
						got1.Code(), got1.Error(), got1.Message(), got1.Details(), tt.want1)
				}
			}
			t.Logf("res: %v\n", util.JSONIgnoreErr(got))
		})
	}
}
