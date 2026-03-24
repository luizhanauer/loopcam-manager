package domain

// VideoSource define a porta para leitura de frames em sequência
type VideoSource interface {
	ReadFrame() (Frame, error)
	Close() error
}
