package txlog

type Event struct {
    Key   string
    Value string
    Op    string
}


type Log interface {
    Append(e Event) error
}