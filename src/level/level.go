package level

import (
	"global"
	"kv"

	"strconv"
)

type LevelInterface interface {
	Assign(string, int, string)
}

type Level struct { // 1 page
	start   int
	end     int
	level   int // 1, 2, 3
	key     string
	bufSize int
	buf     []byte
}

func (this *Level) SetStart(byteOff int) {
	this.buf[0] = byte((byteOff >> 8) & 0xFF)
	this.buf[1] = byte(byteOff & 0xFF)
	this.start = byteOff
}

func (this *Level) SetEnd(byteOff int) {
	this.buf[2] = byte(byteOff >> 8)
	this.buf[3] = byte(byteOff)
	this.end = byteOff
}

func (this *Level) SetLevels(levels int) {
	this.level = levels
}

func (this *Level) SetKey(key string) {
	this.key = key
}

func (this *Level) Assign(key string, levels int, val string) {
	copy(this.buf[:], []byte(val))
	this.start = (int(this.buf[0]) << 8) + int(this.buf[1])
	this.end = (int(this.buf[2]) << 8) + int(this.buf[3])
	this.SetKey(key)
	this.SetLevels(levels)
}

func (this *Level) Start() int {
	return this.start
}

func (this *Level) End() int {
	return this.end
}

/* Range[0, global.xxxBufSize) */
func (this *Level) Range() (start int, end int) {
	start = this.start
	end = this.end
	return
}

type IndexLevel struct {
	Level
}

func NewIndexLevel() *IndexLevel {
	l := new(IndexLevel)
	l.bufSize = global.IndexBufSize
	l.buf = make([]byte, l.bufSize)
	l.SetStart(global.IndexBufSize - 1)
	l.SetEnd(0)
	return l
}

func (this *IndexLevel) SetBit(bit int) bool {
	if this.CheckBitFlat(bit) {
		return false
	}
	byteOff := bit / 8
	bitOff := bit % 8
	mask := byte(0x80) >> uint(bitOff)
	this.buf[global.Base+byteOff] |= mask
	start, end := this.Range()
	if byteOff < start {
		this.SetStart(byteOff)
	}
	if byteOff >= end {
		this.SetEnd(byteOff + 1)
	}
	return true
}

func (this *IndexLevel) CheckByte(idx int) bool {
	/* byte[0~3] is used to store Range */
	if idx < 0 || idx >= (global.IndexBufSize-global.Base) {
		panic("CheckByte idx error: " + strconv.Itoa(idx))
	}
	if this.buf[idx+global.Base]&0xFF == 0 {
		return false
	} else {
		return true
	}
}

func (this *IndexLevel) CheckBitFlat(bit int) bool {
	return this.CheckBit(bit/8, bit%8)
}

func (this *IndexLevel) CheckBit(idx int, bit int) bool {
	/* byte[0~3] is used to store Range */
	if idx < 0 || idx >= ((global.IndexBufSize-global.Base)<<3) {
		panic("CheckBit idx error: " + strconv.Itoa(idx))
	}
	if bit < 0 || bit >= 8 {
		panic("CheckBit bit error: " + strconv.Itoa(bit))
	}
	mask := byte(0x80 >> uint(bit))
	if this.buf[idx+global.Base]&mask == 0 {
		return false
	} else {
		return true
	}
}

func (this *IndexLevel) BitStart() int {
	s := this.Start()
	for i := 0; i != 8; i++ {
		if this.CheckBit(s, i) {
			start := (s << 3) + i
			return start
		}
	}
	panic(this.key + ", IndexLevel.BitStart() all zero")
}

func (this *IndexLevel) BitEnd() int {
	e := this.End() - 1
	for i := 7; i >= 0; i-- {
		if this.CheckBit(e, i) {
			end := (e << 3) + i + 1
			return end
		}
	}
	panic(this.key + ", IndexLevel.BitEnd() all zero")
}

func (this *IndexLevel) BitRange() (start int, end int) {
	/* s < e */
	start = this.BitStart()
	end = this.BitEnd()
	return
}

func (this *IndexLevel) NextLevelKey(bitoff int) (key string) {
	switch this.level {
	case 1:
		key = global.L2Key(this.key, bitoff)
	case 2:
		key = global.L3Key(this.key, bitoff)
	default:
		panic("level error: " + strconv.Itoa(this.level))
	}
	return
}

func (this *IndexLevel) SetNextLevel(uid int32, weight int32) *IndexLevel {
	l2bitOff := (int(uid) % global.L1BitUsrs) / global.L2BitUsrs
	switch this.level {
	case 1:
		bitOff := int(uid) / global.L1BitUsrs
		key := this.NextLevelKey(bitOff)
		return SetLevel2(key, l2bitOff)

	case 2:
		key := this.NextLevelKey(l2bitOff)
		itmNum := (int(uid) % global.L1BitUsrs) % global.L2BitUsrs
		SetLevel3(key, itmNum, weight)
		return nil

	default:
		panic("SetNextLevel error level: " + strconv.Itoa(this.level))
	}
}

type DataLevel struct {
	Level
}

func NewDataLevel() *DataLevel {
	l := new(DataLevel)
	l.bufSize = global.DataBufSize
	l.buf = make([]byte, l.bufSize)
	l.SetStart(global.ItemNum - 1)
	l.SetEnd(0)
	return l
}

func (this *DataLevel) SetStart(itm int) {
	this.buf[0] = byte((itm >> 8) & 0xFF)
	this.buf[1] = byte(itm & 0xFF)
	this.start = itm
}

func (this *DataLevel) SetEnd(itm int) {
	this.buf[2] = byte((itm >> 8) & 0xFF)
	this.buf[3] = byte(itm & 0xFF)
	this.end = itm
}

func (this *DataLevel) SetItm(itm int, weight int32) bool {
	if itm < 0 || itm >= global.ItemNum {
		panic("SetItm error: " + strconv.Itoa(itm))
	}
	idx := (itm << 2) + global.Base
	changed := false
	w0 := byte((weight >> 24) & 0xFF)
	w1 := byte((weight >> 16) & 0xFF)
	w2 := byte((weight >> 8) & 0xFF)
	w3 := byte(weight & 0xFF)
	if w0 != this.buf[idx] {
		this.buf[idx] = w0
		changed = true
	}
	if w1 != this.buf[idx+1] {
		this.buf[idx+1] = w1
		changed = true
	}
	if w2 != this.buf[idx+2] {
		this.buf[idx+2] = w2
		changed = true
	}
	if w3 != this.buf[idx+3] {
		this.buf[idx+3] = w3
		changed = true
	}
	this.start = (int(this.buf[0]) << 8) + int(this.buf[1])
	this.end = (int(this.buf[2]) << 8) + int(this.buf[3])
	if itm < this.start {
		this.SetStart(itm)
		changed = true
	}
	if itm >= this.end {
		this.SetEnd(itm + 1)
		changed = true
	}
	return changed
}

/* itmNum: 1, 2, ..., global.ItemNum-1 */
func (this *DataLevel) GetWeight(itmNum int) (weight int32, ok bool) {
	/* byte[0~3] is used to store Range */
	if itmNum < 0 || itmNum >= global.ItemNum {
		panic("GetWeight itmNum error: " + strconv.Itoa(itmNum))
	}
	idx := (itmNum << 2) + global.Base
	weight = (int32(this.buf[idx]) << 24) +
		(int32(this.buf[idx+1]) << 16) +
		(int32(this.buf[idx+2]) << 8) +
		int32(this.buf[idx+3])
	if weight == 0 {
		ok = false
	} else {
		ok = true
	}
	return
}

func GetLevel(key string, levels int, level LevelInterface) bool {
	val, err := kv.Get(&key)
	if err == kv.NotExist {
		return false
	}
	var bufSize int
	if levels == 1 || levels == 2 {
		bufSize = global.IndexBufSize
	} else if levels == 3 {
		bufSize = global.DataBufSize
	} else {
		panic("level error: " + strconv.Itoa(levels))
	}
	if len(val) != bufSize {
		panic("Level: " + strconv.Itoa(levels) +
			", key: " + key +
			", val size: " + strconv.Itoa(len(val)) +
			", correct bufSize: " + strconv.Itoa(bufSize))
	}
	level.Assign(key, levels, val)
	return true
}

func GetLevel1(id int32) *IndexLevel {
	key1 := global.L1Key(id)
	l1 := NewIndexLevel()
	if GetLevel(key1, 1, l1) == false {
		return nil
	}
	return l1
}

func SetLevel1(tid int32, uid int32) *IndexLevel {
	key1 := global.L1Key(tid)
	bitOff := int(uid) / global.L1BitUsrs
	l1 := NewIndexLevel()
	if GetLevel(key1, 1, l1) == false {
		l1.SetKey(key1)
		l1.SetLevels(1)
	}
	if l1.SetBit(bitOff) {
		val1 := string(l1.buf[:])
		kv.Put(&key1, &val1)
	}
	l1.SetKey(key1)
	l1.SetLevels(1)
	return l1
}

func GetLevel2(key string) *IndexLevel {
	l2 := NewIndexLevel()
	if GetLevel(key, 2, l2) == false {
		return nil
	}
	return l2
}

func SetLevel2(key2 string, bitOff int) *IndexLevel {
	l2 := NewIndexLevel()
	if GetLevel(key2, 2, l2) == false {
		l2.SetKey(key2)
		l2.SetLevels(2)
	}

	if l2.SetBit(bitOff) {
		val2 := string(l2.buf[:])
		kv.Put(&key2, &val2)
	}
	return l2
}

func GetLevel3(key string) *DataLevel {
	l3 := NewDataLevel()
	if GetLevel(key, 3, l3) == false {
		return nil
	}
	return l3
}

func SetLevel3(key3 string, itm int, weight int32) *DataLevel {
	l3 := NewDataLevel()
	if GetLevel(key3, 3, l3) == false {
		l3.SetKey(key3)
		l3.SetLevels(3)
	}
	if l3.SetItm(itm, weight) {
		val3 := string(l3.buf[:])
		kv.Put(&key3, &val3)
	}
	return l3
}
