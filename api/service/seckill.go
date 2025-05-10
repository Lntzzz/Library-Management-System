package service

import (
	"Library-Management-System/api/constant"
	"Library-Management-System/api/dto/request"
	"Library-Management-System/api/dto/response"
	"Library-Management-System/api/idgen"
	"Library-Management-System/api/service/messages"
	"Library-Management-System/api/util/xerror"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var Seckill *SeckillService

var ctx = context.Background()

func init() {
	Seckill = NewSeckillService()
}

type SeckillService struct {
	RedisClient *redis.Client
}

func NewSeckillService() *SeckillService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // Redis 密码
		DB:       0,                // Redis 数据库索引
	})

	return &SeckillService{
		RedisClient: redisClient,
	}
}

type SeckillMessage struct {
	Id     string `json:"id"`
	UserID string `json:"user_id"`
	BookID string `json:"book_id"`
	Time   string `json:"time"`
}

func (s *SeckillService) Seckill(option *request.SecKillBorrowRecordOption) (*response.SecKillBorrowRecordResponse, xerror.OpenError) {
	userID := option.UserId
	bookID := option.BookId
	// Redis库存键和用户集合键
	stockKey := "book_stock:" + bookID
	userSetKey := "book_users:" + bookID

	// 检查库存是否充足
	stock, err := s.RedisClient.Get(ctx, stockKey).Int()
	if err != nil {
		return nil, xerror.ErrInternalServer.SetMessage("failed to get stock from redis: " + err.Error())
	}
	if stock <= 0 {
		return nil, xerror.ErrInternalServer.SetMessage("stock not sufficient")
	}

	// 检查用户是否已下单
	isMember, err := s.RedisClient.SIsMember(ctx, userSetKey, userID).Result()
	if err != nil {
		return nil, xerror.ErrInternalServer.SetMessage("failed to check user in redis: " + err.Error())
	}
	if isMember {
		return nil, xerror.ErrInternalServer.SetMessage("user already placed an order")
	}

	// 扣减库存
	_, err = s.RedisClient.Decr(ctx, stockKey).Result()
	if err != nil {
		return nil, xerror.ErrInternalServer.SetMessage("failed to decrement stock: " + err.Error())
	}

	// 将用户ID加入已下单集合
	_, err = s.RedisClient.SAdd(ctx, userSetKey, userID).Result()
	if err != nil {
		return nil, xerror.ErrInternalServer.SetMessage("failed to add user to set: " + err.Error())
	}

	//	script := `
	//local stockKey = KEYS[1]
	//local userSetKey = KEYS[2]
	//local userID = ARGV[1]
	//
	//-- 检查库存是否充足
	//local stock = tonumber(redis.call("GET", stockKey))
	//if not stock or stock <= 0 then
	//    return 1 -- 库存不足
	//end
	//
	//-- 检查用户是否已下单
	//local isMember = redis.call("SISMEMBER", userSetKey, userID)
	//if isMember == 1 then
	//    return 2 -- 用户已下单
	//end
	//
	//-- 扣减库存
	//redis.call("DECR", stockKey)
	//
	//-- 将用户ID加入已下单集合
	//redis.call("SADD", userSetKey, userID)
	//
	//return 0 -- 正常结束操作
	//`
	//
	//	stockKey := "book_stock:" + bookID
	//	userSetKey := "book_users:" + bookID
	//
	//	result, err := s.RedisClient.Eval(ctx, script, []string{stockKey, userSetKey}, userID).Int()
	//	if err != nil {
	//		return nil, xerror.ErrInternalServer.SetMessage("failed to execute lua script: " + err.Error())
	//	}
	//
	//	switch result {
	//	case 1:
	//		return nil, xerror.ErrInternalServer.SetMessage("stock not sufficient")
	//	case 2:
	//		return nil, xerror.ErrInternalServer.SetMessage("user already placed an order")
	//	case 0:
	//		// 正常结束操作
	//	}

	// 构造消息并发送到 Kafka
	message := messages.SeckillMessage{
		Id:     idgen.GenBorrowRecordId(),
		UserID: userID,
		BookID: bookID,
		Time:   time.Now().Format(time.RFC3339),
	}
	messages.GetSecKillProducer().Send(&message)
	return &response.SecKillBorrowRecordResponse{BorrowRecordId: message.Id, SecKillResult: constant.SeckillSuccess}, nil
}

//func (s *SeckillService) KafkaConsumerSetup() {
//
//	// 创建 Kafka Reader
//	reader := kafka.NewReader(kafka.ReaderConfig{
//		Brokers: []string{"localhost:9092"}, // Kafka 地址
//		Topic:   "seckill",                  // Kafka 主题
//		GroupID: "seckill-group",            // 消费者组 ID
//	})
//	defer reader.Close()
//
//	for {
//		// 等待消息到达
//		msg, err := reader.FetchMessage(ctx)
//		if err != nil {
//			fmt.Println("failed to fetch message from kafka: " + err.Error())
//			continue
//		}
//
//		// 处理消息
//		err = s.ProcessKafkaMessage(msg)
//		if err != nil {
//			// 处理消息失败，记录错误
//			fmt.Println("Process message failed:", err.Error())
//		}
//
//		// 提交消息偏移量
//		if err := reader.CommitMessages(ctx, msg); err != nil {
//			fmt.Println("failed to commit message: " + err.Error())
//		}
//	}
//}
//
//func (s *SeckillService) ProcessKafkaMessage(msg kafka.Message) xerror.OpenError {
//	// 解析消息内容
//	var seckillMessage SeckillMessage
//	err := json.Unmarshal(msg.Value, &seckillMessage)
//	if err != nil {
//		return xerror.ErrInternalServer.SetMessage("failed to unmarshal kafka message: " + err.Error())
//	}
//
//	//数据库落库
//	parsedTime, err := time.Parse(time.RFC3339, seckillMessage.Time)
//	if err != nil {
//		return xerror.ErrInternalServer.SetMessage("failed to parse time: " + err.Error())
//	}
//	bookRecord := model.BorrowRecord{
//		Id:         idgen.GenBorrowRecordId(),
//		UserId:     seckillMessage.UserID,
//		BookId:     seckillMessage.BookID,
//		BorrowedAt: parsedTime,
//	}
//	err = db.BorrowRecord.Add(bookRecord)
//	if err != nil {
//		return xerror.ErrInternalServer.SetMessage("failed to add borrow record: " + err.Error())
//	}
//
//	return nil
//}
