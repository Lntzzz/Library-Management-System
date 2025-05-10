package xerror

import (
	"strings"
)

/*
	eg: Vpc.Invalid.name

Error Object like:

	"error": {
	  "code": 429,
	  "message": "Out of resource quota",
	  "status": "QUOTA_EXCEEDED",
	  "details": [
	    {
	      "@type": "error.LocalizedMessage",
	      "locale": "zh-CN",
	      "message": "配额已耗尽"
	    }
	  ]
	}
*/
type OpenErrObj struct {
	C int         // err code
	R string      // err resource
	T string      // err type
	P string      // params
	M string      // message
	D interface{} // detail
}

type OpenError interface {
	Code() int
	Message() string
	Error() string // as Status
	Details() interface{}
}

func (e *OpenErrObj) Code() int {
	return e.C
}

func (e *OpenErrObj) Error() string {
	return strings.Replace(e.Status(), ".", "", -1)
}

func (e *OpenErrObj) Status() string {
	s := e.T
	if e.P != ErrNull {
		s += "." + e.P
	}
	return s
}

func (e *OpenErrObj) Message() string {
	return e.M
}

func (e *OpenErrObj) Details() interface{} {
	return e.D
}

func NE(r, t, p, m string, c int, d interface{}) OpenError {
	return &OpenErrObj{c, r, t, p, m, d}
}

const (
	HTTP_200 = 200
	HTTP_202 = 202

	HTTP_400 = 400
	HTTP_401 = 401
	HTTP_403 = 403
	HTTP_404 = 404
	HTTP_409 = 409
	HTTP_429 = 429
	HTTP_499 = 499

	HTTP_500 = 500
	HTTP_501 = 501
	HTTP_503 = 503
	HTTP_504 = 504
)

const (
	ErrNull      = ""
	ErrFilter    = "Filter"
	ErrOrder     = "Order"
	ErrElasticIp = "elasticip"
	ErrRegion    = "Region"
	ErrPin       = "Pin"
	ErrService   = "service"
)

const (
	ErrInvalidArguments   = "INVALID_ARGUMENT"
	ErrFailedPrecondition = "FAILED_PRECONDITION"
	ErrOutOfRange         = "OUT_OF_RANGE"
	ErrUnAuthenticated    = "UNAUTHENTICATED"
	ErrPermissionDenied   = "PERMISSION_DENIED"
	ErrNotFound           = "NOT_FOUND"
	ErrAborted            = "ABORTED"
	ErrAlreadyExists      = "ALREADY_EXISTS"
	ErrQuotaExceededs     = "QUOTA_EXCEEDED"
	ErrRateLimit          = "RATE_LIMIT"
	ErrCancelled          = "CANCELLED"
	ErrUnknown            = "UNKNOWN"
	ErrInternal           = "INTERNAL"
	ErrNotImplemented     = "NOT_IMPLEMENTED"
	ErrUnAvailable        = "UNAVAILABLE"
	ErrDeadlineExceeded   = "DEADLINE_EXCEEDED"
	ErrConflict           = "CONFLICT"
)

const (
	ServerErrError       = "error"
	ServerErrMantaining  = "mantaining"
	ServerErrInvalid     = "invalid"
	ServerErrMiss        = "miss"
	ServerErrMalformed   = "malformed"
	ServerErrNotFound    = "notfound"
	ServerErrInuse       = "inuse"
	ServerErrConflict    = "conflict"
	ServerErrUnsupported = "unsupported"
	ServerErrForbidden   = "forbidden"
	ServerErrLimit       = "limit"
	ServerErrExceeded    = "exceeded"
	ServerErrDuplicate   = "duplicate"
	ServerErrNorse       = "norse"
	ServerErrExisted     = "existed"
	ServerErrQuota       = "quota"
)

var ErrMsgToCode = map[string]int{
	ErrInvalidArguments:   400,
	ErrFailedPrecondition: 400,
	ErrOutOfRange:         400,
	ErrUnAuthenticated:    401,
	ErrPermissionDenied:   403,
	ErrNotFound:           404,
	ErrAborted:            409,
	ErrConflict:           409,
	ErrAlreadyExists:      409,
	ErrQuotaExceededs:     429,
	ErrRateLimit:          429,
	ErrCancelled:          499,
	ErrUnknown:            500,
	ErrInternal:           500,
	ErrNotImplemented:     501,
	ErrUnAvailable:        503,
	ErrDeadlineExceeded:   504,
}
var Err2OpenErrMap = map[string]string{
	ServerErrNotFound:  ErrNotFound,
	ServerErrInuse:     ErrConflict,
	ServerErrConflict:  ErrConflict,
	ServerErrMiss:      ErrInvalidArguments,
	ServerErrMalformed: ErrInvalidArguments,
	ServerErrInvalid:   ErrInvalidArguments,
	ServerErrError:     ErrInternal,
	ServerErrQuota:     ErrQuotaExceededs,
	ServerErrExisted:   ErrConflict,
	ServerErrForbidden: ErrPermissionDenied,
	ServerErrLimit:     ErrInvalidArguments,
}

func Err2OpenErr(msg string, details interface{}) OpenError {
	_, errType, _ := parseErrMsg(msg)
	if oErrType, ok := Err2OpenErrMap[errType]; ok {
		return NE(ErrNull, oErrType, ErrNull, msg, ErrMsgToCode[oErrType], details)
	}
	return NE(ErrNull, ErrInternal, ErrNull, msg, ErrMsgToCode[ErrInternal], details)
}

func SdkErr2OpenErr(sdkErrMsg, sdkErrDetail string, openErrDetails interface{}) OpenError {
	msg := sdkErrMsg
	if sdkErrDetail != "" {
		msg += ":" + sdkErrDetail
	}
	_, errType, _ := parseErrMsg(sdkErrMsg)
	if oErrType, ok := Err2OpenErrMap[errType]; ok {
		return NE(ErrNull, oErrType, ErrNull, msg, ErrMsgToCode[oErrType], openErrDetails)
	}
	return NE(ErrNull, ErrInternal, ErrNull, msg, ErrMsgToCode[ErrInternal], openErrDetails)
}

func parseErrMsg(msg string) (res string, errType string, resId string) {
	items := strings.FieldsFunc(msg, Split)
	if len(items) >= 2 {
		return items[0], items[1], items[len(items)-1]
	}
	return
}

func Split(r rune) bool {
	return r == ':' || r == '.'
}
