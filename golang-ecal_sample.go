package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang-ecal/ecal"
)

func minimalSnd() {
	var pub ecal.PublisherIf
	var err error
	pub, err = ecal.PublisherCreate("Hello", "base:std::string", "", true)
	if err != nil {
		log.Fatal(err)
	}
	defer pub.Destroy()

	pubChannel := pub.GetInputChannel()
	go func() {
		count := 1
		for ecal.Ok() {
			message := ecal.Message{Content: []byte(fmt.Sprintf("HELLO WORLD FROM GO (%d)", count)),
				Timestamp: -1}
			count += 1
			select {
			case pubChannel <- message:
				log.Printf("Sent \"%s\"\n", message.Content)
			case <-time.After(time.Second):
			}
			<-time.After(250 * time.Millisecond)
		}

		pub.Stop()
	}()

	for !pub.IsStopped() {
		<-time.After(time.Second)
	}
}

func minimalRec() {
	var sub ecal.SubscriberIf
	var err error
	sub, err = ecal.SubscriberCreate("Hello", "base:std::string", "", true, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Destroy()

	subChannel := sub.GetOutputChannel()
	go func() {
		for ecal.Ok() {
			select {
			case message := <-subChannel:
				log.Printf("Received \"%s\"\n", message.Content)
			case <-time.After(time.Second):
			}
		}

		sub.Stop()
	}()

	for !sub.IsStopped() {
		<-time.After(time.Second)
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
			mode = "minimal_snd"
		case "minimal_rec":
			mode = "minimal_rec"
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
