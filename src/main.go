package main

import (
	"fmt"

	"kv"
	"search"
	"update"
)

func dbInit() {
	levelDbName := "TagInfo"
	kv.Open(&levelDbName)
}

func main() {
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
