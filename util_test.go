package main

import "testing"

func TestExist(t *testing.T) {
	yes := exist("goget.go")
	if !yes {
		t.Fail()
	}

	yes = exist("xxoo.go")
	if yes {
		t.Fail()
	}
}

func TestGetFilename(t *testing.T) {
	f1 := getFilename("http://example.com/a.png", "")
	f2 := getFilename("https://example.cn/a.png", "")
	f3 := getFilename("example.com/a.png", "")
	f4 := getFilename("example.cn/", "")
	f5 := getFilename("example.com", "")

	f6 := getFilename("", "")
	f7 := getFilename("http://example.com/a.png?name=xxoo", "")
	f8 := getFilename("http://example.com/a.png#name=xxoo", "")
	f9 := getFilename("http://example.com/a.png?name=xxoo#name=xxoo", "")
	f0 := getFilename("http://example.com/a?name=xxoo#name=xxoo", "png")

	t.Log(f1, f2, f3, f4, f5)
	t.Log(f6, f7, f8, f9, f0)

	equal(f1, "a.png", t)
	equal(f2, "a.png", t)
	equal(f3, "a.png", t)
	equal(f4, "example.cn", t)
	equal(f5, "example.com", t)

	equal(f6, "goget-download", t)
	equal(f7, "a.png", t)
	equal(f8, "a.png", t)
	equal(f9, "a.png", t)
	equal(f0, "a.png", t)
}

func TestCutAfter(t *testing.T) {
	equal(cutAfter("xx/oo", "/"), "xx", t)
	equal(cutAfter("xx", "/"), "xx", t)
}

func TestCutBefore(t *testing.T) {
	equal(cutBefore("xx/oo", "/"), "oo", t)
	equal(cutBefore("xx", "/"), "xx", t)
}

func TestByteUnitString(t *testing.T) {
	if byteUnitString(1000) != "1 KB" {
		t.Fail()
	}
}
