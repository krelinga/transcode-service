package main

import (
    "fmt"
    "log"
    "net"
    "os"

    "github.com/krelinga/transcode-service/pb"
    "go.temporal.io/sdk/client"
    "google.golang.org/grpc"
)

func temporalHostPort() (string, error) {
    get := func(k string) (string, error) {
        v := os.Getenv(k)
        if len(v) == 0 {
            return "", fmt.Errorf("could not read env var %s", k)
        }
        return v, nil
    }

    host, err := get("temporal_host")
    if err != nil {
        return "", err
    }
    port, err := get("temporal_port")
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("%s:%s", host, port), nil
}

func mainOrError() error {
    hp, err := temporalHostPort()
    if err != nil {
        return err
    }
    temporalC, err := client.Dial(client.Options{
        HostPort: hp,
    })
    if err != nil {
        return fmt.Errorf("Could not create temporal client: %w", err)
    }
    defer temporalC.Close()

    lis, err := net.Listen("tcp", ":25003")
    if err != nil {
        return err
    }
    grpcServer := grpc.NewServer()
    pb.RegisterTranscodeServer(grpcServer, NewTranscodeServer(temporalC))
    grpcServer.Serve(lis)  // Runs as long as the server is alive.

    return nil
}

func main() {
    if err := mainOrError(); err != nil {
        log.Fatal(err)
    }
}
