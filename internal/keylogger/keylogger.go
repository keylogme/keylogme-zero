package keylogger

type Keylogger interface {
	Read() chan InputEvent
	Close() error
}
