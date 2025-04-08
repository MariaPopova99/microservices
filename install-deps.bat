@echo off
setlocal

set GOBIN=C:\Users\hallo\GoProjects\micro\microservices\bin

echo Installing protoc-gen-go...
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

echo Installing protoc-gen-go-grpc...
go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

endlocal
