package kofa

import (
	"kofa-pulsar/core/message"

)

type Message interface {
	/*
		通信接口入口
		参数 kofa server
	*/
	Initialize(server Server)

	/*
		需要实现的传输接口
		参数1 message的interface
		参数2 需要接收的topic list 可空
		空的情况下 需要通过通过service manager 获取topic
	*/

	Send(msg message.Message, topic ...string) error
	/*
		监听接口
		参数1 监听组名
		参数2 需要监听的topic list
	*/

	Listen(group string, topic ...string) error
	/*
		关闭接口
	*/
	Close()

	///*
	//	新增kafka自定义处理头
	//*/
	//AddCustomHandle(pre KafkaPreHandle)

	///*
	//	kafka发送消息
	//*/
	//Sync(topic string, key, data []byte, headers ...sarama.RecordHeader) (int32, int64, error)
	//Async(topic string, key, data []byte, headers ...sarama.RecordHeader)
}

//type KafkaPreHandle interface {
//	KafkaHandle(consumerMessage *sarama.ConsumerMessage)
//}
