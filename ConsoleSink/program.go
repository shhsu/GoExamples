package console_sink

import "os/exec"

// Example usage, not tested AT ALL
func main() {
	sink := NewConsoleSink()
	defer sink.Close()
	cmd := exec.Command("bla", "bla", "bla")
	sink.AttachCommand(cmd)
	// go run cmd
	// go display sink content
}
