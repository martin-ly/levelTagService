package update

import (
	"level"
)

func update(tid int32, uid int32, weight int32, toDel bool) {
	if toDel {
		// TODO: waiting for impl
		return
	}

	// Insert
	l1 := level.SetLevel1(tid, uid)
	l2 := l1.SetNextLevel(uid, weight)
	l2.SetNextLevel(uid, weight)
}

func Insert(tid int32, uid int32, weight int32) {
	update(tid, uid, weight, false)
}

func Delete(tid int32, uid int32) {
	update(tid, uid, 0, true)
}
