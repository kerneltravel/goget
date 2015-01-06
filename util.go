package main

import "strings"
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
