FROM golang:1.26 AS builder
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/jungle-gaming ./cmd/api

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /
COPY --from=builder /out/jungle-gaming /jungle-gaming
ENTRYPOINT ["/jungle-gaming"]
