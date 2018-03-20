// go test *.go -bench=".*" -benchmem

package drh

import (
    "testing"
    "gitee.com/johng/gf/g/container/gdrh"
)

var drh  = gdrh.New(10000, 10000)
var gom  = make(map[int]interface{})
var size = 10000000

func BenchmarkDrh_Set(b *testing.B) {
    b.N = size
    for i := 0; i < b.N; i++ {
        drh.Set(i, i)
    }
}

func BenchmarkDrh_Get(b *testing.B) {
    b.N = size
    for i := 0; i < b.N; i++ {
        drh.Get(i)
    }
}

func BenchmarkDrh_Remove(b *testing.B) {
    b.N = size
    for i := 0; i < b.N; i++ {
        drh.Remove(i)
    }
}

func BenchmarkGoMap_Set(b *testing.B) {
    b.N = size
    for i := 0; i < b.N; i++ {
        gom[i] = i
    }
}

func BenchmarkGoMap_Get(b *testing.B) {
    b.N = size
    for i := 0; i < b.N; i++ {
        if _, ok := gom[i]; ok {}
    }
}

func BenchmarkGoMap_Remove(b *testing.B) {
    b.N = size
    for i := 0; i < b.N; i++ {
        delete(gom, i)
    }
}