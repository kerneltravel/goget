package main

import "net/http/httptest"
import "io/ioutil"
import "net/http"
import "testing"
import "os"

func TestGetRange(t *testing.T) {
	ranges := []string{
		"0-10",
		"10-20",
		"20-30",
	}
	flag := 0

	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Range", "bytes "+ranges[flag]+"/30")
		res.Header().Set("Content-Type", "text/plain")
		res.Header().Set("Accept-Ranges", "bytes")
		res.WriteHeader(206)

		res.Write([]byte("abcdefghij"))
	}))

	// get
	res, _ := http.Get(ts.URL)
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	s := string(data[:])
	if s != "abcdefghij" {
		t.Error(s)
	}

	// get range
	r, e := newRange(ts.URL, true, "get-range.txt")
	if e != nil {
		t.Error(e)
	}

	flag = 0
	r.resume(10)
	flag = 1
	r.resume(10)
	flag = 2
	r.resume(10)
}

func TestRecover(t *testing.T) {
	// init test file
	filename := "recover.txt"
	testFile, _ := os.Create(filename)
	size, _ := testFile.WriteString("hello, golang")

	if recover(filename, false) != int64(size) {
		t.Fail()
	}
}

func TestGetRangeString(t *testing.T) {
	if getRangeString(100, 200) != "bytes=100-200" {
		t.Fail()
	}
}

func TestSupportRange(t *testing.T) {
	ts1 := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// none
	}))
	defer ts1.Close()
	ts2 := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("accept-ranges", "by")
	}))
	defer ts2.Close()
	ts3 := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Accept-Ranges", "by")
	}))
	defer ts3.Close()

	ts4 := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("accept-ranges", "bytes")
	}))
	defer ts4.Close()
	ts5 := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Accept-Ranges", "bytes")
	}))
	defer ts5.Close()

	r1, _ := supportRange(ts1.URL)
	r2, _ := supportRange(ts2.URL)
	r3, _ := supportRange(ts3.URL)

	r4, _ := supportRange(ts4.URL)
	r5, _ := supportRange(ts5.URL)

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

// utils for test
func int64equal(start1, end1, total1, start2, end2, total2 int64, t *testing.T) {
	if start1 != start2 || end1 != end2 || total1 != total2 {
		t.Error(start1, end1, total1)
	}
}
