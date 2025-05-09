package response

import (
	"Library-Management-System/api/util/xerror"
	"encoding/json"
	"net/http"
)

type NullRespData struct {
	RequestId string `json:"requestId"`
}

type ErrorRespData struct {
	Error *ErrorData `json:"error"`
}

type ErrorData struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func CondOp(cond bool, left, right interface{}) interface{} {
	if cond {
		return left
	} else {
		return right
	}
}

func Response(w http.ResponseWriter, ce xerror.OpenError, data interface{}) {
	var rsp interface{}
	if ce != nil {
		w.WriteHeader(ce.Code())
		rsp = ErrorRespData{Error: &ErrorData{Code: ce.Code(), Status: ce.Error(), Message: ce.Message(), Details: CondOp(ce.Details() == nil, []string{}, ce.Details())}}
	} else {
		if data == nil {
			rsp = NullRespData{RequestId: ""}
		} else {
			rsp = data
		}
	}
	jsonRsp, _ := json.Marshal(&rsp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRsp)
}
