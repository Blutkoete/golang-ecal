package ecal

/*
#include <stdlib.h>
 */
import "C"
import (
	"errors"
	"log"
	"sync"
	"unsafe"

	"github.com/Blutkoete/golang-ecal/ecalc"
)

type PublisherIf interface {
	GetHandle() uintptr

	Start() error
	Stop() error
	Stopped() bool

	Destroy() error
	Destroyed() bool

	send(message Message) error

	GetChannel() chan<- Message
	GetTopic() string
	GetTopicType() string
	GetTopicDesc() string
}

type publisher struct {
	handle uintptr
	running bool
	destroyed bool
	input chan Message
	topicName string
	topicType string
	topicDesc string
	mutex *sync.Mutex
}

func (pub publisher) GetHandle() uintptr {
	return pub.handle
}

func (pub publisher) Start() error {
	log.Println("Start() waiting for lock ...")
	pub.mutex.Lock()
	defer pub.mutex.Unlock()

	if pub.destroyed {
		return errors.New("publisher already destroyed")
	}
	pub.running = true

	go func(pub publisher) {
		for !pub.destroyed && pub.running {
			message := <-pub.input

			log.Println("goroutine waiting for lock ...")
			pub.mutex.Lock()
			if !pub.destroyed && pub.running {
				pub.mutex.Unlock()
				pub.send(message)
			} else {
				pub.mutex.Unlock()
			}
		}
	}(pub)
	return nil
}

func (pub publisher) Stop() error {
	log.Println("Stop() waiting for lock ...")
	pub.mutex.Lock()
	defer pub.mutex.Unlock()

	if pub.destroyed {
		return errors.New("publisher already destroyed")
	}

	pub.running = false
	return nil
}

func (pub publisher) Stopped() bool {
	return !pub.running
}

func (pub publisher) Destroy() error {
	if pub.running {
		pub.Stop()
	}

	log.Println("Destroy() waiting for lock ...")
	pub.mutex.Lock()
	defer pub.mutex.Unlock()

	rc := ecalc.ECAL_Pub_Destroy(pub.handle)
	if rc == 0 {
		return errors.New("could not destroy publisher")
	}

	pub.destroyed = true
	return nil
}

func (pub publisher) Destroyed() bool {
	pub.mutex.Lock()
	defer pub.mutex.Unlock()

	return pub.destroyed
}

func (pub publisher) send(message Message) error {
	log.Println("Send() waiting for lock ...")
	pub.mutex.Lock()
	defer pub.mutex.Unlock()

	if pub.destroyed {
		return errors.New("publisher already destroyed")
	}

	if !pub.running {
		return errors.New("publisher stopped")
	}

	bytesSent := ecalc.ECAL_Pub_Send(pub.handle, (uintptr)(unsafe.Pointer(&message.Content[0])), len(message.Content), message.Timestamp)
	if bytesSent < len(message.Content) {
		log.Println("error sending")
		return errors.New("error sending")
	}

	log.Println("success sending", bytesSent, "bytes")
	return nil
}

func (pub publisher) GetChannel() chan<- Message {
	return pub.input
}

func (pub publisher) GetTopic() string {
	return pub.topicName
}

func (pub publisher) GetTopicType() string {
	return pub.topicType
}

func (pub publisher) GetTopicDesc() string {
	return pub.topicDesc
}

func PublisherCreate(topicName string, topicType string, topicDesc string, start bool) (PublisherIf, error) {
	handle := ecalc.ECAL_Pub_New()
	if handle == 0 {
		return nil, errors.New("could not create new publisher")
	}

	rc := ecalc.ECAL_Pub_Create(handle, topicName, topicType, topicDesc, len(topicDesc))
	if rc == 0 {
		return nil, errors.New("could not create new publisher")
	}

	pub := publisher { handle:    handle,
		               running:   false,
		               destroyed: false,
		               input:     make(chan Message),
		               topicName: topicName,
		               topicType: topicType,
		               topicDesc: topicDesc,
		               mutex:     &sync.Mutex{} }
	if start {
		err := pub.Start()
		if err != nil {
			return nil, err
		}
	}

	return pub, nil
}