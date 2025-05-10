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

func CreateBorrowRecord(w http.ResponseWriter, r *http.Request) {
	userId := r.PostFormValue("userId")
	bookId := r.PostFormValue("bookId")
	option := &request.CreateBorrowRecordOption{
		UserId: userId,
		BookId: bookId,
	}
	ret, oe := service.BorrowRecord.Create(option)
	if oe != nil {
		respUtil.Response(w, oe, nil)
		return
	}
	respUtil.Response(w, nil, response.RespData{Result: ret})
}

func DeleteBorrowRecord(w http.ResponseWriter, r *http.Request) {
	borrowRecordId := mux.Vars(r)["borrowRecordId"] // 从 URL 中获取 borrowRecordId
	option := &request.DeleteBorrowRecordOption{
		Id: borrowRecordId,
	}
	if ret, oe := service.BorrowRecord.Delete(option); oe != nil {
		respUtil.Response(w, oe, nil)
		return
	} else {
		respUtil.Response(w, nil, response.RespData{Result: ret})
		return
	}
}

func UpdateBorrowRecord(w http.ResponseWriter, r *http.Request) {
	borrowRecordId := mux.Vars(r)["borrowRecordId"] // 从 URL 中获取 borrowRecordId

	// 从请求体中解析更新的字段
	userId := r.PostFormValue("userId")
	bookId := r.PostFormValue("bookId")

	// 构造更新选项
	option := &request.UpdateBorrowRecordOption{
		Id:     borrowRecordId,
		UserId: userId,
		BookId: bookId,
	}

	// 调用服务层更新方法
	ret, oe := service.BorrowRecord.Update(option)
	if oe != nil {
		respUtil.Response(w, oe, nil)
		return
	}

	// 返回成功响应
	respUtil.Response(w, nil, response.RespData{Result: ret})
}

func DescribeBorrowRecord(w http.ResponseWriter, r *http.Request) {
	borrowRecordId := mux.Vars(r)["borrowRecordId"] // 从 URL 中获取 borrowRecordId
	option := &request.DescribeBorrowRecordOption{
		Id: borrowRecordId,
	}
	ret, oe := service.BorrowRecord.Describe(option)
	if oe != nil {
		respUtil.Response(w, oe, nil)
		return
	}
	respUtil.Response(w, nil, response.RespData{Result: ret})
}

func DescribeBorrowRecords(w http.ResponseWriter, r *http.Request) {
	option := &request.DescribeBorrowRecordsOption{
		PageNumber: constant.DefaultPageNum,
		PageSize:   constant.DefaultPageSize,
	}
	if ret, oe := service.BorrowRecord.Describes(option); oe != nil {
		respUtil.Response(w, oe, nil)
		return
	} else {
		respUtil.Response(w, nil, response.RespData{Result: ret})
		return
	}
}

func SecKillBooks(w http.ResponseWriter, r *http.Request) {
	userId := r.PostFormValue("userId")
	bookId := r.PostFormValue("bookId")
	option := &request.SecKillBorrowRecordOption{
		BookId: bookId,
		UserId: userId,
	}
	ret, oe := service.Seckill.Seckill(option)
	if oe != nil {
		respUtil.Response(w, oe, nil)
		return
	}
	respUtil.Response(w, nil, response.RespData{Result: ret})
}
