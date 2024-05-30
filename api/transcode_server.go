package main

import (
    "context"

    "github.com/krelinga/transcode-service/pb"
    "go.temporal.io/sdk/client"
)

type TranscodeServer struct {
    pb.UnimplementedTranscodeServer
    temporalC client.Client
}

func NewTranscodeServer(temporalC client.Client) *TranscodeServer {
    return &TranscodeServer{
        temporalC: temporalC,
    }
}

func (_ *TranscodeServer) BeginOneFile(ctx context.Context, req *pb.BeginOneFileRequest) (*pb.BeginOneFileReply, error) {
    return &pb.BeginOneFileReply{}, nil
}

func (_ *TranscodeServer) CheckOneFile(ctx context.Context, req *pb.CheckOneFileRequest) (*pb.CheckOneFileReply, error) {
    return &pb.CheckOneFileReply{}, nil
}
