package main

import (
    "bytes"
    "fmt"
    "os/exec"
    "sync"
    "testing"

    "github.com/google/uuid"
)

func captureOutput(cmd *exec.Cmd) *bytes.Buffer {
    cmdOutput := &bytes.Buffer{}
    cmd.Stdout = cmdOutput
    cmd.Stderr = cmdOutput
    return cmdOutput
}

type testContainer struct {
    dockerFile string
    containerId string
}

func newTestContainer(name, dockerFile string) testContainer {
    return testContainer{
        dockerFile: dockerFile,
        containerId: fmt.Sprintf("tcservice-docker-test-%s-%s", name, uuid.NewString()),
    }
}

func (tc *testContainer) BuildImage(t *testing.T) {
    t.Helper()
    args := []string{
        "image", "build",
        "-f", tc.dockerFile,
        "-t", tc.containerId,
        ".",
    }
    cmd := exec.Command("docker", args...)
    cmdOutput := captureOutput(cmd)
    if err := cmd.Run(); err != nil {
        const template = "Could not build docker image %s.  Error was %s.  Output:\n%s"
        t.Fatalf(template, tc.containerId, err, cmdOutput)
    }
}

func (tc *testContainer) DeleteImage(t *testing.T) {
    t.Helper()
    cmd := exec.Command("docker", "image", "rm", tc.containerId)
    cmdOutput := captureOutput(cmd)
    if err := cmd.Run(); err != nil {
        const template = "Could not delete docker image %s.  Error was %s.  Output:\n%s"
        t.Fatalf(template, tc.containerId, err, cmdOutput)
    }
}

func TestDocker(t *testing.T) {
    t.Parallel()
    apiTc := newTestContainer("api", "api.Dockerfile")
    workerTc := newTestContainer("worker", "worker.Dockerfile")
    func() {
        // build containers in parallel
        wg := &sync.WaitGroup{}
        wg.Add(2)
        go func() {
            apiTc.BuildImage(t)
            wg.Done()
        }()
        go func() {
            workerTc.BuildImage(t)
            wg.Done()
        }()
        wg.Wait()
    }()
    defer apiTc.DeleteImage(t)
    defer workerTc.DeleteImage(t)
}
