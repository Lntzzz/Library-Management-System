package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Init() *mux.Router {
	r := mux.NewRouter()

	// 路由配置
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
	return r
}
