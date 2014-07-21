package search

import (
	"strconv"

	"global"
	"level"
)

func FuzzyRange(tid int32) (start int, end int) {
	l1 := level.GetLevel1(tid)
	if l1 == nil {
		return
	}

	/* s < e, BitRange[s, e) */
	s, e := l1.BitRange()
	start = s * global.L1BitUsrs
	end = e*global.L1BitUsrs - 1
	return
}

func Range(tid int32) (start int, end int) {
	l1 := level.GetLevel1(tid)
	if l1 == nil {
		return
	}

	s, e := l1.BitRange()
	if s == e {
		return
	}

	key21 := l1.NextLevelKey(s)
	startBase := s * global.L1BitUsrs
	key22 := l1.NextLevelKey(e - 1)
	endBase := (e - 1) * global.L1BitUsrs

	l21 := level.GetLevel2(key21)
	if l21 == nil {
		panic("tid error in level l21: " + strconv.Itoa(int(tid)))
	}

	if s+1 == e {
		s, e = l21.BitRange()
		start = s*global.L2BitUsrs + startBase
		end = e*global.L2BitUsrs - 1 + endBase
		return
	}

	s = l21.BitStart()
	start = s*global.L2BitUsrs + startBase

	l22 := level.GetLevel2(key22)
	if l22 == nil {
		panic("tid error in level l22: " + strconv.Itoa(int(tid)))
	}
	e = l22.BitEnd()
	end = e*global.L2BitUsrs - 1 + endBase

	return
}

type Usr struct {
	Uid   int32
	Score int32
}

func GetUsrByRange(tid int32, start int, end int, limit int) (u []Usr) {
	var uSize int
	if limit > 1024 {
		uSize = 1024
	} else {
		uSize = limit + 1
	}
	u = make([]Usr, 0, uSize)
	l1 := level.GetLevel1(tid)
	if l1 == nil {
		return
	}

	l1Start, l1End := l1.BitRange()
	if l1Start < start/global.L1BitUsrs {
		l1Start = start / global.L1BitUsrs
	}
	if l1End >= end/global.L1BitUsrs {
		l1End = (end / global.L1BitUsrs) + 1
	}

	for i := l1Start; i != l1End; i++ {
		if !l1.CheckBitFlat(i) {
			continue
		}
		l2Key := l1.NextLevelKey(i)
		l2 := level.GetLevel2(l2Key)
		if l2 == nil {
			panic("l2Key = " + l2Key + ", l2 == NULL")
		}
		l2Start, l2End := l2.BitRange()

		for j := l2Start; j != l2End; j++ {
			if !l2.CheckBitFlat(j) {
				continue
			}
			base := i*global.L1BitUsrs + j*global.L2BitUsrs
			if base+global.ItemNum < start || base > end {
				continue
			}
			l3Key := l2.NextLevelKey(j)
			l3 := level.GetLevel3(l3Key)
			if l3 == nil {
				panic("l3Key = " + l3Key + ", l3 == NULL")
			}

			/* l3Range: [0, ItemNum) */
			l3Start, l3End := l3.Range()
			for k := l3Start; k < l3End; k++ {
				var usr Usr
				weight, ok := l3.GetWeight(k)
				if !ok {
					continue
				}
				uid := base + k
				if uid < start || uid >= end {
					continue
				}
				usr.Uid = int32(uid)
				usr.Score = weight
				u = append(u, usr)
				if len(u) > limit {
					return
				}
			}
		}
	}
	return
}
