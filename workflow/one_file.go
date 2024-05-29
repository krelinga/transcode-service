package workflow

import (
    "fmt"
    "path/filepath"

    "github.com/krelinga/transcode-service/activity"
    "go.temporal.io/sdk/workflow"
)

type OneFileRequest struct {
    InPath string
    OutPath string
}

type OneFileReply struct {
}

func OneFile(ctx workflow.Context, req *OneFileRequest) (*OneFileReply, error) {
    // Create output directory if necessary.
    outDir := filepath.Dir(req.OutPath)
    err := workflow.ExecuteLocalActivity(ctx, activity.Mkdir, &activity.MkdirRequest{
        Path: outDir,
    }).Get(ctx, nil)
    if err != nil {
        return nil, err
    }

    // Transcode the .mkv file.
    var tcReply activity.HBTranscodeReply
    err = workflow.ExecuteActivity(ctx, activity.HBTranscode, &activity.HBTranscodeRequest{
        InPath: req.InPath,
        OutPath: req.OutPath,
    }).Get(ctx, &tcReply)
    if err != nil {
        return nil, err
    }

    // Discover all files with the input path prefix.
    inBase := filepath.Base(req.InPath)
    var matchReply activity.MatchReply
    err = workflow.ExecuteLocalActivity(ctx, activity.Match, &activity.MatchRequest{
        PathBase: inBase,
    }).Get(ctx, &matchReply)

    // Create a mapping from input path to output path for all non-mkv files.
    type pathPair struct {
        inPath string
        outPath string
    }
    pathMap := []*pathPair{}
    outBase := filepath.Base(req.OutPath)
    for _, p := range matchReply.Paths {
        if p == req.InPath {
            continue
        }
        pathMap = append(pathMap, &pathPair{
            inPath: p,
            outPath: fmt.Sprintf("%s.%s", outBase, filepath.Ext(p)),

        })
    }

    // Copy all non-MKV files to the corresponding output location.
    wg := workflow.NewWaitGroup(ctx)
    copyErrs := workflow.NewBufferedChannel(ctx, len(pathMap))
    for _, p := range pathMap {
        p := p

        wg.Add(1)
        workflow.Go(ctx, func(ctx workflow.Context) {
            err := workflow.ExecuteActivity(ctx, activity.Copy, &activity.CopyRequest{
                InPath: p.inPath,
                OutPath: p.outPath,
            }).Get(ctx, nil)
            if err != nil {
                copyErrs.Send(ctx, err)
            }
            wg.Done()
        })
    }
    wg.Wait(ctx)
    copyErrs.Close()
    if copyErrs.Len() > 0 {
        var err error
        copyErrs.Receive(ctx, &err)
        return nil, err
    }

    // Find the `.nfo` file, if it exists.
    var nfoPathPair *pathPair
    for _, p := range pathMap {
        if filepath.Ext(p.inPath) == ".nfo" {
            nfoPathPair = p
            break
        }
    }
    if nfoPathPair == nil {
        return &OneFileReply{}, nil
    }

    // Update .nfo file, if it exists.
    var nfoReply activity.UpdateNfoReply
    err = workflow.ExecuteActivity(ctx, activity.UpdateNfo, activity.UpdateNfoRequest{
        Path: nfoPathPair.outPath,
    }).Get(ctx, &nfoReply)
    if err != nil {
        return nil, err
    }

    return &OneFileReply{}, nil
}
