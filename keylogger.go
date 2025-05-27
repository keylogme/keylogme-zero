package k0

type Keylogger interface {
	Read() chan InputEvent
	Close() error
}
