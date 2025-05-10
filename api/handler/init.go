package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Init() *mux.Router {
	r := mux.NewRouter()

	// Book路由配置
	r.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			AddBook(w, r)
		case http.MethodGet:
			DescribeBooks(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodPost, http.MethodGet)
	r.HandleFunc("/book/describer/{bookId}", DescribeBook)
	r.HandleFunc("/book/updater/{bookId}", UpdateBook)
	r.HandleFunc("/book/deleter/{bookId}", DeleteBook)

	// BorrowRecord路由配置
	r.HandleFunc("/borrowRecords", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateBorrowRecord(w, r)
		case http.MethodGet:
			DescribeBorrowRecords(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	r.HandleFunc("/borrowRecords/describer/{borrowRecordId}", DescribeBorrowRecord)
	r.HandleFunc("/borrowRecords/updater/{borrowRecordId}", UpdateBorrowRecord)
	r.HandleFunc("/borrowRecords/deleter/{borrowRecordId}", DeleteBorrowRecord)

	// User路由配置

	// SecKill路由配置
	r.HandleFunc("/borrowRecords/secKill}", SecKillBooks)

	return r
}
