package main

import "github.com/mitchellh/ioprogress"
import "github.com/docopt/docopt-go"
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
	usage := `goget.
	Usage:
		goget <uri>
		goget <uri> [--replace]
		goget <uri> [--name=<filename>]
		goget <uri> [--replace|--name=<filename>]
		goget -h | --help
		goget --version

	Options:
		-h --help            Show this screen.
		--version            Show version.
		-r --replace         Replace exist file.
		-n --name=<filename> Set file name.
	`

	args, _ := docopt.Parse(usage, os.Args[1:], true, "v0.1.0", false)
	fmt.Println(args)

	var uri, filename string
	var replace bool

	i, ok := args["<uri>"]
	fmt.Println(i)
	uri = i.(string)

	i, ok = args["--replace"].(bool)
	if ok {
		replace = i.(bool)
	}

	i, ok = args["--name"]
	if ok {
		filename = i.(string)
	}

	get(uri, replace, filename)
}

func get(uri string, replace bool, fname string) {
	_, e := url.ParseRequestURI(uri)
	if e != nil {
		fmt.Print("invalid url")
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

	if fname == "" {
		filename = getFilename(uri, mediaType)
	} else {
		filename = fname
	}

	if exist(filename) {
		if replace {
			os.Remove(filename)
		} else {
			fmt.Errorf("file: %s exist", filename)
			os.Exit(1)
		}
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
