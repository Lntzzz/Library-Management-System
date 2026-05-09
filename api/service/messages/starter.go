package messages

var kafkaConf = KafkaConf{
	Brokers:              "localhost:9092",
	GroupId:              "1",
	FromOffsets:          "Newest",
	AppCode:              "test",
	ConsumerTopicSecKill: "seckill",
	ProducerTopicSecKill: "seckill",
}

func Init() {
	initSecKillConsumers(kafkaConf)
	initSecKillProducers(kafkaConf)
}

func initSecKillConsumers(conf KafkaConf) {
	SecKillConsumer(&conf).Start()
}

func initSecKillProducers(conf KafkaConf) {
	NewSecKillProducer(&conf)
}
