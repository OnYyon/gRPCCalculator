export PATH="$PATH:$(go env GOPATH)/bin" 

protoc  -I ./proto/google \
        -I ./proto \
        --go_out=./proto/gen/ \
        --go_opt=paths=source_relative \
        --go-grpc_out=./proto/gen/ \
        --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=./proto/gen \
        --grpc-gateway_opt paths=source_relative \
        --grpc-gateway_opt logtostderr=true \
        ./proto/orchestrator.proto
