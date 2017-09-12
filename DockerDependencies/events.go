package dependencies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type DockerImageEventsListener struct {
	cmd *exec.Cmd
	buf *bytes.Buffer
}

type DockerDependencies struct {
	Image             string   `json:"image"`
	RuntimeDependency string   `json:"runtime-dependency"`
	BuildDependencies []string `json:"build-dependencies,omitempty"`
}

type dockerImageLogEntry struct {
	Action string `json:"action"`
	Name   string `json:"name"`
	ID     string `json:"id"`
}

func NewDockerImageEventsListener() (*DockerImageEventsListener, error) {
	l := &DockerImageEventsListener{}
	formatString := "{\"action\":\"{{.Status}}\",\"id\":\"{{.ID}}\",\"name\":\"{{index .Actor.Attributes \"name\"}}\"}"
	l.cmd = exec.Command("docker", "events", "--filter", "type=image", "--format", formatString)
	l.buf = &bytes.Buffer{}
	l.cmd.Stdout = l.buf
	err := l.cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("Failed to start Docker event listener: %s", err)
	}
	return l, nil
}

func (l *DockerImageEventsListener) ResolveDependencies() ([]DockerDependencies, error) {
	outputs, err := l.parseOutputs()
	if err != nil {
		return nil, err
	}
	var imagesPulled []string
	results := []DockerDependencies{}
	for _, entry := range outputs {
		if entry.Action == "pull" {
			imagesPulled = append(imagesPulled, entry.ID)
		} else if entry.Action == "tag" {
			image := entry.Name
			if len(imagesPulled) == 0 {
				logrus.Errorf("Failed to determine dependency of %s, this image's dependency will be ignored", image)
			} else {
				cutoff := len(imagesPulled) - 1
				result := DockerDependencies{Image: image, RuntimeDependency: imagesPulled[cutoff]}
				if cutoff > 0 {
					result.BuildDependencies = imagesPulled[0:cutoff]
				}
				results = append(results, result)
			}
			imagesPulled = []string{}
		} else {
			logrus.Info("Irrelevant docker image action: %s ID: %s, name: %s was performed, ignoring", entry.Action, entry.ID, entry.Name)
		}
	}
	return results, nil
}

func (l *DockerImageEventsListener) parseOutputs() ([]dockerImageLogEntry, error) {
	err := l.cmd.Process.Kill()
	if err != nil {
		releaseErr := l.cmd.Process.Release()
		if releaseErr != nil {
			logrus.Errorf("Failed to release resources for Docker event listener process, error: %s", releaseErr)
		}
		pid := l.cmd.Process.Pid
		return nil, fmt.Errorf("Failed to kill Docker event listener process, pid: %d, error: %s", pid, err)
	}
	err = l.cmd.Wait()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok || exitErr.Error() != "exit status 1" {
			logrus.Warnf("Docker event listener ended with error: %s", exitErr)
		}
	}
	lines := split(l.buf.Bytes(), '\n')
	results := make([]dockerImageLogEntry, len(lines), len(lines))
	for i, line := range lines {
		var entry dockerImageLogEntry
		err := json.Unmarshal(line, &entry)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse docker log, line: %s, error: %s", string(line), err)
		}
		results[i] = entry
	}

	return results, nil
}

func split(s []byte, sep byte) [][]byte {
	result := [][]byte{}
	head := 0
	tail := 0
	for tail < len(s) {
		if s[tail] == sep {
			if tail > head+1 {
				result = append(result, s[head:tail])
			}
			head = tail + 1
		}
		tail++
	}
	return result
}
