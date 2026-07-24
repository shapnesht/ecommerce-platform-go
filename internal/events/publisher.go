package events

type Publisher interface {
	Publish(eventName string, payload interface{}, metadata map[string]string) error
	Close() error
}
