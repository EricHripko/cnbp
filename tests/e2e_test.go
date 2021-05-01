package e2e

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func build(t *testing.T, project string) {
	// Arrange
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd := exec.Command(
		"docker",
		"build",
		"-t",
		project,
		"-f", "project.toml",
		".",
	)
	cmd.Dir = filepath.Join(".", project)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	cmd.Env = append(os.Environ(), "DOCKER_BUILDKIT=1")

	// Act
	err := cmd.Run()

	// Assert
	if err != nil {
		t.Log("stdout: ", stdout.String())
		t.Log("stderr: ", stderr.String())
		t.Log(err)
		t.Fail()
	}
}

func TestJavaMaven(t *testing.T) {
	// Arrange
	project := "java-maven"
	build(t, project)

	// Act
	cmd := exec.Command("docker", "run", "--rm", "-p", "8080:8080", project, "/cnb/lifecycle/launcher")
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	for i := 0; i <= 10; i++ {
		res, err := http.Get("http://localhost:8080")
		if err == nil && res.StatusCode == 200 {
			return
		}

		time.Sleep(time.Second)
	}
	t.Fatal("Cannot reach the endpoint")
}
