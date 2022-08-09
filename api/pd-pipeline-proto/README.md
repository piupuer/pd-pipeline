# Overview
Proto files for paddle pipeline

exec it in your project root:
```shell
mkdir -p api/build/pd-pipeline
/usr/local/bin/protoc \
    --proto_path=api/pd-pipeline-proto \
    --go_out=api/build/pd-pipeline --go_opt=paths=source_relative \
    --go-grpc_out=api/build/pd-pipeline --go-grpc_opt=paths=source_relative \
    api/pd-pipeline-proto/**.proto
```