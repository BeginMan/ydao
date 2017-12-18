package main

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

func TestIsChinese(t *testing.T) {
    Convey("Init a Chinese string", t, func() {
        Convey("Call isChinese func", func() {
            result := isChinese("测试")
            Convey("result should be: true", func() {
                So(result, ShouldEqual, true)
            })

            resultFalse := isChinese("test")
            Convey("result should by: false", func() {
                So(resultFalse, ShouldEqual, false)
            })
        })
    })
}

