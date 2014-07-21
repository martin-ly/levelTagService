package kv

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include "leveldb/c.h"

#cgo LDFLAGS: -lleveldb -lpthread

static char *errmsg;
static char *dbpath;
static leveldb_t *db;
static leveldb_options_t *opt;
static leveldb_readoptions_t *ropt;
static leveldb_writeoptions_t *wopt;

static void openDb(const char *name)
{
    leveldb_options_t  *opt;

    opt = leveldb_options_create();
    leveldb_options_set_create_if_missing(opt, (unsigned char)1);
    ropt = leveldb_readoptions_create();
    wopt = leveldb_writeoptions_create();

    dbpath = strdup(name);
    db = leveldb_open(opt, dbpath, &errmsg);

    if (errmsg != NULL) {
        printf("open levelDb %s error: %s\n", dbpath, errmsg);
        leveldb_free(errmsg);
        errmsg = NULL;
        abort();
    }
}

static void closeDb()
{
    if (errmsg != NULL) {
        leveldb_free(errmsg);
        errmsg = NULL;
    }
    if (opt != NULL) {
        leveldb_options_destroy(opt);
        opt = NULL;
    }
    if (ropt != NULL) {
        leveldb_readoptions_destroy(ropt);
        ropt = NULL;
    }
    if (wopt != NULL) {
        leveldb_writeoptions_destroy(wopt);
        wopt = NULL;
    }
    if (dbpath != NULL) {
        free(dbpath);
        dbpath = NULL;
    }
    if (db != NULL) {
        leveldb_close(db);
        db = NULL;
    }
}

static int put(const char *k, size_t k_len, const char *v, size_t v_len)
{
    leveldb_put(db, wopt, k, k_len, v, v_len, &errmsg);

    if (errmsg != NULL) {
        printf("put k, v failed: %s\n", errmsg);
        leveldb_free(errmsg);
        errmsg = NULL;
        return -1;
    }
    return 0;
}

static char *get(const char *k, size_t k_len, size_t *v_len)
{
    char                   *val;

    val = leveldb_get(db, ropt, k, k_len, v_len, &errmsg);

    if (errmsg != NULL) {
        printf("get k failed: %s\n", errmsg);
        leveldb_free(errmsg);
        errmsg = NULL;
    }
    return val;
}

static int del(const char *k, size_t k_len)
{
    leveldb_delete(db, wopt, k, k_len, &errmsg);

    if (errmsg != NULL) {
        printf("delete k failed: %s\n", errmsg);
        leveldb_free(errmsg);
        errmsg = NULL;
        return -1;
    }
    return 0;
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

var PutError error = errors.New("levelDb put failed")
var NotExist error = errors.New("key not exist")
var DelError error = errors.New("levelDb del failed")

func Open(name *string) {
	cname := C.CString(*name)
	defer C.free(unsafe.Pointer(cname))
	C.openDb(cname)
}

func Close() {
	C.closeDb()
}

func Put(k *string, v *string) error {
	key := C.CString(*k)
	defer C.free(unsafe.Pointer(key))
	val := C.CString(*v)
	defer C.free(unsafe.Pointer(val))

	ok := C.put(key, C.size_t(len(*k)), val, C.size_t(len(*v)))
	if ok != 0 {
		return PutError
	}
	return nil
}

func Get(k *string) (v string, err error) {
	var v_len C.size_t

	key := C.CString(*k)
	defer C.free(unsafe.Pointer(key))

	val := C.get(key, C.size_t(len(*k)), &v_len)
	if val == nil {
		err = NotExist
		return
	}
	defer C.leveldb_free(unsafe.Pointer(val))

	v = C.GoStringN(val, C.int(v_len))
	return
}

func Delete(k *string) error {
	key := C.CString(*k)
	defer C.free(unsafe.Pointer(key))

	ok := C.del(key, C.size_t(len(*k)))
	if ok != 0 {
		return DelError
	}
	return nil
}
