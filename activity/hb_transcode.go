package activity

import (
    "context"
)

type HBTranscodeRequest struct {
    InPath string
    OutPath string
}

type HBTranscodeReply struct {
}

func HBTranscode(ctx context.Context, req *HBTranscodeRequest) (*HBTranscodeReply, error) {
    // TODO: implement a wrapper around transcoding w/ Handbrake.
    return nil, nil
}
