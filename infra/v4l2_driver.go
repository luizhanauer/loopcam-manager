package infra

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"

	// Adapte para o nome do seu módulo real configurado no go.mod
	"loop-cam/domain"
)

const (
	// IOCTLs específicos para Kernel x86_64 (208 bytes de payload)
	VIDIOC_S_FMT = 0xc0d05605 // Set Format
	VIDIOC_G_FMT = 0xc0d05604 // Get Format

	V4L2_BUF_TYPE_OUTPUT = 2
	V4L2_PIX_FMT_YUYV    = 0x56595559
)

// v4l2PixFormat reflete o struct v4l2_pix_format em C (48 bytes)
type v4l2PixFormat struct {
	Width        uint32
	Height       uint32
	PixelFormat  uint32
	Field        uint32
	BytesPerLine uint32
	SizeImage    uint32
	Colorspace   uint32
	Priv         uint32
	Flags        uint32
	YcbcrEnc     uint32
	Quantization uint32
	XferFunc     uint32
}

// v4l2Format reflete o struct v4l2_format do Kernel Linux em 64-bits
type v4l2Format struct {
	Type uint32
	_    uint32 // Padding de 4 bytes para alinhamento (Memory Padding)

	Pix v4l2PixFormat

	_ [152]byte // Restante do union original em C para fechar 208 bytes
}

type LinuxV4L2Driver struct {
	device *os.File
}

func NewLinuxV4L2Driver(path string) (*LinuxV4L2Driver, error) {
	// Abertura restrita a O_WRONLY garante compatibilidade com exclusive_caps=1
	// Isso permite que navegadores reconheçam o device como Webcam.
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir device %s: %w", path, err)
	}

	return &LinuxV4L2Driver{device: f}, nil
}

func (d *LinuxV4L2Driver) Setup(width, height int) error {
	f := v4l2Format{
		Type: V4L2_BUF_TYPE_OUTPUT,
		Pix: v4l2PixFormat{
			Width:        uint32(width),
			Height:       uint32(height),
			PixelFormat:  V4L2_PIX_FMT_YUYV,
			Field:        1, // V4L2_FIELD_NONE (Vídeo progressivo)
			BytesPerLine: uint32(width * 2),
			SizeImage:    uint32(width * height * 2),
			Colorspace:   8, // V4L2_COLORSPACE_SRGB
		},
	}

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, d.device.Fd(), VIDIOC_S_FMT, uintptr(unsafe.Pointer(&f)))

	// Early return em caso de sucesso absoluto
	if errno == 0 {
		return nil
	}

	// Se o erro for diferente de EINVAL (22), falhamos imediatamente
	if errno != 22 {
		return fmt.Errorf("ioctl VIDIOC_S_FMT falhou (errno %d)", errno)
	}

	// Fallback para lidar com device travado pelo ciclo de vida do Wails
	return d.checkAndReuseFormat(width, height)
}

func (d *LinuxV4L2Driver) checkAndReuseFormat(expectedWidth, expectedHeight int) error {
	var currentFmt v4l2Format
	currentFmt.Type = V4L2_BUF_TYPE_OUTPUT

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, d.device.Fd(), VIDIOC_G_FMT, uintptr(unsafe.Pointer(&currentFmt)))

	if errno != 0 {
		return fmt.Errorf("device rejeitou S_FMT e falhou ao ler G_FMT (errno %d)", errno)
	}

	if currentFmt.Pix.Width != uint32(expectedWidth) || currentFmt.Pix.Height != uint32(expectedHeight) {
		return fmt.Errorf(
			"device travado em resolução incompatível (%dx%d). Mate os processos lendo o device ou recarregue o módulo do kernel",
			currentFmt.Pix.Width,
			currentFmt.Pix.Height,
		)
	}

	// Resolução bate com o esperado, podemos ignorar o erro do Setup inicial
	fmt.Printf("[V4L2] Device %s previamente alocado. Reutilizando formato existente.\n", d.device.Name())
	return nil
}

func (d *LinuxV4L2Driver) Write(frame domain.Frame) error {
	_, err := d.device.Write(frame.Data)
	return err
}

func (d *LinuxV4L2Driver) Close() error {
	return d.device.Close()
}
