package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Executor interface {
	Execute() ([]Event, error)
}

type Type int

const (
	Unknown Type = iota
	Message
	TaskSending
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
