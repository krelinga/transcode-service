package activity

import (
    "context"
)

type UpdateNfoRequest struct {
    Path string
}

type UpdateNfoReply struct {
}

func UpdateNfo(ctx context.Context, req *UpdateNfoRequest) (*UpdateNfoReply, error) {
    // TODO: implement a wrapper around updating .nfo files.
    return nil, nil
}
