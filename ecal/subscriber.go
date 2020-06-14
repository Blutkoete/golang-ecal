package ecal

/*
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"log"
	"sync"

	"golang-ecal/ecalc"
)

type SubscriberIf interface {
	GetHandle() uintptr
	GetBufferSize() int

	Start() error
	Stop() error
	Stopped() bool

	Destroy() error
	Destroyed() bool

	GetChannel() <-chan Message
	GetTopic() string
	GetTopicType() string
	GetTopicDesc() string
}

type subscriber struct {
	handle     uintptr
	bufferSize int
	running    bool
	destroyed  bool
	output     chan Message
	topicName  string
	topicType  string
	topicDesc  string
	mutex      *sync.Mutex
}

func (sub subscriber) GetHandle() uintptr {
	return sub.handle
}

func (sub subscriber) GetBufferSize() int {
	return sub.bufferSize
}

func (sub subscriber) Start() error {
	log.Println("Start() waiting for lock ...")
	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	if sub.destroyed {
		return errors.New("subscriber already destroyed")
	}
	sub.running = true

	go func(sub subscriber) {
		cBuffer := C.malloc(C.ulong(sub.bufferSize))
		defer C.free(cBuffer)

		for !sub.destroyed && sub.running {
			message := Message{Content: nil,
				Timestamp: 0}
			bytesReceived := ecalc.ECAL_Sub_Receive(sub.handle, uintptr(cBuffer), sub.bufferSize, &message.Timestamp, 100)
			if bytesReceived <= 0 {
				continue
			} else if bytesReceived > sub.bufferSize {
				log.Println("received more data than pre-allocated")
				continue
			}
			log.Println("received", bytesReceived, "bytes")

			message.Content = make([]byte, bytesReceived, bytesReceived)
			gBuffer := (*[1 << 30]byte)(cBuffer)
			bytesCopied := copy(message.Content, gBuffer[:bytesReceived])
			log.Println("copied", bytesCopied, "bytes")
			sub.output <- message
		}

	}(sub)
	return nil
}

func (sub subscriber) Stop() error {
	log.Println("Stop() waiting for lock ...")
	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	if sub.destroyed {
		return errors.New("subscriber already destroyed")
	}

	sub.running = false
	return nil
}

func (sub subscriber) Stopped() bool {
	return !sub.running
}

func (sub subscriber) Destroy() error {
	if sub.running {
		sub.Stop()
	}

	log.Println("Destroy() waiting for lock ...")
	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	rc := ecalc.ECAL_Sub_Destroy(sub.handle)
	if rc == 0 {
		return errors.New("could not destroy subscriber")
	}

	sub.destroyed = true
	return nil
}

func (sub subscriber) Destroyed() bool {
	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	return sub.destroyed
}

func (sub subscriber) GetChannel() <-chan Message {
	return sub.output
}

func (sub subscriber) GetTopic() string {
	return sub.topicName
}

func (sub subscriber) GetTopicType() string {
	return sub.topicType
}

func (sub subscriber) GetTopicDesc() string {
	return sub.topicDesc
}

func SubscriberCreate(topicName string, topicType string, topicDesc string, start bool, bufferSize int) (SubscriberIf, error) {
	if bufferSize <= 0 {
		return nil, errors.New("bufferSize must be larger than zero")
	}

	handle := ecalc.ECAL_Sub_New()
	if handle == 0 {
		return nil, errors.New("could not create new subscriber")
	}

	rc := ecalc.ECAL_Sub_Create(handle, topicName, topicType, topicDesc, len(topicDesc))
	if rc == 0 {
		return nil, errors.New("could not create new subscriber")
	}

	sub := subscriber{handle: handle,
		bufferSize: bufferSize,
		running:    false,
		destroyed:  false,
		output:     make(chan Message),
		topicName:  topicName,
		topicType:  topicType,
		topicDesc:  topicDesc,
		mutex:      &sync.Mutex{}}
	if start {
		err := sub.Start()
		if err != nil {
			return nil, err
		}
	}

	return sub, nil
}
