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
	"unsafe"

	"github.com/Blutkoete/golang-ecal/ecalc"
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
	GetEventChannel() <-chan bool
	GetTopic() string
	GetType() string
	GetDescription() string
	GetQoS() (ReaderQOS, error)
	GetIDs() []int64
	GetTimeout() int

	SetQoS(qos ReaderQOS) error
	SetIDs(id []int64) error
	SetTimeout(timeout int) error

	Dump() ([]byte, error)
}

type subscriber struct {
	handle     uintptr
	bufferSize int
	running    bool
	destroyed  bool
	outputSink chan Message
	eventSink  chan bool
	topicName  string
	topicType  string
	topicDesc  string
	ids        []int64
	timeout    int
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
			sub.outputSink <- message
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
	return sub.outputSink
}

func (sub *subscriber) GetEventChannel() <-chan bool {
	log.Println("events not supported")
	return sub.eventSink
}

func (sub *subscriber) GetTopic() string {
	return sub.topicName
}

func (sub *subscriber) GetType() string {
	return sub.topicType
}

func (sub *subscriber) GetDescription() string {
	return sub.topicDesc
}

func (sub *subscriber) GetQoS() (ReaderQOS, error) {
	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	if sub.destroyed {
		return ReaderQOS{BestEffortReliability, KeepLastHistoryQOS}, errors.New("subscriber already destroyed")
	}

	var cQOS ecalc.SReaderQOSC
	rc := ecalc.ECAL_Sub_GetQOS(sub.handle, cQOS)
	if rc == 0 {
		return ReaderQOS{BestEffortReliability, KeepLastHistoryQOS}, errors.New("getting QOS failed")
	}

	return ReaderQOS{(int)(cQOS.GetReliability()), (int)(cQOS.GetHistory_kind())}, nil
}

func (sub *subscriber) GetIDs() []int64 {
	return sub.ids
}

func (sub *subscriber) GetTimeout() int {
	return sub.timeout
}

func (sub *subscriber) SetQoS(qos ReaderQOS) error {
	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	if sub.destroyed {
		return errors.New("subscriber already destroyed")
	}

	var cQOS ecalc.SReaderQOSC
	cQOS.SetReliability((ecalc.Enum_SS_eQOSPolicy_ReliabilityC)(qos.Reliability))
	cQOS.SetHistory_kind((ecalc.Enum_SS_eQOSPolicy_HistoryKindC)(qos.HistoryKind))
	rc := ecalc.ECAL_Sub_SetQOS(sub.handle, cQOS)
	if rc == 0 {
		return errors.New("setting QOS failed")
	}

	return nil
}

func (sub *subscriber) SetIDs(ids []int64) error {
	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	if sub.destroyed {
		return errors.New("subscriber already destroyed")
	}

	rc := ecalc.ECAL_Sub_SetID(sub.handle, (*int64)(unsafe.Pointer(&ids[0])), len(ids))
	if rc == 0 {
		return errors.New("setting ids failed")
	}

	sub.ids = ids
	return nil
}

func (sub *subscriber) SetTimeout(timeout int) error {
	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	if sub.destroyed {
		return errors.New("subscriber already destroyed")
	}

	rc := ecalc.ECAL_Sub_SetTimeout(sub.handle, timeout)
	if rc == 0 {
		return errors.New("setting QOS failed")
	}

	return nil
}

func (sub *subscriber) Dump() ([]byte, error) {
	sub.mutex.Lock()
	defer sub.mutex.Unlock()

	if sub.destroyed {
		return nil, errors.New("subscriber already destroyed")
	}

	const bufferSize = 4096
	cBuffer := C.malloc(bufferSize)
	defer C.free(cBuffer)

	bytesInDump := ecalc.ECAL_Sub_Dump(sub.handle, (uintptr)(cBuffer), bufferSize)
	if bytesInDump <= 0 {
		return nil, errors.New("dump failed")
	}

	dump := make([]byte, bytesInDump, bytesInDump)
	gBuffer := (*[1 << 30]byte)(cBuffer)
	copy(dump, gBuffer[:bytesInDump])
	return dump, nil
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
		outputSink: make(chan Message),
		eventSink:  make(chan bool),
		topicName:  topicName,
		topicType:  topicType,
		topicDesc:  topicDesc,
		ids:        make([]int64, 0),
		timeout:    0,
		mutex:      &sync.Mutex{}}
	if start {
		err := sub.Start()
		if err != nil {
			return nil, nil, err
		}
	}

	return &sub, sub.GetOutputChannel(), nil
}
