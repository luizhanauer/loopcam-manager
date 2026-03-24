package domain

// Frame representa um buffer de imagem bruto
type Frame struct {
	Data   []byte
	Width  int
	Height int
}

// VideoDriver define a interface para interagir com o hardware/kernel
type VideoDriver interface {
	Setup(width, height int) error
	Write(frame Frame) error
	Close() error
}
