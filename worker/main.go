package main

import (
    "log"

    "github.com/krelinga/transcode-service/activity"
    "github.com/krelinga/transcode-service/common"
    "github.com/krelinga/transcode-service/workflow"
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
)

func main() {
    c, err := client.Dial(client.Options{})
    if err != nil {
        log.Fatalln("unable to create Temporal client", err)
    }
    defer c.Close()

    // This worker hosts both Workflow and Activity functions
    w := worker.New(c, common.TaskQueue, worker.Options{})

    // Workflows
    w.RegisterWorkflow(workflow.TranscodeOneFile)

    // Activities
    w.RegisterActivity(activity.Copy)
    w.RegisterActivity(activity.HBTranscode)
    w.RegisterActivity(activity.Match)
    w.RegisterActivity(activity.Mkdir)
    w.RegisterActivity(activity.UpdateNfo)

    // Start listening to the Task Queue
    err = w.Run(worker.InterruptCh())
    if err != nil {
        log.Fatalln("unable to start Worker", err)
    }
}
