package kafkacli

import (
	cfg "api/config"
	"github.com/IBM/sarama"
	"strings"
	"sync"
)

var KafkaPool = sync.Pool{
	New: func() interface{} {
		config := sarama.NewConfig()
		config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
		config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
		config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
		broker := strings.Split(cfg.KafkaBroker, ",")
		var err error
		kafkaclient, err := sarama.NewSyncProducer(broker, config)
		if err != nil {
			panic(err)
		}
		return kafkaclient
	},
}

func PushData2Kafka(msg string) error {
	cli := KafkaPool.Get().(sarama.SyncProducer)
	defer KafkaPool.Put(cli)

	producerMsg := &sarama.ProducerMessage{}
	producerMsg.Topic = cfg.Topic
	producerMsg.Value = sarama.StringEncoder(msg)

	_, _, err := cli.SendMessage(producerMsg)
	return err

}
