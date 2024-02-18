FROM golang:1.21-bullseye AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api/main.go 

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app/api /api

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/api"]
