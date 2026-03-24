package infra

import (
	"fmt"
	"os/exec"
	"strings"
)

// ConfigureVirtualCameras cria a string de configuração do Kernel dinamicamente e executa via pkexec
func ConfigureVirtualCameras(count int) error {
	if count < 1 || count > 10 {
		return fmt.Errorf("número de câmeras deve ser entre 1 e 10")
	}

	var devices []string
	var labels []string
	var caps []string
	var fuserTargets []string

	for i := 0; i < count; i++ {
		videoNr := 10 + i
		devices = append(devices, fmt.Sprintf("%d", videoNr))
		labels = append(labels, fmt.Sprintf("LoopCam-%d", i+1))
		caps = append(caps, "1")
		fuserTargets = append(fuserTargets, fmt.Sprintf("/dev/video%d", videoNr))
	}

	targetsStr := strings.Join(fuserTargets, " ")

	// Adicionamos um loop de estabilização (Polling) dentro do bash.
	// O root aguarda a criação assíncrona do device e libera o acesso para o User-Space.
	script := fmt.Sprintf(`
		fuser -k %s || true
		rmmod v4l2loopback || true
		modprobe v4l2loopback devices=%d video_nr=%s card_label="%s" exclusive_caps=%s
		
		# UDEV Race Condition Fix: Aguarda os nós /dev serem criados fisicamente
		for dev in %s; do
			for i in {1..20}; do
				if [ -c "$dev" ]; then
					chmod 666 "$dev"
					break
				fi
				sleep 0.1
			done
		done
	`,
		targetsStr,
		count,
		strings.Join(devices, ","),
		strings.Join(labels, ","),
		strings.Join(caps, ","),
		targetsStr,
	)

	cmd := exec.Command("pkexec", "bash", "-c", script)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao configurar kernel: %v - %s", err, string(output))
	}

	return nil
}
