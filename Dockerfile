FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o grpc-admin ./cmd/grpc-admin/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/grpc-admin ./

EXPOSE 50052

ENTRYPOINT [ "./grpc-admin" ]