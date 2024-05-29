package main

import (
    "log"
    "net"

    "google.golang.org/grpc"
    "github.com/krelinga/transcode-service/pb"
)

func mainOrError() error {
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
