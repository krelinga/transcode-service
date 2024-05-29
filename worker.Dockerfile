FROM golang:1.21 AS build_stage

WORKDIR /app
COPY go.mod go.sum ./
COPY activity/*.go ./activity/
COPY common/*.go ./common/
COPY pb/*.go ./pb/
COPY worker/*.go ./worker/
COPY workflow/*.go ./workflow/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o worker-server ./worker

FROM gcr.io/distroless/static-debian12 AS build-release-stage
WORKDIR /
COPY --from=build_stage /app/worker-server /worker-server

ENTRYPOINT ["/worker-server"]
