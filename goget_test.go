package main

import "testing"
import "os"

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
	f1 := getFilename("http://example.com/a.png")
	f2 := getFilename("https://example.cn/a.png")
	f3 := getFilename("example.com/a.png")
	f4 := getFilename("example.cn/")
	f5 := getFilename("example.com")

	f6 := getFilename("")
	f7 := getFilename("http://example.com/a.png?name=xxoo")
	f8 := getFilename("http://example.com/a.png#name=xxoo")
	f9 := getFilename("http://example.com/a.png?name=xxoo#name=xxoo")

	t.Log(f1, f2, f3, f4, f5)
	t.Log(f6, f7, f8, f9)

	equal(f1, "a.png", t)
	equal(f2, "a.png", t)
	equal(f3, "a.png", t)
	equal(f4, "example.cn", t)
	equal(f5, "example.com", t)

	equal(f6, "goget-download", t)
	equal(f7, "a.png", t)
	equal(f8, "a.png", t)
	equal(f9, "a.png", t)
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

func TestGetByHttp(t *testing.T) {
	url := "http://h.hiphotos.baidu.com/image/pic/item/e1fe9925bc315c60cb71c4748fb1cb1348547741.jpg"
	if exist(getFilename(url)) {
		os.Remove(getFilename(url))
	}

	get(url)

	if !exist(getFilename(url)) {
		t.Fail()
	}
}

func TestGetByHttps(t *testing.T) {
	url := "https://avatars1.githubusercontent.com/u/2569835"
	filename := "2569835.png"
	if exist(filename) {
		os.Remove(filename)
	}

	get(url)

	if !exist(filename) {
		t.Fail()
	}
}

// utils
func equal(s1, s2 string, t *testing.T) {
	if s1 != s2 {
		t.Error(s1, "not equal", s2)
	}
}
