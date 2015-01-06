package main

import "github.com/mitchellh/ioprogress"
import "github.com/docopt/docopt-go"
import "net/http"
import "net/url"
import "strconv"
import "mime"
import "fmt"
import "io"
import "os"

func main() {
	usage := `
	Usage:
		goget <uri> [--replace] [--name=<filename>]
		goget --help
		goget --version

	Options:
		-r --replace         Replace exist file
		-n --name=<filename> Set file name
		--resume             Resume from break point
		--help               Show this screen
		--version            Show version
	`

	args, _ := docopt.Parse(usage, os.Args[1:], true, "v0.1.0", false)

	uri := args["<uri>"].(string)
	replace := args["--replace"].(bool)

	name, ok := args["--name"]
	var filename string
	if ok && name != nil {
		filename = name.(string)
	}

	get(uri, replace, filename)
}

func get(uri string, replace bool, fname string) {
	_, e := url.ParseRequestURI(uri)
	if e != nil {
		fmt.Println("invalid url")
		os.Exit(1)
	}

	res, err := http.Get(uri)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("http get error: %s \n", err.Error())
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
			fmt.Printf("file: %s exist \n", filename)
			os.Exit(1)
		}
	}

	out, err := os.Create(filename)
	defer out.Close()
	if err != nil {
		fmt.Printf("create file error: %s \n", err.Error())
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
			fmt.Printf("finished, size: %s. \n", byteUnitString(size))
			return
		}
	}

	fmt.Println("start")
	size, err := io.Copy(out, res.Body)
	fmt.Printf("finished, size: %s. \n", byteUnitString(size))
}
