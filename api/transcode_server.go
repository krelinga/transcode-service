package main

import (
    "context"
    "fmt"

    "github.com/krelinga/transcode-service/common"
    "github.com/krelinga/transcode-service/pb"
    "github.com/krelinga/transcode-service/workflow"
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

func tcOneFileWorkflowID(inPath, outPath string) string {
    return fmt.Sprintf("%s to %s", inPath, outPath)
}

func (tcs *TranscodeServer) BeginOneFile(ctx context.Context, req *pb.BeginOneFileRequest) (*pb.BeginOneFileReply, error) {
    options := client.StartWorkflowOptions{
        ID: tcOneFileWorkflowID(req.InPath, req.OutPath),
        TaskQueue: common.TaskQueue,
    }
    workflowReq := &workflow.TranscodeOneFileRequest{
        InPath: req.InPath,
        OutPath: req.OutPath,
    }
    workflowRun, err := tcs.temporalC.ExecuteWorkflow(ctx, options, workflow.TranscodeOneFile, workflowReq)
    if err != nil {
        return nil, err
    }
    return &pb.BeginOneFileReply{
        Key: &pb.OneFileKey{
            InPath: req.InPath,
            OutPath: req.OutPath,
            Instance: workflowRun.GetRunID(),
        },
    }, nil
}

func (_ *TranscodeServer) CheckOneFile(ctx context.Context, req *pb.CheckOneFileRequest) (*pb.CheckOneFileReply, error) {
    return &pb.CheckOneFileReply{}, nil
}
