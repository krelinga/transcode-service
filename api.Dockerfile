FROM golang:1.21 AS build_stage

WORKDIR /app
COPY go.mod go.sum ./
COPY pb/*.go ./pb/
COPY common/*.go ./common/
COPY api/*.go ./api/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o api-server ./api

FROM gcr.io/distroless/static-debian12 AS build-release-stage
WORKDIR /
COPY --from=build_stage /app/api-server /api-server

EXPOSE 25003

ENTRYPOINT ["/api-server"]
