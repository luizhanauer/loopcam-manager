package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"loop-cam/application"
	"loop-cam/domain"
	"loop-cam/infra" // Lembre-se de ajustar 'myproject' para o nome do seu módulo Go
)

// CameraSession mantém o estado e a função de cancelamento isolada de cada câmera
type CameraSession struct {
	Cancel    context.CancelFunc
	IsRunning bool
}

// App orquestra as dependências e o estado do Wails
type App struct {
	ctx      context.Context
	mu       sync.Mutex
	drivers  map[string]domain.VideoDriver
	sessions map[string]*CameraSession
}

func NewApp() *App {
	return &App{
		drivers:  make(map[string]domain.VideoDriver),
		sessions: make(map[string]*CameraSession),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// SystemSetup aciona o pkexec (solicita senha do Ubuntu) e aloca N câmeras dinamicamente
func (a *App) SystemSetup(count int) ([]string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// 1. Configura o Kernel Linux dinamicamente (Infraestrutura)
	if err := infra.ConfigureVirtualCameras(count); err != nil {
		return nil, fmt.Errorf("falha na infraestrutura do kernel: %w", err)
	}

	// 2. Limpa e fecha qualquer driver/sessão antiga caso o usuário refaça o setup
	for path, driver := range a.drivers {
		_ = driver.Close()
		fmt.Printf("[Core] Hardware %s liberado.\n", path)
	}

	a.drivers = make(map[string]domain.VideoDriver)
	a.sessions = make(map[string]*CameraSession)

	var activeDevices []string

	// 3. Pré-aloca os hardwares em Go para proteger contra o WebKitGTK
	for i := 0; i < count; i++ {
		path := fmt.Sprintf("/dev/video%d", 10+i)

		driver, err := infra.NewLinuxV4L2Driver(path)
		if err != nil {
			return nil, fmt.Errorf("falha ao conectar na câmera %s: %w", path, err)
		}

		if err := driver.Setup(640, 480); err != nil {
			driver.Close()
			return nil, fmt.Errorf("kernel rejeitou setup inicial da câmera %s: %w", path, err)
		}

		a.drivers[path] = driver
		a.sessions[path] = &CameraSession{IsRunning: false}
		activeDevices = append(activeDevices, path)

		fmt.Printf("[Core] Hardware %s trancado e configurado (640x480).\n", path)
	}

	return activeDevices, nil
}

// SelectVideoFile abre o diálogo nativo do SO para escolher o arquivo MP4
func (a *App) SelectVideoFile() (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "Selecione o arquivo de vídeo",
		Filters: []runtime.FileFilter{
			{DisplayName: "Vídeos (*.mp4)", Pattern: "*.mp4"},
		},
	}
	return runtime.OpenFileDialog(a.ctx, options)
}

// StartCamera inicia a transmissão de um vídeo para um dispositivo específico
func (a *App) StartCamera(devicePath, videoPath string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	session, exists := a.sessions[devicePath]
	if !exists {
		return fmt.Errorf("câmera %s não inicializada no sistema", devicePath)
	}

	if session.IsRunning {
		return fmt.Errorf("a câmera %s já está em execução", devicePath)
	}

	if videoPath == "" {
		return fmt.Errorf("caminho do vídeo é obrigatório")
	}

	driver := a.drivers[devicePath]

	// Cria um contexto de cancelamento isolado para esta câmera
	ctx, cancel := context.WithCancel(context.Background())

	session.Cancel = cancel
	session.IsRunning = true

	// Inicia a I/O pesada em uma Goroutine dedicada (Application Layer)
	go func(dp string) {
		err := application.StreamToDriver(ctx, videoPath, driver, 640, 480)
		if err != nil {
			fmt.Printf("Stream da %s encerrado/falhou: %v\n", dp, err)
		}

		// Libera a trava de execução quando o stream termina (ou é cancelado)
		a.mu.Lock()
		a.sessions[dp].IsRunning = false
		a.mu.Unlock()
	}(devicePath)

	return nil
}

// StopCamera interrompe a transmissão de um dispositivo específico
func (a *App) StopCamera(devicePath string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	session, exists := a.sessions[devicePath]

	// Guard Clauses (No Else)
	if !exists {
		return nil
	}
	if !session.IsRunning {
		return nil
	}

	// Aciona o ctx.Done() na Goroutine, parando o FFmpeg instantaneamente
	session.Cancel()
	session.IsRunning = false

	return nil
}
