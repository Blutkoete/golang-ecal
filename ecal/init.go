package ecal

import "C"
import (
	"errors"
	"unsafe"

	"github.com/Blutkoete/golang-ecal/ecalc"
)

const InitPublisher uint = uint(ecalc.ECAL_Init_Publisher)
const InitSubscriber uint = uint(ecalc.ECAL_Init_Subscriber)
const InitService uint = uint(ecalc.ECAL_Init_Service)
const InitMonitoring uint = uint(ecalc.ECAL_Init_Monitoring)
const InitLogging uint = uint(ecalc.ECAL_Init_Logging)
const InitTimeSync uint = uint(ecalc.ECAL_Init_TimeSync)
const InitRPC uint = uint(ecalc.ECAL_Init_RPC)
const InitProcessReg uint = uint(ecalc.ECAL_Init_ProcessReg)

const InitAll = InitPublisher | InitSubscriber | InitService | InitMonitoring | InitLogging | InitTimeSync | InitRPC | InitProcessReg

const InitDefault = InitPublisher | InitSubscriber | InitService | InitLogging | InitTimeSync | InitProcessReg

func Initialize(args []string, unitName string, components uint) error {
	cArgs := C.malloc(C.size_t(len(args)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	goArgs := (*[1<<30 - 1]*C.char)(cArgs)
	for idx, arg := range args {
		goArgs[idx] = C.CString(arg)
	}

	switch rc := ecalc.ECAL_Initialize(len(args), (*string)(cArgs), unitName, components); rc {
	case 0:
		return nil
	case 1:
		return errors.New("already initialized")
	default:
		return errors.New("failed")
	}
}

func Finalize(components uint) error {
	switch rc := ecalc.ECAL_Finalize(components); rc {
	case 0:
		return nil
	case 1:
		return errors.New("already finalized")
	default:
		return errors.New("failed")
	}
}

func Ok() bool {
	return ecalc.ECAL_Ok() != 0
}
