package protocol

type Protocol interface {
	Listen() error
}
