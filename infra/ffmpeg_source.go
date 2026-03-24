package infra

import (
	"fmt"
	"io"
	"os/exec"

	"loop-cam/domain"
)

type FFmpegSource struct {
	cmd    *exec.Cmd
	stdout io.ReadCloser
	width  int
	height int
	buffer []byte
}

func NewFFmpegSource(filePath string, width, height int) (*FFmpegSource, error) {
	// Filtro avançado para Center Crop:
	// 1. force_original_aspect_ratio=increase: Garante que o vídeo preencha toda a tela sem esticar.
	// 2. crop: Corta exatamente no tamanho desejado, centralizando a imagem.
	vfFilter := fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=increase,crop=%d:%d", width, height, width, height)

	cmd := exec.Command(
		"ffmpeg", "-re", "-v", "quiet",
		"-i", filePath,
		"-vf", vfFilter, // Substitui o "-s" pelo pipeline de filtros
		"-f", "rawvideo",
		"-pix_fmt", "yuyv422",
		"-an", "-",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("falha ao criar pipe stdout: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("falha ao iniciar ffmpeg: %w", err)
	}

	return &FFmpegSource{
		cmd:    cmd,
		stdout: stdout,
		width:  width,
		height: height,
		buffer: make([]byte, width*height*2), // YUYV usa 2 bytes por pixel
	}, nil
}

func (s *FFmpegSource) ReadFrame() (domain.Frame, error) {
	_, err := io.ReadFull(s.stdout, s.buffer)

	if err != nil {
		return domain.Frame{}, err
	}

	return domain.Frame{
		Data:   s.buffer,
		Width:  s.width,
		Height: s.height,
	}, nil
}

func (s *FFmpegSource) Close() error {
	s.stdout.Close()
	return s.cmd.Wait()
}
