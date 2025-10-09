
FROM golang:1.22 AS builder
WORKDIR /workspace
COPY go.mod .
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/extender ./cmd/extender


FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /
COPY --from=builder /out/extender /extender
USER nonroot:nonroot
EXPOSE 8000
ENTRYPOINT ["/extender"]
