package activity

import (
    "context"
)

type CopyRequest struct {
    InPath string
    OutPath string
}

func Copy(ctx context.Context, req *CopyRequest) error {
    // TODO: implement copy of a single file.
    return nil
}
