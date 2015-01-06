package main

import "net/http/httptest"
import "io/ioutil"
import "testing"
import "os"

func TestExist(t *testing.T) {
	yes := exist("main.go")
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

func TestInt64toString(t *testing.T) {
	i := 64 // int64
	if int64toString(int64(i)) != "64" {
		t.Fail()
	}
}

func TestString2int64(t *testing.T) {
	if string2int64("64") != 64 {
		t.Fail()
	}
}

func TestString2int(t *testing.T) {
	if string2int("32") != 32 {
		t.Fail()
	}
}

func TestGetRangeString(t *testing.T) {
	if getRangeString(100, 200) != "bytes=100-200" {
		t.Fail()
	}
}

func TestCheckRangeSupport(t *testing.T) {
	h1 := httptest.NewRecorder().Header()
	h2 := httptest.NewRecorder().Header()
	h3 := httptest.NewRecorder().Header()

	h4 := httptest.NewRecorder().Header()
	h5 := httptest.NewRecorder().Header()

	h2.Add("accept-ranges", "by")
	h3.Add("Accept-Ranges", "by")

	h4.Add("accept-ranges", "bytes")
	h5.Add("Accept-Ranges", "bytes")

	r1 := checkRangeSupport(h1)
	r2 := checkRangeSupport(h2)
	r3 := checkRangeSupport(h3)

	r4 := checkRangeSupport(h4)
	r5 := checkRangeSupport(h5)

	result := [5]bool{r1, r2, r3, r4, r5}
	if result != [5]bool{false, false, false, true, true} {
		t.Error(result)
	}
}

func TestParseRangeString(t *testing.T) {
	start, end, total := parseRangeString("")
	int64equal(start, end, total, 0, 0, 0, t)

	start, end, total = parseRangeString("bytes 100-200/500")
	int64equal(start, end, total, 100, 200, 500, t)

	start, end, total = parseRangeString("bytes 100-200/*")
	int64equal(start, end, total, 100, 200, 0, t)

	start, end, total = parseRangeString("bytes 100-200")
	int64equal(start, end, total, 100, 200, 0, t)

	start, end, total = parseRangeString("bytes 200-200/500")
	int64equal(start, end, total, 0, 0, 500, t)

	start, end, total = parseRangeString("bytes 100-200/150")
	int64equal(start, end, total, 100, 150, 150, t)
}

func TestAppendFile(t *testing.T) {
	filename := "append-file.txt"

	if exist(filename) {
		os.Remove(filename)
	}

	offset, err := appendFile(filename, []byte("abcdef"), 0)
	t.Logf("offset: %d, error: %v", offset, err)
	offset, err = appendFile(filename, []byte("ghijkl"), offset)
	t.Logf("offset: %d, error: %v", offset, err)
	offset, err = appendFile(filename, []byte("mnopqr"), offset)
	t.Logf("offset: %d, error: %v", offset, err)

	if err != nil {
		t.Error(err)
	}

	data, _ := ioutil.ReadFile(filename)

	if string(data[:]) != "abcdefghijklmnopqr" {
		t.Fail()
	}
}

// utils for test
func equal(s1, s2 string, t *testing.T) {
	if s1 != s2 {
		t.Error(s1, "not equal", s2)
	}
}

func int64equal(start1, end1, total1, start2, end2, total2 int64, t *testing.T) {
	if start1 != start2 || end1 != end2 || total1 != total2 {
		t.Error(start1, end1, total1)
	}
}
