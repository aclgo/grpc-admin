push:
    sudo docker compose up --build -d
migrate_up:migrate -database postgres://grpc-admin:grpc-admin@localhost:5432/grpc-admin?sslmode=disable -path migrations up 1
    
migrate_down:
    migrate -database postgres://grpc-jwt:grpc-jwt@localhost:5432/grpc-jwt?sslmode=disable -path migrations down 1

protoc:
    protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    user.proto