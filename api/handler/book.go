package handler

import (
	"Library-Management-System/api/constant"
	"Library-Management-System/api/dto/request"
	"Library-Management-System/api/dto/response"
	"Library-Management-System/api/service"
	respUtil "Library-Management-System/api/util/response"
	"github.com/gorilla/mux"
	"net/http"
)

func AddBook(w http.ResponseWriter, r *http.Request) {
	bookName := r.PostFormValue("bookName")
	bookAuthor := r.PostFormValue("bookAuthor")
	option := &request.AddBookOption{
		Name:   bookName,
		Author: bookAuthor,
	}
	ret, oe := service.Book.Add(option)
	if oe != nil {
		respUtil.Response(w, oe, nil)
		return
	}
	respUtil.Response(w, nil, response.RespData{Result: ret})
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["bookId"]
	option := &request.DeleteBookOption{
		Id: bookId,
	}
	if ret, oe := service.Book.Delete(option); oe != nil {
		respUtil.Response(w, oe, nil)
		return
	} else {
		respUtil.Response(w, nil, response.RespData{Result: ret})
		return
	}
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["bookId"] // 从 URL 中获取 bookId

	// 从请求体中解析更新的字段
	bookName := r.FormValue("bookName")
	bookAuthor := r.FormValue("bookAuthor")

	// 构造更新选项
	option := &request.UpdateBookOption{
		Id:     bookId,
		Name:   bookName,
		Author: bookAuthor,
	}

	// 调用服务层更新方法
	ret, oe := service.Book.Update(option)
	if oe != nil {
		respUtil.Response(w, oe, nil)
		return
	}

	// 返回成功响应
	respUtil.Response(w, nil, response.RespData{Result: ret})
}

func DescribeBook(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["bookId"] // 从 URL 中获取 bookId
	option := &request.DescribeBookOption{
		Id: bookId,
	}
	ret, oe := service.Book.Describe(option)
	if oe != nil {
		respUtil.Response(w, oe, nil)
		return
	}
	respUtil.Response(w, nil, response.RespData{Result: ret})
}

func DescribeBooks(w http.ResponseWriter, r *http.Request) {
	option := &request.DescribeBooksOption{
		PageNumber: constant.DefaultPageNum,
		PageSize:   constant.DefaultPageSize,
	}
	if ret, oe := service.Book.Describes(option); oe != nil {
		respUtil.Response(w, oe, nil)
		return
	} else {
		respUtil.Response(w, nil, response.RespData{Result: ret})
		return
	}
}
