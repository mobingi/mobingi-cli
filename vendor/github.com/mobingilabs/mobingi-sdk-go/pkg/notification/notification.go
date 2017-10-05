package notification

type Notifier interface {
	Notify([]byte) error
}
