// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// go test *.go -bench=".*"

package gdrh_test

import (
    "testing"
    "gitee.com/johng/gf/g/container/gdrh"
)

var drh  = gdrh.New(100000, 100000)
var gom  = make(map[int]interface{})

func BenchmarkDrh_Set(b *testing.B) {
    b.N = 1000000
    for i := 0; i < b.N; i++ {
        drh.Set(i, i)
    }
}

func BenchmarkDrh_Get(b *testing.B) {
    b.N = 1000000
    for i := 0; i < b.N; i++ {
        drh.Get(i)
    }
}

func BenchmarkGoMap_Set(b *testing.B) {
    b.N = 1000000
    for i := 0; i < b.N; i++ {
        gom[i] = i
    }
}

func BenchmarkGoMap_Get(b *testing.B) {
    b.N = 1000000
    for i := 0; i < b.N; i++ {
        if _, ok := gom[i]; ok {}
    }
}