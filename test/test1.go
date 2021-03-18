package main

import (
	"kofa-pulsar/core/prehandle"
	"kofa-pulsar/kofa"
	"log"
)

func main(){

	kf:=kofa.New("kofa-test1",true,prehandle.Pulsar("182.254.157.106:6650"))

	kf.AddRouter(10000,"chat",1,&KofaTest{})
	kf.Serve()
}


type KofaTest struct {

}

func (kt *KofaTest)Message (message kofa.Request){
	log.Print("通信测试",message.Message().GetKey())

	p,e:=message.Message().Property().Get("property")
	if e!=nil{
		log.Print(e)
		return
	}
	log.Print("属性",p)
}

