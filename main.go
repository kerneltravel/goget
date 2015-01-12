package main

import "github.com/mitchellh/ioprogress"
import "github.com/docopt/docopt-go"
import . "github.com/tj/go-debug"
import "net/http"
import "strconv"
import "fmt"
import "io"
import "os"

var debug = Debug("goget")

const duration int64 = 1000 * 1000

func main() {
	usage := `
	Usage:
		goget <uri> [--replace] [--name=<filename>] [--continue]
		goget --help
		goget --version

	Options:
		-r --replace         Replace exist file
		-n --name=<filename> Set file name
		-c --continue        Resume getting a partially-downloaded file
		--resume             Resume from break point
		--help               Show this screen
		--version            Show version
	`

	args, _ := docopt.Parse(usage, os.Args[1:], true, "v0.2.0", false)

	uri := args["<uri>"].(string)
	replace := args["--replace"].(bool)
	partially := args["--continue"].(bool)

	name, ok := args["--name"]
	var filename string
	if ok && name != nil {
		filename = name.(string)
	}

	if partially {
		getPartially(uri, replace, filename)
		return
	}

	get(uri, replace, filename)
}

func get(uri string, replace bool, fname string) {
	checkUri(uri)

	res, err := http.Get(uri)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("http get error: %s \n", err.Error())
	}

	contentLength := res.Header.Get("Content-Length")

	var filename string
	if fname != "" {
		filename = fname
	} else {
		filename = parseFilename(uri, res.Header)
	}

	checkFile(filename, replace)

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

func getPartially(uri string, replace bool, filename string) {
	r, e := newRange(uri, replace, filename)

	if e != nil {
		panic(e)
	}

	for {
		r.resume(duration)
	}
}
