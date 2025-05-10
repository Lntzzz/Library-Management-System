package request

type AddBookOption struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

type DeleteBookOption struct {
	Id string `json:"id"`
}

type UpdateBookOption struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

type DescribeBooksOption struct {
	PageNumber int64 `json:"pageNumber"`
	PageSize   int64 `json:"pageSize"`
}

type DescribeBookOption struct {
	Id string `json:"id"`
}

type CreateBorrowRecordOption struct {
	BookId string `json:"bookId"`
	UserId string `json:"userId"`
}

type DeleteBorrowRecordOption struct {
	Id string `json:"id"`
}

type UpdateBorrowRecordOption struct {
	Id     string `json:"id"`
	BookId string `json:"bookId"`
	UserId string `json:"userId"`
}

type DescribeBorrowRecordsOption struct {
	PageNumber int64 `json:"pageNumber"`
	PageSize   int64 `json:"pageSize"`
}

type DescribeBorrowRecordOption struct {
	Id string `json:"id"`
}

type SecKillBorrowRecordOption struct {
	BookId string `json:"bookId"`
	UserId string `json:"userId"`
}
