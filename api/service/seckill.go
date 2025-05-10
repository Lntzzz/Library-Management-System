package service

import (
	"Library-Management-System/api/constant"
	"Library-Management-System/api/db"
	"Library-Management-System/api/dto/request"
	"Library-Management-System/api/dto/response"
	"Library-Management-System/api/idgen"
	"Library-Management-System/api/model"
	"Library-Management-System/api/util/xerror"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
)

var Seckill *SeckillService

var ctx = context.Background()

func init() {
	Seckill = NewSeckillService()
}

type SeckillService struct {
	RedisClient *redis.Client
	KafkaWriter *kafka.Writer
}

func NewSeckillService() *SeckillService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // Redis 密码
		DB:       0,                // Redis 数据库索引
	})

	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"), // Kafka 地址
		Topic:    "seckill",                   // Kafka 主题
		Balancer: &kafka.Hash{},               // 分区分配器
	}

	return &SeckillService{
		RedisClient: redisClient,
		KafkaWriter: kafkaWriter,
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

	// 构造消息并发送到 Kafka
	message := SeckillMessage{
		Id:     idgen.GenBorrowRecordId(),
		UserID: userID,
		BookID: bookID,
		Time:   time.Now().Format(time.RFC3339),
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, xerror.ErrInternalServer.SetMessage("failed to marshal message: " + err.Error())
	}

	err = s.KafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(userID),
		Value: messageBytes,
	})
	if err != nil {
		return nil, xerror.ErrInternalServer.SetMessage("failed to send message to kafka: " + err.Error())
	}

	return &response.SecKillBorrowRecordResponse{BorrowRecordId: message.Id, SecKillResult: constant.SeckillSuccess}, nil
}

func (s *SeckillService) ConsumeKafkaMessages() {
	// 创建上下文，用于优雅退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 捕获系统信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		fmt.Println("Received shutdown signal, exiting...")
		cancel()
	}()

	// 创建 Kafka Reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"}, // Kafka 地址
		Topic:   "seckill",                  // Kafka 主题
		GroupID: "seckill-group",            // 消费者组 ID
	})
	defer reader.Close()

	for {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down Kafka consumer...")
			return
		default:
			// 读取消息
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				fmt.Println("failed to read message from kafka: " + err.Error())
				continue
			}

			// 处理消息
			err = s.ProcessKafkaMessage(msg)
			if err != nil {
				// 处理消息失败，记录错误
				fmt.Println("Process message failed:", err.Error())
			}
		}
	}
}

func (s *SeckillService) ProcessKafkaMessage(msg kafka.Message) xerror.OpenError {
	// 解析消息内容
	var seckillMessage SeckillMessage
	err := json.Unmarshal(msg.Value, &seckillMessage)
	if err != nil {
		return xerror.ErrInternalServer.SetMessage("failed to unmarshal kafka message: " + err.Error())
	}

	//数据库落库
	parsedTime, err := time.Parse(time.RFC3339, seckillMessage.Time)
	if err != nil {
		return xerror.ErrInternalServer.SetMessage("failed to parse time: " + err.Error())
	}
	bookRecord := model.BorrowRecord{
		Id:         idgen.GenBorrowRecordId(),
		UserId:     seckillMessage.UserID,
		BookId:     seckillMessage.BookID,
		BorrowedAt: parsedTime,
	}
	err = db.BorrowRecord.Add(bookRecord)
	if err != nil {
		return xerror.ErrInternalServer.SetMessage("failed to add borrow record: " + err.Error())
	}

	return nil
}
