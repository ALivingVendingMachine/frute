package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/alivingvendingmachine/frute/art"
	"github.com/alivingvendingmachine/frute/fuzzer"
	"github.com/alivingvendingmachine/frute/requester"
	"github.com/alivingvendingmachine/frute/util"
	"github.com/kortschak/ct"
)

//TODO: if fuzzing and you have the same sentinel twice, something bad happens

type headerFlags []string

func (h *headerFlags) String() string {
	return ""
}

func (h *headerFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

var (
	//traceLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger

	traceColor = ct.Fg(ct.Cyan).Paint
	infoColor  = ct.Fg(ct.Green).Paint
	warnColor  = ct.Fg(ct.Yellow).Paint
	errorColor = ct.Fg(ct.Red).Paint

	helpFlag     bool
	fuzzFlag     bool
	asciiFlag    bool
	generateFlag bool
	itersFlag    int
	seedFlag     int64
	urlFlag      string
	methodFlag   string
	bodyFlag     string
	requestFlag  string
	sentsFlag    string
	headers      headerFlags
)

func printUsage() {
	fmt.Println("usage:")
	fmt.Println("  \tfrute [-h] [--url url --method method [--body body] [-header header_pair]] request_output_file")
	fmt.Println("or\tfrute -r request_file -s sentinel_file -f")
	fmt.Printf("\n")
	fmt.Println("\t-h --help\n\t  print this usage message, then quit")
	fmt.Printf("\n")
	fmt.Println("\t-u --url\n\t  specified url to use while generating a new request text file")
	fmt.Println("\t-m --method\n\t  the method to use on the new request")
	fmt.Println("\t-b --body\n\t  the body of the request to use")
	fmt.Println("\t-H --header\n\t  headers to include while generating a request")
	fmt.Printf("\n")
	fmt.Println("\t-r --request\n\t  specify path to text file which contains a request to send")
	fmt.Println("\t-s --sentinel\n\t  specify path to text file which contains a list of sentinels to replace while bruting/fuzzing")
	fmt.Println("\t-f\n\t  fuzz string at sentinels")
}

func init() {
	flag.Usage = printUsage
	//traceLog = log.New(os.Stderr, fmt.Sprint(traceColor("TRACE: ")), log.Ldate|log.Ltime|log.Lshortfile)
	infoLog = log.New(os.Stderr, fmt.Sprint(infoColor("INFO: ")), log.Ldate|log.Ltime)
	warnLog = log.New(os.Stderr, fmt.Sprint(warnColor("WARN: ")), log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, fmt.Sprint(errorColor("ERROR: ")), log.Ldate|log.Ltime)

	const (
		helpUsage     = "prints usage, then exits"
		urlUsage      = "url to generate a request to"
		methodUsage   = "method to use while generating request"
		bodyUsage     = "request body to use while generating request"
		headerUsage   = "headers to include while generating request"
		requestUsage  = "request to use while fuzzing/bruteforcing"
		sentsUsage    = "filepath to file containing sentinels"
		fuzzUsage     = "set to fuzzing"
		generateUsage = "let fuzzer generate new strings at sentinels"
		asciiUsage    = "limit fuzzer to ascii only"
		itersUsage    = "number of iterations for the fuzzer to use"
		seedUsage     = "seed to pass to fuzzer, default is current time"
	)
	flag.BoolVar(&helpFlag, "h", false, helpUsage+" (shorthand)")
	flag.BoolVar(&helpFlag, "help", false, helpUsage)

	flag.StringVar(&urlFlag, "u", "", urlUsage+" (shorthand)")
	flag.StringVar(&urlFlag, "url", "", urlUsage)

	flag.StringVar(&methodFlag, "m", "", methodUsage+" (shorthand)")
	flag.StringVar(&methodFlag, "method", "", methodUsage)

	flag.StringVar(&bodyFlag, "b", "", bodyUsage+" (shorthand)")
	flag.StringVar(&bodyFlag, "body", "", bodyUsage)

	flag.Var(&headers, "H", headerUsage+" (shorthand)")
	flag.Var(&headers, "header", headerUsage)

	flag.StringVar(&requestFlag, "r", "", requestUsage)
	flag.StringVar(&requestFlag, "request", "", requestUsage)

	flag.StringVar(&sentsFlag, "s", "", sentsUsage)
	flag.StringVar(&sentsFlag, "sentinel", "", sentsUsage)

	flag.BoolVar(&fuzzFlag, "f", false, fuzzUsage)
	flag.BoolVar(&asciiFlag, "A", false, asciiUsage)
	flag.BoolVar(&generateFlag, "G", false, generateUsage)

	flag.IntVar(&itersFlag, "I", 3, itersUsage)
	flag.Int64Var(&seedFlag, "S", 1234567890, seedUsage)
}

func main() {
	flag.Parse()
	art.DrawArt()

	var sents []string
	//found := false

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}
	if len(flag.Args()) == 0 && urlFlag == "" && requestFlag == "" {
		flag.Usage()
		os.Exit(0)
	}

	// generate request
	if urlFlag != "" && methodFlag == "" {
		errorLog.Println("method cannot be blank!")
		flag.Usage()
		os.Exit(1)
	}
	if urlFlag != "" && len(flag.Args()) == 0 {
		errorLog.Println("no output file for request")
		flag.Usage()
		os.Exit(1)
	}
	if urlFlag != "" && bodyFlag == "" {
		warnLog.Println("request body is blank")
	}

	if urlFlag != "" {
		infoLog.Println("generating request")
		err := util.GenerateRequest(methodFlag, urlFlag, bodyFlag, headers, flag.Args()[0])
		if err != nil {
			errorLog.Printf("%v\n", err)
		}
		os.Exit(0)
	}

	infoLog.Println("fuzzing/bruting")

	// brute/fuzz
	var ret string
	var request *http.Request
	if requestFlag == "" {
		errorLog.Println("request path not set")
		os.Exit(1)
	}
	if sentsFlag == "" {
		warnLog.Println("sentinel path not set, using defaults")
		sents = []string{"!!!", "@@@", "###", "$$$", "^^^", "&&&", "***", "(((", ")))", "___"}
	} else if sentsFlag != "" {
		var err error
		sents, err = readSents(sentsFlag)
		if err != nil {
			errorLog.Printf("%v\n", err)
			os.Exit(1)
		}
	}

	if fuzzFlag { // fuzzing!
		if seedFlag == 123456790 {
			seedFlag = time.Now().UTC().UnixNano()
		}
		if asciiFlag {
			errorLog.Println("not implemented")
			os.Exit(1)
		}
		if generateFlag {
			errorLog.Println("not implemented")
			os.Exit(1)
		}

		req, err := ioutil.ReadFile(requestFlag)
		if err != nil {
			errorLog.Printf("%v\n", err)
			os.Exit(1)
		}

		ret = string(req)
		for i := 0; i < len(sents); i++ {
			var err error
			thisRet, err := fuzzer.MutateSelection(ret, sents[i], seedFlag, itersFlag)
			if err != nil {
				break
			}
			ret = thisRet
		}
		request, err = requester.ReadStringReq(ret)
		if err != nil {
			errorLog.Printf("%v\n", err)
			os.Exit(1)
		}
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		errorLog.Printf("%v\n", err)
		os.Exit(1)
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		errorLog.Printf("%v\n", err)
	}
	fmt.Printf("%q\n", dump)
}

func readSents(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ret []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	return ret, scanner.Err()
}
