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
go build -o bin/levelTagService src/main.go
