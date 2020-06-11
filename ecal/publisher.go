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
	Start() error
	Stop() error
	Destroy() error

	IsStopped() bool
	IsDestroyed() bool
	IsSubscribed() bool

	GetHandle() uintptr
	GetInputChannel() chan<- Message
	GetEventChannel() <-chan bool
	GetTopic() string
	GetType() string
	GetDescription() string
	GetQoS() error
	GetLayerMode() error
	GetMaxBandwidthUDP() int64
	GetID() int64

	SetDescription(topicDesc string) error
	SetQoS() error
	SetLayerMode() error
	SetMaxBandwidthUDP(bandwidth int64) error
	SetID(id int64) error

	ShareType(state int) error
	ShareDescription(state int) error

	Dump() ([]byte, error)

	send(message Message) error
}

type publisher struct {
	handle          uintptr
	running         bool
	destroyed       bool
	inputSource     chan Message
	eventSink       chan bool
	topicName       string
	topicType       string
	topicDesc       string
	maxBandwidthUDP int64
	id              int64
	mutex           *sync.Mutex
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
			message := <-pub.inputSource

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

func (pub publisher) IsStopped() bool {
	return !pub.running
}

func (pub publisher) IsDestroyed() bool {
	pub.mutex.Lock()
	defer pub.mutex.Unlock()

	return pub.destroyed
}

func (pub publisher) IsSubscribed() bool {
	pub.mutex.Lock()
	defer pub.mutex.Unlock()

	return ecalc.ECAL_Pub_IsSubscribed(pub.handle) != 0
}

func (pub publisher) GetHandle() uintptr {
	return pub.handle
}

func (pub publisher) GetInputChannel() chan<- Message {
	return pub.inputSource
}

func (pub publisher) GetEventChannel() <-chan bool {
	return pub.eventSink
}

func (pub publisher) GetTopic() string {
	return pub.topicName
}

func (pub publisher) GetType() string {
	return pub.topicType
}

func (pub publisher) GetDescription() string {
	return pub.topicDesc
}

func (pub publisher) GetQoS() error {
	return errors.New("not implemented")
}

func (pub publisher) GetLayerMode() error {
	return errors.New("not implemented")
}

func (pub publisher) GetMaxBandwidthUDP() int64 {
	return pub.maxBandwidthUDP
}

func (pub publisher) GetID() int64 {
	return pub.id
}

func (pub publisher) SetDescription(topicDesc string) error {
	rc := ecalc.ECAL_Pub_SetDescription(pub.handle, topicDesc, len(topicDesc))
	if rc == 0 {
		return errors.New("setting description failed")
	}
	pub.topicDesc = topicDesc
	return nil
}

func (pub publisher) SetQoS() error {
	return errors.New("not implemented")
}

func (pub publisher) SetLayerMode() error {
	return errors.New("not implemented")
}

func (pub publisher) SetMaxBandwidthUDP(bandwidth int64) error {
	rc := ecalc.ECAL_Pub_SetMaxBandwidthUDP(pub.handle, bandwidth)
	if rc == 0 {
		return errors.New("setting maximum UDP bandwith failed")
	}
	pub.maxBandwidthUDP = bandwidth
	return nil
}

func (pub publisher) SetID(id int64) error {
	rc := ecalc.ECAL_Pub_SetID(pub.handle, id)
	if rc == 0 {
		return errors.New("setting ID failed")
	}
	pub.id = id
	return nil
}

func (pub publisher) ShareType(state int) error {
	return errors.New("not implemented")
}

func (pub publisher) ShareDescription(state int) error {
	return errors.New("not implemented")
}

func (pub publisher) Dump() ([]byte, error) {
	const bufferSize = 4096
	cBuffer := C.malloc(bufferSize)
	defer C.free(cBuffer)

	bytesInDump := ecalc.ECAL_Pub_Dump(pub.handle, (uintptr)(cBuffer), bufferSize)
	if bytesInDump <= 0 {
		return nil, errors.New("dump failed")
	}

	dump := make([]byte, bytesInDump, bytesInDump)
	gBuffer := (*[1 << 30]byte)(cBuffer)
	copy(dump, gBuffer[:bytesInDump])
	return dump, nil
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

func PublisherCreate(topicName string, topicType string, topicDesc string, start bool) (PublisherIf, error) {
	handle := ecalc.ECAL_Pub_New()
	if handle == 0 {
		return nil, errors.New("could not create new publisher")
	}

	rc := ecalc.ECAL_Pub_Create(handle, topicName, topicType, topicDesc, len(topicDesc))
	if rc == 0 {
		return nil, errors.New("could not create new publisher")
	}

	pub := publisher{handle: handle,
		running:         false,
		destroyed:       false,
		inputSource:     make(chan Message),
		eventSink:       make(chan bool),
		topicName:       topicName,
		topicType:       topicType,
		topicDesc:       topicDesc,
		maxBandwidthUDP: -1,
		id:              -1,
		mutex:           &sync.Mutex{}}
	if start {
		err := pub.Start()
		if err != nil {
			return nil, err
		}
	}

	return pub, nil
}
