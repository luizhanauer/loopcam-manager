package application

import (
	"context"
	"fmt"
	"io"

	"loop-cam/domain"
	"loop-cam/infra"
)

// StreamToDriver delega a I/O bruta e mantém o fluxo de negócio do loop
func StreamToDriver(ctx context.Context, videoPath string, driver domain.VideoDriver, width, height int) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := processVideo(ctx, videoPath, driver, width, height)
			if err != nil {
				return err
			}
		}
	}
}

func processVideo(ctx context.Context, videoPath string, driver domain.VideoDriver, width, height int) error {
	source, err := infra.NewFFmpegSource(videoPath, width, height)
	if err != nil {
		return err
	}
	defer source.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			frame, err := source.ReadFrame()

			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("erro lendo frame: %w", err)
			}

			if err := driver.Write(frame); err != nil {
				return fmt.Errorf("erro escrevendo frame: %w", err)
			}
		}
	}
}
