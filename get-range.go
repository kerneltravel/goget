package main

import "io/ioutil"
import "net/http"
import "errors"
import "fmt"
import "os"

type Range struct {
	uri      string
	filename string
	point    int64
	total    int64
}

func newRange(uri string, replace bool, fname string) (r *Range, err error) {
	checkUri(uri)

	supported, header := supportRange(uri)

	if !supported {
		return nil, errors.New("range not support")
	}

	var total int64
	contentRange := header.Get("Content-Range")
	if contentRange != "" {
		_, _, total = parseRangeString(contentRange)
	}

	var filename string
	if fname != "" {
		filename = fname
	} else {
		filename = parseFilename(uri, header)
	}

	recover(filename, replace)

	return &Range{
		uri:      uri,
		filename: filename,
		total:    total,
		point:    0,
	}, nil
}

func (r *Range) resume(size int64) (success, finished bool) {
	start := r.point
	end := start + size

	if r.total != 0 && end > r.total {
		end = r.total
	}
	if start >= end {
		debug("total: %d, start: %d, end: %d", r.total, start, end)
		return true, true
	}

	rangeString := getRangeString(start, end)

	client := &http.Client{}

	req, err := http.NewRequest("GET", r.uri, nil)
	if err != nil {
		debug("http get error: %v", err)
		return false, false
	}
	req.Header.Set("Range", rangeString)

	res, err := client.Do(req)
	if err != nil {
		debug("client do error: %v", err)
		return false, false
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		debug("read all error: %v", err)
		return false, false
	}
	_, err = appendFile(r.filename, data, start)

	// set point
	r.point = end

	// check end
	if res.Header.Get("Content-Range") != "" {
		_, e, t := parseRangeString(res.Header.Get("Content-Range"))

		if t != 0 && r.point >= t {
			return true, true
		}

		if e <= r.point {
			return true, true
		}
	}

	return true, false
}

// recover from the exist file
func recover(filename string, replace bool) int64 {
	if !exist(filename) {
		return 0
	}

	if replace {
		os.Remove(filename)
		return 0
	} else {
		info, _ := os.Stat(filename)
		return info.Size()
	}
}

func supportRange(uri string) (supported bool, header http.Header) {
	res, err := http.Head(uri)

	if err != nil {
		return false, nil
	}
	if res.StatusCode > 300 {
		return false, nil
	}

	if res.Header.Get("Accept-Ranges") == "bytes" {
		return true, res.Header
	}

	return false, nil
}

// Example:
//   "Content-Range": "bytes 100-200/1000"
//   "Content-Range": "bytes 100-200/*"
func parseRangeString(r string) (start, end, total int64) {
	fmt.Sscanf(r, "bytes %d-%d/%d", &start, &end, &total)

	if total != 0 && end > total {
		end = total
	}
	if start >= end {
		start = 0
		end = 0
	}

	return
}

// Example:
//   "Range": "bytes=100-200"
func getRangeString(start, end int64) string {
	return "bytes=" + int64toString(start) + "-" + int64toString(end)
}
