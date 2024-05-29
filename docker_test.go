package main

import (
    "bytes"
    "fmt"
    "os"
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
    DockerFile string
    ContainerId string
}

func newTestContainer(name, dockerFile string) *testContainer {
    return &testContainer{
        DockerFile: dockerFile,
        ContainerId: fmt.Sprintf("tcservice-docker-test-%s-%s", name, uuid.NewString()),
    }
}

func (tc *testContainer) BuildImage(t *testing.T) {
    t.Helper()
    args := []string{
        "image", "build",
        "-f", tc.DockerFile,
        "-t", tc.ContainerId,
        ".",
    }
    cmd := exec.Command("docker", args...)
    cmdOutput := captureOutput(cmd)
    if err := cmd.Run(); err != nil {
        const template = "Could not build docker image %s.  Error was %s.  Output:\n%s"
        t.Fatalf(template, tc.ContainerId, err, cmdOutput)
    }
}

func (tc *testContainer) DeleteImage(t *testing.T) {
    t.Helper()
    cmd := exec.Command("docker", "image", "rm", tc.ContainerId)
    cmdOutput := captureOutput(cmd)
    if err := cmd.Run(); err != nil {
        const template = "Could not delete docker image %s.  Error was %s.  Output:\n%s"
        t.Fatalf(template, tc.ContainerId, err, cmdOutput)
    }
}

type testProject string

func newTestProject() testProject {
    return testProject(fmt.Sprintf("tcservice-%s", uuid.NewString()))
}

func (tp testProject) Up(t *testing.T, api, worker *testContainer) {
    t.Helper()
    args := []string {
        "compose", "-p", string(tp),
        "up", "-d",
    }
    cmd := exec.Command("docker", args...)
    apiEnv := fmt.Sprintf("API_IMAGE=%s", api.ContainerId)
    workerEnv := fmt.Sprintf("WORKER_IMAGE=%s", worker.ContainerId)
    cmd.Env = append(os.Environ(), apiEnv, workerEnv)
    cmdOutput := captureOutput(cmd)
    if err := cmd.Run(); err != nil {
        const template = "Could not up containers.  Error was %s.  Output:\n%s"
        t.Fatalf(template, err, cmdOutput)
    }
}

func (tp testProject) Down(t *testing.T) {
    t.Helper()
    args := []string {
        "compose", "-p", string(tp),
        "down",
    }
    cmd := exec.Command("docker", args...)
    cmdOutput := captureOutput(cmd)
    if err := cmd.Run(); err != nil {
        const template = "Could not down containers.  Error was %s.  Output:\n%s"
        t.Fatalf(template, err, cmdOutput)
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
    tp := newTestProject()
    tp.Up(t, apiTc, workerTc)
    defer tp.Down(t)
}
