package prehandle

import (
	"context"
	"kofa-pulsar/core/message"
	"kofa-pulsar/kofa"

	"github.com/apache/pulsar-client-go/pulsar"
	"log"
	"time"
)

type IPulsar struct {
	Addr           string
	closeFun       func()
	kofa.Server
	client pulsar.Client
}



func Pulsar(addr string) kofa.Message {
	p := &IPulsar{Addr: addr}
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://" + addr,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}
	p.client = client
	return p

}


func (p *IPulsar) Initialize(server kofa.Server) {
	p.Server = server
	p.closeFun = func() {
		//p.client.Close()
	}
}

func (p *IPulsar) Send(msg message.Message, topic ...string) error {


	ser, err := p.ServiceManager().Get(msg.Header().GetMsgId())
		//log.Print("开始发送消息",msg.Header().GetMsgId())
	if err != nil {
		return err
	}

	if len(topic) < 1 {

		//log.Print("发送topic",ser.GetTopic())
		return p.MakeMessage(msg,ser.GetTopic())
	}

	for _,v:=range topic{
		p.MakeMessage(msg,v)

	}
	return nil

}

func (p *IPulsar) MakeMessage(msg message.Message, topic string) error {

	msg.Header().SetTimesTamp(time.Now().UnixNano())
	producer, err := p.client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
		Name: p.Info().GetId(),
	})
	if err != nil {
		return err
	}
	defer producer.Close()


	msgPull := &pulsar.ProducerMessage{
		Key:     msg.GetKey(),
		Payload: msg.GetData(),
	}
	if len(msg.Property().GetMap())>0{
		msgPull.Properties = msg.Property().GetMap()
	}else{
		msgPull.Properties = make(map[string]string)
	}
	head,err:=msg.Header().Encode()
	if err!=nil{
		log.Print("发送压缩头失败",err)
	}
	msgPull.Properties["head"] = string(head)

	_, err = producer.Send(context.Background(),msgPull)
	if err != nil {
		return err
	}
	return nil

}


func (p *IPulsar) Close() {
	p.closeFun()
}

func (p *IPulsar) Listen(group string, topic ...string) error {
	go func() {
		consumer, err := p.client.Subscribe(pulsar.ConsumerOptions{
			Topics:            topic,
			SubscriptionName: p.Info().GetId(),
			Type:             0,
		})
		if err != nil {
			log.Fatal(err)
		}

		defer consumer.Close()

		for  {
			msg, err := consumer.Receive(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			kMsg := message.NewMessage()
			kMsg.Property().SetMap(msg.Properties())
			kMsg.SetKey(msg.Key())
			kMsg.SetData(msg.Payload())
			kMsg.Header().SetProducer(msg.ProducerName())
			if kMsg.Header().Decode([]byte(msg.Properties()["head"]))!=nil{


			//	log.Print("解压头出错:",err)
				continue
			}
			//log.Print("接受头",kMsg.Header().GetMsgId())
			ser, err := p.ServiceManager().Get(kMsg.Header().GetMsgId())
			if err != nil {
				log.Println("pulsar msg error id=",kMsg.Header().GetMsgId(), err)
			} else {
				_ = p.Server.Reflect().Call(
					ser.GetAlias(),
					ser.GetMethod(),
					kofa.NewRequest(p.Server, kMsg),
				)
			}

			//log.Print("msgId:",msg.ID())
			//log.Print("Key:",msg.Key())
			//log.Print("Properties:",msg.Properties())
			//log.Print("GetReplicatedFrom:",msg.GetReplicatedFrom())
			//log.Print("OrderingKey:",msg.OrderingKey())
			//log.Print("ProducerName:",msg.ProducerName())

			consumer.Ack(msg)
		}
	}()
	return nil
}
