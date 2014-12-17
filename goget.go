package main

import "github.com/mitchellh/ioprogress"
import "net/http"
import "net/url"
import "strings"
import "strconv"
import "path"
import "mime"
import "fmt"
import "io"
import "os"

func main() {
	uri := os.Args[1]
	get(uri)
}

func get(uri string) {
	_, e := url.ParseRequestURI(uri)
	if e != nil {
		fmt.Print("invalid url")
		os.Exit(1)
	}

	name := getFilename(uri)
	if exist(name) {
		fmt.Printf("file: %s exist", name)
		os.Exit(1)
	}

	res, err := http.Get(uri)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("http get error: %s", err.Error())
	}

	contentLength := res.Header.Get("Content-Length")
	contentType := res.Header.Get("Content-Type")
	mimeType, _, _ := mime.ParseMediaType(contentType)
	mediaType := cutBefore(mimeType, "/")

	var filename string
	if path.Ext(name) == "" && mediaType != "" {
		filename = name + "." + mediaType
	} else {
		filename = name
	}

	out, err := os.Create(filename)
	defer out.Close()
	if err != nil {
		fmt.Printf("create file error: %s", err.Error())
	}

	if contentLength != "" {
		size, err := strconv.ParseInt(contentLength, 10, 64)
		if err == nil && size > 0 {
			process := &ioprogress.Reader{
				DrawFunc: ioprogress.DrawTerminalf(os.Stdout, ioprogress.DrawTextFormatBytes),
				Reader:   res.Body,
				Size:     size,
			}

			fmt.Println("start")
			io.Copy(out, process)
			fmt.Printf("finished, size: %s.\n", byteUnitString(size))
			return
		}
	}

	fmt.Println("start")
	size, err := io.Copy(out, res.Body)
	fmt.Printf("finished, size: %s.\n", byteUnitString(size))
}

// utils
func exist(filename string) bool {
	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func getFilename(url string) string {
	name := path.Base(url)

	if name == "" || name == "." {
		return "goget-download"
	}

	name1 := cutAfter(name, "#")
	name2 := cutAfter(name1, "?")

	return name2
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
