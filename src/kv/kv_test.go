package kv

import (
	"sync"
	"testing"
)

var once sync.Once

func TestKv(t *testing.T) {
	testDb := "testDb"
	once.Do(func() {
		Open(&testDb)
	})
	defer Close()

	k1, v1 := "liuwei-git:", "git@github.com:brg-liuwei"
	k2, v2 := "liuwei-com:", "www.qingtingfm.com"

	err := Put(&k1, &v1)
	if err != nil {
		t.Error("Put k1, v1 error: ", err)
	}

	var v string
	v, err = Get(&k1)
	if err != nil {
		t.Error("Get k1 error: ", err)
	}
	if v != v1 {
		t.Error("v = ", v, ", should be ", v1)
	}

	err = Put(&k1, &v2)
	if err != nil {
		t.Error("Put k1, v2 error: ", err)
	}
	v, err = Get(&k1)
	if err != nil {
		t.Error("Get k1 error: ", err)
	}
	if v != v2 {
		t.Error("v = ", v, ", should be ", v2)
	}

	err = Put(&k1, &v1)
	if err != nil {
		t.Error("Put k1, v1 error: ", err)
	}
	v, err = Get(&k1)
	if err != nil {
		t.Error("Get k1 error: ", err)
	}
	if v != v1 {
		t.Error("v = ", v, ", should be ", v1)
	}

	err = Delete(&k1)
	if err != nil {
		t.Error("delete k1 error: ", err)
	}

	v, err = Get(&k1)
	if err != NotExist {
		t.Error("k1 is deleted, err = ", err, ", v = ", v)
	}

	err = Delete(&k2)
	if err != nil {
		t.Error("delete k2(not exist) err = ", err, ", should be nil")
	}
}
