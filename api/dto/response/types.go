package response

import "Library-Management-System/api/model"

type RespData struct {
	Result interface{} `json:"result"`
}

type AddBookResponse struct {
	BookId string `json:"bookId"`
}

type DeleteBookResponse struct {
	BookId string `json:"bookId"`
}

type UpdateBookResponse struct {
	BookId string `json:"bookId"`
}

type DescribeBookResponse struct {
	Book model.Book `json:"book"`
}

type DescribesBooksResponse struct {
	Books []model.Book `json:"books"`
}

type CreateBorrowRecordResponse struct {
	BorrowRecordId string `json:"borrowRecordId"`
}

type DescribeBorrowRecordResponse struct {
	BorrowRecord model.BorrowRecord `json:"borrowRecord"`
}

type DescribeBorrowRecordsResponse struct {
	BorrowRecords []model.BorrowRecord `json:"borrowRecords"`
}

type UpdateBorrowRecordResponse struct {
	BorrowRecordId string `json:"borrowRecordId"`
}

type DeleteBorrowRecordResponse struct {
	BorrowRecordId string `json:"borrowRecordId"`
}

type SecKillBorrowRecordResponse struct {
	BorrowRecordId string `json:"borrowRecordId"`
	SecKillResult  string `json:"secKillResult"`
}
