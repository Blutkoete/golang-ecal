package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Blutkoete/golang-ecal/ecal"
)

func minimalSnd() {
	pub, pubChannel, err := ecal.PublisherCreate("Hello", "base:std::string", "", true)
	if err != nil {
		log.Fatal(err)
	}
	defer pub.Destroy()

	count := 0
	for ecal.Ok() {
		message := ecal.Message{Content: []byte(fmt.Sprintf("Hello World from Go (%d)", count)),
				Timestamp: -1}
		count += 1
		select {
		case pubChannel <- message:
			log.Printf("Sent \"%s\"\n", message.Content)
		case <-time.After(time.Second):
		}
		<-time.After(250 * time.Millisecond)
	}
}

func minimalRec() {
	sub, subChannel, err := ecal.SubscriberCreate("Hello", "base:std::string", "", true, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Destroy()

	for ecal.Ok() {
		select {
		case message := <-subChannel:
			log.Printf("Received \"%s\"\n", message.Content)
		case <-time.After(time.Second):
		}
	}
}

func main() {
	var mode string
	if len(os.Args) <= 1 {
		log.Print("No sample type given. Assuming \"minimal_snd\".\n")
		mode = "minimal_snd"
	} else {
		switch os.Args[1] {
		case "minimal_snd":
			mode = os.Args[1]
		case "minimal_rec":
			mode = os.Args[1]
		default:
			log.Printf("Unknown sample type \"%s\". Assuming \"minimal_snd\".\n", os.Args[1])
			mode = "minimal_snd"
		}
	}

	err := ecal.Initialize(os.Args, fmt.Sprintf("golang-ecal_sample_%s", mode), ecal.InitDefault)
	if err != nil {
		log.Fatal(err)
	}
	defer ecal.Finalize(ecal.InitAll)

	switch mode {
	case "minimal_snd":
		minimalSnd()
	case "minimal_rec":
		minimalRec()
	default:
		log.Fatalf("Unknown mode \"%s\".", mode)
	}
}
