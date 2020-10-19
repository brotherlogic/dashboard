protoc --proto_path ../../../ -I=./proto --go_out=plugins=grpc:./proto proto/dashboard.proto
mv proto/github.com/brotherlogic/dashboard/proto/* ./proto
