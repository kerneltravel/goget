package main

import "net/http"
import "net/url"
import "strconv"
import "strings"
import "mime"
import "path"
import "fmt"
import "os"

func exist(filename string) bool {
	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func getFilename(url string, mediaType string) string {
	n := path.Base(url)

	if mediaType == "" && (n == "" || n == ".") {
		return "goget-download"
	}

	name1 := cutAfter(n, "#")
	name := cutAfter(name1, "?")

	if path.Ext(name) == "" && mediaType != "" {
		return name + "." + mediaType
	}

	return name
}

func parseFilename(uri string, header http.Header) string {
	contentType := header.Get("Content-Type")
	mimeType, _, _ := mime.ParseMediaType(contentType)
	mediaType := cutBefore(mimeType, "/")

	filename := getFilename(uri, mediaType)

	return filename
}

func cutAfter(s, sep string) string {
	if strings.Contains(s, sep) {
		return strings.Split(s, sep)[0]
	}

	return s
}

func cutBefore(s, sep string) string {
	if strings.Contains(s, sep) {
		return strings.Split(s, sep)[1]
	}

	return s
}

var byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB"}

func byteUnitString(n int64) string {
	var unit string
	size := float64(n)
	for i := 1; i < len(byteUnits); i++ {
		if size < 1000 {
			unit = byteUnits[i-1]
			break
		}

		size = size / 1000
	}

	return fmt.Sprintf("%.3g %s", size, unit)
}

func int64toString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func string2int64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func string2int(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func appendFile(filename string, data []byte, offset int64) (int64, error) {
	var file *os.File
	var err error

	if !exist(filename) {
		file, err = os.Create(filename)
	} else {
		file, err = os.OpenFile(filename, os.O_RDWR, 0666)
	}
	defer file.Close()

	if err != nil {
		debug("file: %v, error: %v", file, err)
		return 0, err
	}
	debug("append file, data length: %d, offset: %d", len(data), offset)
	size, err := file.WriteAt(data, offset)

	if err != nil {
		debug("error: %v", err)
		return 0, err
	}

	return int64(size) + offset, nil
}

func checkFile(filename string, replace bool) {
	if exist(filename) {
		if replace {
			os.Remove(filename)
		} else {
			fmt.Printf("file: %s exist \n", filename)
			os.Exit(1)
		}
	}
}

func checkUri(uri string) {
	_, e := url.ParseRequestURI(uri)
	if e != nil {
		fmt.Println("invalid url")
		os.Exit(1)
	}
}
