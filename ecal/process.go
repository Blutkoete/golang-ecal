package ecal

import (
	"github.com/Blutkoete/golang-ecal/ecalc"
)

func ProcessSleepMS(sleepTimeMs int64) {
	ecalc.ECAL_Process_SleepMS(sleepTimeMs)
}


