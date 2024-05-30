package main

import (
    "fmt"
    "log"
    "net"

    "github.com/krelinga/transcode-service/pb"
    "go.temporal.io/sdk/client"
    "google.golang.org/grpc"
)

func mainOrError() error {
    temporalC, err := client.Dial(client.Options{})
    if err != nil {
        return fmt.Errorf("Could not create temporal client: %w", err)
    }
    defer temporalC.Close()

    lis, err := net.Listen("tcp", ":25003")
    if err != nil {
        return err
    }
    grpcServer := grpc.NewServer()
    pb.RegisterTranscodeServer(grpcServer, &TranscodeServer{})
    grpcServer.Serve(lis)  // Runs as long as the server is alive.

    return nil
}

func main() {
    if err := mainOrError(); err != nil {
        log.Fatal(err)
    }
}
