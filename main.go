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

	"github.com/alivingvendingmachine/frute/brute"

	"github.com/alivingvendingmachine/frute/art"
	"github.com/alivingvendingmachine/frute/fuzzer"
	"github.com/alivingvendingmachine/frute/requester"
	"github.com/alivingvendingmachine/frute/util"
	"github.com/kortschak/ct"
)

//TODO: if fuzzing and you have the same sentinel twice, something bad happens
//TODO: if the sentinels are out of order, boom

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
	timesFlag    int
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
	fmt.Println("or\tfrute -f -r request_file [-s sentinel_file] [-times number_of_times")
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
	fmt.Println("\t-times\n\t  times to repeat the fuzzing and sending")
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
		timesUsage    = "number of times to repeat the fuzzing and sending"
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
	flag.IntVar(&timesFlag, "times", 1, timesUsage)
	flag.Int64Var(&seedFlag, "S", 1234567890, seedUsage)
}

func main() {
	flag.Parse()
	art.DrawArt()

	var sents []string
	newSeed := false
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
	if requestFlag == "" {
		errorLog.Println("request path not set")
		os.Exit(1)
	}
	if sentsFlag == "" {
		warnLog.Println("sentinel path not set, using defaults")
		sents = []string{"~~~", "!!!", "@@@", "###", "&&&", "<<<", ">>>", "___", ",,,", "'''"}
	} else if sentsFlag != "" {
		var err error
		sents, err = readSents(sentsFlag)
		if err != nil {
			errorLog.Printf("%v\n", err)
			os.Exit(1)
		}
	}
	if timesFlag < 1 {
		errorLog.Println("times to fuzz cannot be less than 1")
		printUsage()
	}

	if fuzzFlag { // fuzzing!
		if seedFlag == 1234567890 {
			seedFlag = time.Now().UTC().UnixNano()
			newSeed = true
		}
		if asciiFlag {
			errorLog.Println("not implemented")
			os.Exit(1)
		}
		if generateFlag {
			errorLog.Println("not implemented")
			os.Exit(1)
		}

		// loop here
		for i := 0; i < timesFlag; i++ {
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
			resp, err := doRequest(ret)
			if err != nil {
				errorLog.Printf("%v\n", err)
			}
			dump, err := httputil.DumpResponse(resp, true)
			if err != nil {
				errorLog.Printf("%v\n", err)
			}
			if newSeed {
				seedFlag = time.Now().UTC().UnixNano()
			}
			fmt.Printf("%s\n", dump)
			// end loop
		}
	} else { //bruting
		req, err := ioutil.ReadFile(requestFlag)
		if err != nil {
			errorLog.Printf("%v\n", err)
			os.Exit(1)
		}

		ret = string(req)

		offs := make([]int64, len(flag.Args()))
		fps, err := util.BulkOpen(flag.Args())
		if err != nil {
			errorLog.Printf("brute: %v", err)
			os.Exit(1)
		}

		var exhausted = false
		var out []string
		out, exhausted, err = util.ReadInputs(fps, offs)
		for !exhausted {
			if err != nil {
				errorLog.Printf("brute: %v", err)
				os.Exit(1)
			}
			ret, err = brute.Forcer(ret, out, sents)
			fmt.Println("doing:")
			fmt.Println(ret)
			resp, err := doRequest(ret)
			if err != nil {
				errorLog.Printf("%v", err)
				os.Exit(1)
			}
			dump, err := httputil.DumpResponse(resp, true)
			if err != nil {
				errorLog.Printf("%v", err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", dump)
			out, exhausted, err = util.ReadInputs(fps, offs)
		}
	}
}

func doRequest(request string) (*http.Response, error) {
	req, err := requester.ReadStringReq(request)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	return resp, err
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
