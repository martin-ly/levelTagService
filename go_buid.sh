cur_path=`pwd`
export GOPATH=$cur_path:$GOPATH
export CGO_CFLAGS="-I${cur_path}/lib/leveldb/include"
export CGO_LDFLAGS="-L${cur_path}/lib/leveldb"

#cd lib/leveldb/
#make
#cd ../../

cd src/
go test kv
cd ..
#thrift -o src/thrift --gen go src/thrift/tagSearchService.thrift
go build -o bin/levelTagService src/main.go
go build -o bin/client src/client/tag_search_service-remote.go
