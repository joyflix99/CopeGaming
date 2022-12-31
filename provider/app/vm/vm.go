package vm

import (
	"bufio"
	"log"
	"os/exec"
	"strconv"
)

func StartVM(id string, appName string, videoRelayPort, audioRelayPort, syncPort int) error {
	log.Printf("[%s] Spinning off VM\n", id)

	params := []string{
		id,
		strconv.Itoa(videoRelayPort),
		strconv.Itoa(audioRelayPort),
		strconv.Itoa(syncPort),
		appName,
	}

	for _, value := range params {
		log.Printf("[%s] params\n", value)
	}
	cmd := exec.Command("./startVM.sh", params...)

	stderr, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)

	go func() {
		for scanner.Scan() {
			m := scanner.Text()
			log.Printf("[%s]", m)
		}
	}()

	return nil
}

func StopVM(id, appName string) error {
	log.Printf("[%s] Stopping VM\n", id)

	params := []string{
		id,
		appName,
	}
	cmd := exec.Command("./stopVM.sh", params...)
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}
