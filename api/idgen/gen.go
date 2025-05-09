package idgen

import (
	"crypto/rand"
	"math/big"
	mathrandom "math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	bookPrefix   = "book-"
	UserIdPrefix = "user-"
)

var character = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
var chLen = len(character)

func NewBookId() string {
	return bookPrefix + Uuid()
}

func NewUserId() string {
	return UserIdPrefix + Uuid()
}

func Uuid() string {
	var uuidLen = 10
	buf := make([]byte, uuidLen, uuidLen)
	max := big.NewInt(int64(chLen))
	for i := 0; i < uuidLen; i++ {
		random, err := rand.Int(rand.Reader, max)
		if err != nil {
			mathrandom.Seed(time.Now().UnixNano())
			buf[i] = character[mathrandom.Intn(chLen)]
			continue
		}
		buf[i] = character[random.Int64()]
	}
	return string(buf)
}

const (
	recvMsgPrefix     = "recvmsg-"
	recvBillMsgPrefix = "recvbillmsg-"
	sendMsgPrefix     = "sendmsg-"
)

func NewTraceId() string {
	return Uuid()
}

func NewRecvMsgTraceId() string {
	return recvMsgPrefix + Uuid()
}

func NewRecvBillMsgTraceId() string {
	return recvBillMsgPrefix + Uuid()
}

func NewSendMsgTraceId() string {
	return sendMsgPrefix + Uuid()
}

func GenBindId(id string, protocol string, port int32) string {
	return id + "_" + strings.ToLower(protocol) + ":" + strconv.Itoa(int(port))
}
