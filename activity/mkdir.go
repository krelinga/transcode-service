package activity

import (
    "context"
    "os"
)

type MkdirRequest struct {
    Path string
}

func Mkdir(ctx context.Context, req *MkdirRequest) error {
    return os.MkdirAll(req.Path, 0755)
}
