package activity

import (
    "context"
)

type MatchRequest struct {
    PathBase string
}

type MatchReply struct {
    Paths []string
}

func Match(ctx context.Context, req *MatchRequest) (*MatchReply, error) {
    // TODO: implement a wrapper around matching paths with a given prefix.
    return nil, nil
}
