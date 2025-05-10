package messages

type KafkaConf struct {
	Brokers              string `mapstructure:"brokers"`
	GroupId              string `mapstructure:"group_id"`
	FromOffsets          string `mapstructure:"from_offsets"` // 从哪个offset开始消费, 支持（Newest，Oldest）二种，默认使用Oldest
	AppCode              string `mapstructure:"app_code"`
	ConsumerTopicSecKill string `mapstructure:"sec_kill_topic"`
	ProducerTopicSecKill string `mapstructure:"sec_kill_topic"`
}
