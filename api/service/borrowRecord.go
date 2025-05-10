package service

import (
	"Library-Management-System/api/dto/request"
	"Library-Management-System/api/dto/response"
	"Library-Management-System/api/util/xerror"
)

var BorrowRecord *BorrowRecordService

func init() {
	BorrowRecord = NewBorrowRecordService()
}

type BorrowRecordService struct {
}

func NewBorrowRecordService() *BorrowRecordService {
	return &BorrowRecordService{}
}

func (service *BorrowRecordService) Create(option *request.CreateBorrowRecordOption) (*response.CreateBorrowRecordResponse, xerror.OpenError) {
	panic("implement me")
}

func (service *BorrowRecordService) Describe(option *request.DescribeBorrowRecordOption) (*response.DescribeBorrowRecordResponse, xerror.OpenError) {
	panic("implement me")
}

func (service *BorrowRecordService) Describes(option *request.DescribeBorrowRecordsOption) (*response.DescribeBorrowRecordsResponse, xerror.OpenError) {
	panic("implement me")
}

func (service *BorrowRecordService) Update(option *request.UpdateBorrowRecordOption) (*response.UpdateBorrowRecordResponse, xerror.OpenError) {
	panic("implement me")
}

func (service *BorrowRecordService) Delete(option *request.DeleteBorrowRecordOption) (*response.DeleteBorrowRecordResponse, xerror.OpenError) {
	panic("implement me")
}
