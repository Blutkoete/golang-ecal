package ecal

/*
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"log"
	"os"
	"sync"

	"golang-ecal/ecalc"
)

type SubscriberIf interface {
	Start() error
	Stop() error
	Destroy() error

	IsStopped() bool
	IsDestroyed() bool

	GetHandle() uintptr
	GetBufferSize() int
	GetOutputChannel() <-chan Message
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

func (sub *subscriber) Start() error {
	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	if sub.destroyed {
		return errors.New("subscriber already destroyed")
	}
	sub.running = true

	go func(sub *subscriber) {
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

			message.Content = make([]byte, bytesReceived, bytesReceived)
			gBuffer := (*[1 << 30]byte)(cBuffer)
			copy(message.Content, gBuffer[:bytesReceived])
			sub.output <- message
		}

	}(sub)
	return nil
}

func (sub *subscriber) Stop() error {
	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	if sub.destroyed {
		return errors.New("subscriber already destroyed")
	}

	sub.running = false
	return nil
}

func (sub *subscriber) Destroy() error {
	if sub.running {
		sub.Stop()
	}

	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	rc := ecalc.ECAL_Sub_Destroy(sub.handle)
	if rc == 0 {
		return errors.New("could not destroy subscriber")
	}

	sub.destroyed = true
	return nil
}

func (sub *subscriber) IsStopped() bool {
	return !sub.running
}

func (sub *subscriber) IsDestroyed() bool {
	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	return sub.destroyed
}

func (sub *subscriber) GetHandle() uintptr {
	return sub.handle
}

func (sub *subscriber) GetBufferSize() int {
	return sub.bufferSize
}

func (sub *subscriber) GetOutputChannel() <-chan Message {
	return sub.output
}

func (sub *subscriber) GetTopic() string {
	return sub.topicName
}

func (sub *subscriber) GetTopicType() string {
	return sub.topicType
}

func (sub *subscriber) GetTopicDesc() string {
	return sub.topicDesc
}

func SubscriberCreate(topicName string, topicType string, topicDesc string, start bool, bufferSize int) (SubscriberIf, <-chan Message, error) {
	var err error
	if ecalc.ECAL_IsInitialized(InitSubscriber) == 0 {
		err = Initialize(os.Args, os.Args[0], InitSubscriber)
		if err != nil {
			return nil, nil, err
		}
	}

	if bufferSize <= 0 {
		return nil, nil, errors.New("bufferSize must be larger than zero")
	}

	handle := ecalc.ECAL_Sub_New()
	if handle == 0 {
		return nil, nil, errors.New("could not create new subscriber")
	}

	rc := ecalc.ECAL_Sub_Create(handle, topicName, topicType, topicDesc, len(topicDesc))
	if rc == 0 {
		return nil, nil, errors.New("could not create new subscriber")
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
			return nil, nil, err
		}
	}

	return &sub, sub.GetOutputChannel(), nil
}
