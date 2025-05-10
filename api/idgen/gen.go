package idgen

import (
	"crypto/rand"
	"math/big"
	mathrandom "math/rand"
	"time"
)

const (
	bookPrefix         = "book-"
	UserIdPrefix       = "user-"
	BorrowRecordPrefix = "borrow-"
)

var character = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
var chLen = len(character)

func GenBookId() string {
	return bookPrefix + Uuid()
}

func GenBorrowRecordId() string {
	return BorrowRecordPrefix + Uuid()
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
