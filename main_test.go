package main

import "testing"
import "os"

func TestGetByHttp(t *testing.T) {
	url := "http://h.hiphotos.baidu.com/image/pic/item/e1fe9925bc315c60cb71c4748fb1cb1348547741.jpg"
	if exist(getFilename(url, "")) {
		os.Remove(getFilename(url, ""))
	}

	get(url, false, "")
	get(url, true, "")
	get(url, true, "10086.png")

	if !exist(getFilename(url, "")) {
		t.Fail()
	}
}

func TestGetByHttps(t *testing.T) {
	url := "https://avatars1.githubusercontent.com/u/2569835"
	filename := "2569835.png"
	if exist(filename) {
		os.Remove(filename)
	}

	get(url, false, "")
	get(url, true, "")

	if !exist(filename) {
		t.Fail()
	}
}
