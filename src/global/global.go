package global

import (
	"strconv"
)

var L1 string = "@"
var L2 string = "#"
var L3 string = "$"

var Base int = 4                            /* start: 2 byte, end: 2 bytes*/
var ItemNum int = 1000                      /* number of items in each level */
var IndexBufSize int = (ItemNum / 8) + Base /* 1000 user */
var DataBufSize int = (ItemNum * 4) + Base  /* 1000 user */

var L3Space int = (DataBufSize - Base) / 4  /* 第三层可用空间（每个用户权重为4Bytes） */
var L2Space int = (IndexBufSize - Base) * 8 /* 第二层可用空间（每个空间为1bit）*/

var L2BitUsrs int = L3Space
var L1BitUsrs int = L3Space * L2Space

func L1Key(id int32) string {
	return L1 + strconv.Itoa(int(id))
}

func L2Key(l1key string, bitoff int) string {
	if bitoff >= (IndexBufSize<<3) || bitoff < 0 {
		panic("bit off error: " + strconv.Itoa(bitoff))
	}
	return l1key + L2 + strconv.FormatInt(int64(bitoff), 16)
}

func L3Key(l2key string, bitoff int) string {
	if bitoff >= (IndexBufSize<<3) || bitoff < 0 {
		panic("bit off error: " + strconv.Itoa(bitoff))
	}
	return l2key + L3 + strconv.FormatInt(int64(bitoff), 16)
}
