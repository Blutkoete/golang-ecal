package ecal

const (
	KeepLastHistoryQOS = iota
	KeepAllHistoryQOS  = iota
)

const (
	BestEffortReliability = iota
	ReliableReliability   = iota
)

type WriterQOS struct {
	HistoryKind int
	Reliability int
}

type ReaderQOS struct {
	HistoryKind int
	Reliability int
}
