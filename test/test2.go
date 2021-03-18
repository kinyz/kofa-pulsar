package main

import (
	message2 "kofa-pulsar/core/message"
	"kofa-pulsar/core/prehandle"
	"kofa-pulsar/kofa"
	"log"
	"time"
)

func main(){
	message:=prehandle.Pulsar("182.254.157.106:6650")
	kf:=kofa.New("kofa-test",true,message)


	msg:=message2.NewMessage()

	msg.SetKey("im key")
	msg.Header().SetMsgId(10001)
	msg.SetData([]byte("im data"))
	msg.Property().Set("property","testdata")

	go func() {

		time.Sleep(10 * time.Second)

		err:=kf.Message().Send(msg)
	if err!=nil{
		log.Print(err)
	}
	}()
	kf.Serve()
}