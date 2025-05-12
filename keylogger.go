package keylog

type Keylogger interface {
	Read() chan inputEvent
	Close() error
}
