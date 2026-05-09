package messages

import (
	"Library-Management-System/api/db"
	"Library-Management-System/api/internal/mq"
	"Library-Management-System/api/model"
	"Library-Management-System/api/util"
	"Library-Management-System/api/util/xerror"
	"encoding/json"
	"fmt"
	"time"
)

var _defaultSecKillConsumer *secKillConsumer

func NewSecKillConsumer(kConf *KafkaConf) mq.Starter {
	if _defaultSecKillConsumer == nil {
		_defaultSecKillConsumer = newSecKillConsumer(kConf)
	}
	return _defaultSecKillConsumer
}

type secKillConsumer struct {
	consumer mq.Consumer
	topic    string
}

func newSecKillConsumer(kConf *KafkaConf) *secKillConsumer {
	c := &secKillConsumer{}
	c.topic = kConf.ConsumerTopicSecKill
	c.consumer = mq.NewKafkaConsumer(&mq.KafkaConfig{
		Name:         c.topic,
		Url:          kConf.Brokers,
		GroupId:      kConf.GroupId,
		Topics:       c.topic,
		ConsumerSize: 10,
		FromOffsets:  kConf.FromOffsets,
	})
	return c
}

func (o *secKillConsumer) Start() {
	o.consumer.Listen(o.event)
}

var c_secKill *secKillConsumer

func SecKillConsumer(kConf *KafkaConf) *secKillConsumer {
	if c_secKill == nil {
		c_secKill = newSecKillConsumer(kConf)
	}
	return c_secKill
}

type SeckillMessage struct {
	Id     string `json:"id"`
	UserID string `json:"user_id"`
	BookID string `json:"book_id"`
	Time   string `json:"time"`
}

func (o *secKillConsumer) event(kafkaMessage interface{}) {
	if message, ok := kafkaMessage.(*mq.KafkaMessage); ok {
		var secKillMsg SeckillMessage
		if err := json.Unmarshal(message.Value, &secKillMsg); err != nil {
			return
		}
		if secKillMsg.Id != "" {
			o.sendBorrowRecord2DB(secKillMsg.Id, secKillMsg.UserID, secKillMsg.BookID, secKillMsg.Time)
		}
	} else {
		fmt.Println("Convert kafkaMessage failed. %s", util.JSONIgnoreErr(kafkaMessage))
	}
}

func (o *secKillConsumer) sendBorrowRecord2DB(id, userId, bookId, secKillTime string) xerror.OpenError {
	parsedTime, err := time.Parse(time.RFC3339, secKillTime)
	if err != nil {
		return xerror.ErrInternalServer.SetMessage("failed to parse time: " + err.Error())
	}
	bookRecord := model.BorrowRecord{
		Id:         id,
		UserId:     userId,
		BookId:     bookId,
		BorrowedAt: parsedTime,
	}
	err = db.BorrowRecord.Add(bookRecord)
	if err != nil {
		return xerror.ErrInternalServer.SetMessage("failed to add borrow record: " + err.Error())
	}
	return nil
}
