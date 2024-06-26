package main

import (
    "bytes"
    "context"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "testing"

    "github.com/google/uuid"
    "github.com/krelinga/transcode-service/pb"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
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

type testProject struct {
    name string
    dir string

    // If set, these will be passed with `-f option to docker compose command.
    ComposeFiles []string
}

func newTestProject(dir, namePrefix string) *testProject {
    return &testProject{
        name: fmt.Sprintf("%s-%s", namePrefix, uuid.NewString()),
        dir: dir,
    }
}

func (tp *testProject) Up(t *testing.T, envEdits... string) {
    t.Helper()
    args := []string {"compose", "-p", tp.name}
    for _, f := range tp.ComposeFiles {
        args = append(args, "-f", f)
    }
    args = append(args, "up", "-d")
    cmd := exec.Command("docker", args...)
    cmd.Dir = tp.dir
    cmd.Env = append(os.Environ(), envEdits...)
    cmdOutput := captureOutput(cmd)
    if err := cmd.Run(); err != nil {
        const template = "Could not up containers.  Error was %s.  Output:\n%s"
        t.Fatalf(template, err, cmdOutput)
    }
}

func (tp *testProject) Down(t *testing.T) {
    t.Helper()
    args := []string {
        "compose", "-p", tp.name,
        "down",
    }
    cmd := exec.Command("docker", args...)
    cmdOutput := captureOutput(cmd)
    if err := cmd.Run(); err != nil {
        const template = "Could not down containers.  Error was %s.  Output:\n%s"
        t.Fatalf(template, err, cmdOutput)
    }
}

func makeTempDir(t *testing.T) string {
    dir, err := os.MkdirTemp("", "")
    if err != nil {
        t.Fatalf("Could not make temp directory. Error was %s", err)
    }
    return dir
}

func deleteTempDir(t *testing.T, dir string) {
    if err := os.RemoveAll(dir); err != nil {
        t.Fatalf("Could not delete temp directory %s.  Error was %s", dir, err)
    }
}

func cloneTemporalGitRepo(t *testing.T, dir string) {
    cmd := exec.Command("git", "clone", "https://github.com/temporalio/docker-compose.git")
    cmd.Dir = dir
    cmdOutput := captureOutput(cmd)
    if err := cmd.Run(); err != nil {
        t.Fatalf("Could not clone temporal git repo.  Error was %s.  Output:\n%s", err, cmdOutput)
    }
}

func getWorkingDir(t *testing.T) string {
    dir, err := os.Getwd()
    if err != nil {
        t.Fatalf("Could not read working directory.  Error was %s", err)
    }
    return dir
}

func createApiStub(t *testing.T) (pb.TranscodeClient, func()) {
    const target = "docker-daemon:25003"
    creds := grpc.WithTransportCredentials(insecure.NewCredentials())
    conn, err := grpc.DialContext(context.Background(), target, creds)
    if err != nil {
        t.Fatalf("Could not dial API server: %e", err)
    }
    return pb.NewTranscodeClient(conn), func() {
        conn.Close()
    }
}

func TestDocker(t *testing.T) {
    t.Parallel()
    tmpDir := makeTempDir(t)
    defer deleteTempDir(t, tmpDir)
    cloneTemporalGitRepo(t, tmpDir)
    temporal := newTestProject(filepath.Join(tmpDir, "docker-compose"), "temporal")
    temporal.Up(t)
    defer temporal.Down(t)
    apiTc := newTestContainer("api", "api.Dockerfile")
    apiTc.BuildImage(t)
    defer apiTc.DeleteImage(t)
    workerTc := newTestContainer("worker", "worker.Dockerfile")
    workerTc.BuildImage(t)
    defer workerTc.DeleteImage(t)
    tp := newTestProject(getWorkingDir(t), "tcservice")
    tp.ComposeFiles = []string{"compose.yaml", "compose-for-test.yaml"}
    apiEnv := fmt.Sprintf("API_IMAGE=%s", apiTc.ContainerId)
    workerEnv := fmt.Sprintf("WORKER_IMAGE=%s", workerTc.ContainerId)
    workerUidEnv := fmt.Sprintf("WORKER_UID=%d", os.Getuid())
    workerGidEnv := fmt.Sprintf("WORKER_GID=%d", os.Getgid())
    temporalHostEnv := "temporal_host=docker-daemon"
    tp.Up(t, apiEnv, workerEnv, workerUidEnv, workerGidEnv, temporalHostEnv)
    defer tp.Down(t)
    stub, stubCleanup := createApiStub(t)
    defer stubCleanup()

    req := &pb.BeginOneFileRequest{
        InPath: "/testdata/sample_640x360.mkv",
        OutPath: "/testdata/out.mkv",
    }
    _, err := stub.BeginOneFile(context.Background(), req)
    if err != nil {
        t.Fatalf("Could not call BeginOneFile(): %s", err)
    }
}
