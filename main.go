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
