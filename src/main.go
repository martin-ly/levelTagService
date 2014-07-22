package main

import (
	"errors"
	"fmt"
	"os"

	"git.apache.org/thrift.git/lib/go/thrift"
	"kv"
	"search"
	"thrift/gen-go-modified/tagSearchService"
	"update"
)

func dbInit() {
	levelDbName := "TagInfo"
	kv.Open(&levelDbName)
}

func _main_example() {
	dbInit()
	tid := int32(12345)
	update.Insert(tid, 5999990, 123)
	update.Insert(tid, 7121999, 887)
	update.Insert(tid, 6818989, 909)
	update.Insert(tid, 4800000, 1299)
	update.Insert(tid, 5800998, 1999)
	update.Insert(tid, 7234568, 299)

	start, end := search.FuzzyRange(tid)
	fmt.Printf("FuzzyRange (%d, %d)\n", start, end)

	start, end = search.Range(tid)
	fmt.Printf("Range (%d, %d)\n", start, end)

	uList := search.GetUsrByRange(tid, start, end, 1000)

	for i, u := range uList {
		fmt.Println("usr ", i, ": uid: ", u.Uid, ", weight: ", u.Score)
	}
}

const (
	NetworkAddr = "127.0.0.1:9912"
)

type TagSearchServiceImpl struct{}

func (this *TagSearchServiceImpl) GetRange(tagId int32) (r *tagSearchService.Range, err error) {
	start, end := search.Range(tagId)
	if start > end {
		end = start
	}
	r = new(tagSearchService.Range)
	r.StartUid = int32(start)
	r.EndUid = int32(end)
	return
}

func (this *TagSearchServiceImpl) GetUsrs(tagId int32, r *tagSearchService.Range, limitSize int32) (usrInfo []tagSearchService.UsrInfo, err error) {
	if limitSize <= 0 || limitSize >= 50000 {
		err = errors.New("illegal limitSize, please range limitSize in (0, 50000)")
		return
	}
	if r.StartUid >= r.EndUid {
		err = errors.New("range error, startUid should less than endUid")
		return
	}
	uList := search.GetUsrByRange(tagId, int(r.StartUid), int(r.EndUid), int(limitSize))
	usrInfo = make([]tagSearchService.UsrInfo, 0, len(uList))
	for _, x := range uList {
		var info tagSearchService.UsrInfo
		info.Uid = x.Uid
		info.Weight = x.Score
		usrInfo = append(usrInfo, info)
	}
	return
}

func main() {

	dbInit()

	//transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	transportFactory := thrift.NewTBufferedTransportFactory(4096 * 4) // 4 PageSize
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	//protocolFactory := thrift.NewTCompactProtocolFactory()

	serverTransport, err := thrift.NewTServerSocket(NetworkAddr)
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	handler := &TagSearchServiceImpl{}
	processor := tagSearchService.NewTagSearchServiceProcessor(handler)

	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("thrift server in", NetworkAddr)
	server.Serve()
}
