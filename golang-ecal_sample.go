package main

import (
	"log"
	"os"
	"time"

	"golang-ecal/ecal"
)

func main() {

	err := ecal.Initialize(os.Args, "golang-ecal_sample", ecal.InitDefault)
	if err != nil {
		log.Fatal(err)
	}
	defer ecal.Finalize(ecal.InitAll)

	var pub ecal.PublisherIf
	pub, err = ecal.PublisherCreate("Hello", "base:std::string", "", true)
	if err != nil {
		log.Fatal(err)
	}
	defer pub.Destroy()

	<-time.After(time.Second)

	pubChannel := pub.GetInputChannel()
	go func(channel chan<- ecal.Message) {
		message := ecal.Message { Content: []byte("HELLO WORLD FROM GO"),
			                      Timestamp: -1 }
		for ecal.Ok() {
			select {
			case channel <- message:
			case <- time.After(time.Second):
			}

			<- time.After(250 * time.Millisecond)
		}

		pub.Stop()
	}(pubChannel)

	for !pub.IsStopped() {
		<- time.After(time.Second)
	}

}
