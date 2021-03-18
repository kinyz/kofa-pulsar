package kofa

import (
	"kofa-pulsar/core/message"
	"log"
	"strconv"
	"time"
)

type Discovery interface {
	Initialize(server Server)
	Ping()
	Register()
	Logout()
	Check()
	Close()
}

type IDiscovery struct {
	Server
}

const DiscoveryTopic = "Kofa-Discovery"

func NewDiscovery() Discovery {
	return &IDiscovery{}
}

func (d *IDiscovery) Initialize(server Server) {
	d.Server = server
	d.Server.AddRouter(100, "Discovery", 0, &DiscoveryListen{serverId: make(map[uint64]map[string]bool)}, d.Server)

}

func (d *IDiscovery) Register() {
	serMap := d.ServiceManager().GetAll()
	msg := message.NewMessage()
	msg.Header().SetMsgId(101)

	log.Println("discovery register begin...")

	for _, v := range serMap {
		if v.GetMsgId() <= 103 && v.GetMsgId() > 100 {
			continue
		}

		msg.Property().Set("MsgId", strconv.Itoa(int(v.GetMsgId())))
		msg.Property().Set("Topic", v.GetTopic())
		msg.Property().Set("Alias", v.GetAlias())
		msg.Property().Set("Method", v.GetMethod())
		msg.Property().Set("Level", strconv.Itoa(int(v.GetLevel())))
		msg.Header().SetMsgId(101)
		err := d.Send(msg, DiscoveryTopic)

		if err != nil {
			log.Println(err)
		}
	}
}

func (d *IDiscovery) Logout() {
	serMap := d.ServiceManager().GetAll()
	msg := message.NewMessage()
	msg.Header().SetMsgId(103)
	//b64 := make([]byte, 20)
	for _, v := range serMap {
		//binary.BigEndian.PutUint64(b64, v.GetMsgId())
		msg.Property().Set("MsgId", strconv.Itoa(int(v.GetMsgId())))
		err := d.Send(msg, DiscoveryTopic)
		if err != nil {
			log.Println("discovery logout err:", err)
		}
	}
}

func (d *IDiscovery) Check() {
	msg := message.NewMessage()
	msg.Header().SetMsgId(102)
	err := d.Send(msg, DiscoveryTopic)
	if err != nil {
		log.Println("discovery check err:", err)
	}
}

func (d *IDiscovery) Close() {
	d.Logout()
}
func (d *IDiscovery) Ping() {
	go func() {

		//	services:=d.ServiceManager().GetAll()
		msg := message.NewMessage()
		msg.Header().SetMsgId(104)
		d.Send(msg, DiscoveryTopic)
		time.Sleep(time.Second * 30)
	}()
}

type DiscoveryListen struct {
	serverId map[uint64]map[string]bool
}

func (dl *DiscoveryListen) Add(request Request, server Server) {

	msgId, err := request.Message().Property().Get("MsgId")

	if err != nil {
		log.Println("discovery add service err1:", err)
	}
	topic, err2 := request.Message().Property().Get("Topic")
	if err2 != nil {
		log.Println("discovery add service err2:", err)
	}
	alias, err3 := request.Message().Property().Get("Alias")
	if err3 != nil {
		log.Println("discovery add service err3:", err)
	}
	method, err4 := request.Message().Property().Get("Method")
	if err4 != nil {
		log.Println("discovery add service err4:", err)
	}
	level, err5 := request.Message().Property().Get("Level")

	if err5 != nil {
		log.Println("discovery add service err5:", err)
	}

	intNum, _ := strconv.Atoi(msgId)
	intNumlv, _ := strconv.Atoi(level)
	if _, ok := dl.serverId[uint64(intNum)]; !ok {
		dl.serverId[uint64(intNum)] = make(map[string]bool)
	}

	dl.serverId[uint64(intNum)][request.Message().Header().GetProducer()] = true
	_ = server.ServiceManager().Add(uint64(intNum), topic, alias, method, uint32(intNumlv))

}

func (dl *DiscoveryListen) Logout(request Request, server Server) {

	if data, err := request.Message().Property().Get("MsgId"); err != nil {
		log.Println("discovery logout service err:", err)
	} else {
		intNum, _ := strconv.Atoi(data)
		msgId := uint64(intNum)


		if _, ok := dl.serverId[msgId][request.Message().Header().GetProducer()]; ok {
			delete(dl.serverId[msgId], request.Message().Header().GetProducer())
		}
		if len(dl.serverId[msgId]) <= 0 {
			switch msgId {
			case 101:
				return
			case 102:
				return
			case 103:
				return
			case 104:
				return

			}
			err := server.ServiceManager().Remove(msgId)
			if err != nil {
				log.Println(err)
			}
		}
	}

}

func (dl *DiscoveryListen) GetAllService(request Request, server Server) {

	serMap := server.ServiceManager().GetAll()
	request.Message().Header().SetMsgId(101)
	//b64 := make([]byte, 20)
	//b32 := make([]byte, 10)
	for _, v := range serMap {
		//binary.BigEndian.PutUint64(b64, v.GetMsgId())
		//binary.BigEndian.PutUint32(b32, v.GetLevel())

		request.Message().Property().Set("MsgId", strconv.Itoa(int(v.GetMsgId())))
		request.Message().Property().Set("Topic", v.GetTopic())
		request.Message().Property().Set("Alias", v.GetAlias())
		request.Message().Property().Set("Method", v.GetMethod())
		request.Message().Property().Set("Level", strconv.Itoa(int(v.GetLevel())))
		request.Message().Header().SetMsgId(101)
		//log.Print("查询人:",request.Message().Header().GetProducer())
		err := request.Send(request.Message(), request.Message().Header().GetProducer())
		if err != nil {
			log.Println(err)
		}
	}

}

func (dl *DiscoveryListen) Ping(request Request, server Server) {

	for k, v := range dl.serverId {

		for sn, _ := range v {
			if sn == request.Message().Header().GetProducer() {
				dl.serverId[k][request.Message().Header().GetProducer()] = true

			}
		}

	}
}
