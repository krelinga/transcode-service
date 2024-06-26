FROM golang:1.21 AS build_stage

WORKDIR /app
COPY go.mod go.sum ./
COPY pb/*.go ./pb/
COPY common/*.go ./common/
COPY workflow/*.go ./workflow/
COPY activity/*.go ./activity/
COPY api/*.go ./api/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api-server ./api

FROM gcr.io/distroless/static-debian12 AS build-release-stage
WORKDIR /
COPY --from=build_stage /app/api-server /api-server

EXPOSE 25003

ENV temporal_host=localhost
ENV temporal_port=7233

ENTRYPOINT ["/api-server"]
