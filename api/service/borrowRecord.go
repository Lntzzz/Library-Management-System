package service

import (
	"Library-Management-System/api/constant"
	"Library-Management-System/api/db"
	"Library-Management-System/api/dto/request"
	"Library-Management-System/api/dto/response"
	"Library-Management-System/api/idgen"
	"Library-Management-System/api/model"
	"Library-Management-System/api/util/xerror"
	"database/sql"
	"errors"
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
	bookById, err := db.Book.Get(option.BookId)
	if err != nil {
		return nil, xerror.ErrBookNotFound
	}
	if bookById.Stock <= 0 {
		return nil, xerror.ErrInvalidBook.SetMessage("book stock not sufficient")
	}
	_, err = db.BorrowRecord.GetByBookIdAndUserId(option.BookId, option.UserId)
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, xerror.ErrInvalidArgument.SetMessage("user already borrowed this book")
	}
	borrowRecord := model.BorrowRecord{
		Id:     idgen.GenBorrowRecordId(),
		UserId: option.UserId,
		BookId: option.BookId,
		Status: constant.Borrowing,
	}
	err = db.BorrowRecord.Add(borrowRecord)
	if err != nil {
		return nil, xerror.ErrInternalServer
	}
	return &response.CreateBorrowRecordResponse{BorrowRecordId: borrowRecord.Id}, nil
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
