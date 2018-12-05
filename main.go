package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"

	"github.com/ALivingVendingMachine/frute/brute"
	"github.com/ALivingVendingMachine/frute/decoder"

	"github.com/ALivingVendingMachine/frute/art"
	"github.com/ALivingVendingMachine/frute/fuzzer"
	"github.com/ALivingVendingMachine/frute/requester"
	"github.com/ALivingVendingMachine/frute/util"
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

type respAndStrings struct {
	resp    *http.Response
	strings []string
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

	helpFlag        bool
	fuzzFlag        bool
	asciiFlag       bool
	generateFlag    bool
	randomWaitFlag  bool
	randomWaitRange int
	itersFlag       int
	timesFlag       int
	threadsFlag     int
	waitFlag        time.Duration
	urlFlag         string
	methodFlag      string
	bodyFlag        string
	requestFlag     string
	sentsFlag       string

	headers headerFlags
)

func printUsage() {
	fmt.Println("usage:")
	fmt.Println("  \tfrute [-h] [--url url --method method [--body body] [-header header_pair]] request_output_file")
	fmt.Println("or\tfrute -f -r request_file [-s sentinel_file] [-times number_of_times")
	fmt.Println("or\tfrute -r request_file [-s sentinel_file] [-R random_wait] [-RT random_time_range] [-W wait time] [-T threads]")
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
	fmt.Println("\t-I\n\t  iterations for the fuzzer to use")
	fmt.Printf("\n")
	fmt.Println("\t-r --request\n\t  specify path to text file which contains a request to send")
	fmt.Println("\t-R --random-wait\n\t  wait a random fraction of a second")
	fmt.Println("\t-RT --random-wait-range\n\t  define a scalar for the random wait time (ie randomWaitTime * range)")
	fmt.Println("\t-W --wait\n\t  set a static time to wait")
	fmt.Println("\t-T\n\t  number of threads (requests in air at any time) (10 by default)")
}

func init() {
	flag.Usage = printUsage
	//traceLog = log.New(os.Stderr, fmt.Sprint(traceColor("TRACE: ")), log.Ldate|log.Ltime|log.Lshortfile)
	infoLog = log.New(os.Stderr, fmt.Sprint(infoColor("INFO: ")), log.Ldate|log.Ltime)
	warnLog = log.New(os.Stderr, fmt.Sprint(warnColor("WARN: ")), log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, fmt.Sprint(errorColor("ERROR: ")), log.Ldate|log.Ltime)

	const (
		helpUsage            = "prints usage, then exits"
		urlUsage             = "url to generate a request to"
		methodUsage          = "method to use while generating request"
		bodyUsage            = "request body to use while generating request"
		headerUsage          = "headers to include while generating request"
		requestUsage         = "request to use while fuzzing/bruteforcing"
		sentsUsage           = "filepath to file containing sentinels"
		fuzzUsage            = "set to fuzzing"
		generateUsage        = "let fuzzer generate new strings at sentinels"
		asciiUsage           = "limit fuzzer to ascii only"
		itersUsage           = "number of iterations for the fuzzer to use"
		timesUsage           = "number of times to repeat the fuzzing and sending"
		waitUsage            = "time to wait"
		randomWaitUsage      = "wait a random amount of time between requests"
		randomWaitRangeUsage = "scalar for the random wait time"
		threadsUsage         = "number of threads to use (number of requests at any given time"
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

	flag.DurationVar(&waitFlag, "W", 0, waitUsage)
	flag.DurationVar(&waitFlag, "wait", 0, waitUsage)

	flag.BoolVar(&randomWaitFlag, "R", false, randomWaitUsage)
	flag.BoolVar(&randomWaitFlag, "random-wait", false, randomWaitUsage)

	flag.IntVar(&randomWaitRange, "RT", 1, randomWaitRangeUsage)
	flag.IntVar(&randomWaitRange, "random-wait-range", 1, randomWaitRangeUsage)

	flag.IntVar(&itersFlag, "I", 3, itersUsage)
	flag.IntVar(&timesFlag, "times", 1, timesUsage)
	flag.IntVar(&threadsFlag, "T", 10, threadsUsage)
}

func main() {
	flag.Parse()
	art.DrawArt()

	var sents []string

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
	if randomWaitFlag == false && randomWaitRange != 1 {
		randomWaitFlag = true
	}

	if fuzzFlag { // fuzzing!
		if asciiFlag {
			errorLog.Println("not implemented")
			os.Exit(1)
		}
		if generateFlag {
			errorLog.Println("not implemented")
			os.Exit(1)
		}

		times := 0
		remainder := 0
		responses := make(chan *respAndStrings, timesFlag)
		donePrinting := make(chan struct{})
		var wg sync.WaitGroup

		go printLoop(timesFlag, responses, donePrinting)

		req, err := ioutil.ReadFile(requestFlag)
		if err != nil {
			errorLog.Printf("%v\n", err)
			os.Exit(1)
		}

		ret = string(req)
		// times is the number of times to do threadsFlag threads, and the remainder
		// is the number of remaining threads to do
		if timesFlag == threadsFlag {
			times = timesFlag
		} else if timesFlag > threadsFlag {
			times = int(timesFlag / threadsFlag)
			remainder = int(timesFlag % threadsFlag)
		} else {
			remainder = timesFlag
		}

		for i := 0; i < times; i++ {
			for j := 0; j < threadsFlag; j++ {
				go func(request string) {
					wg.Add(1)
					seed := time.Now().UTC().UnixNano()
					infoLog.Printf("leased thread with seed %d\n", seed)
					var fuzzed []string
					for i := 0; i < len(sents); i++ {
						var err error
						thisRet, fuzz, err := fuzzer.MutateSelection(request, sents[i], seed, itersFlag)
						if err != nil {
							break
						}
						fuzzed = append(fuzzed, fuzz)
						request = thisRet
					}
					resp, err := doRequest(request)
					if err != nil {
						errorLog.Printf("%v\n", err)
						responses <- nil
					}
					ret := respAndStrings{resp: resp, strings: fuzzed}
					responses <- &ret
					wg.Done()
				}(ret)
			}
			infoLog.Printf("leased %d routines, waiting for them to return", threadsFlag)
			wg.Wait()
		}

		for i := 0; i < remainder; i++ {
			go func(request string) {
				wg.Add(1)
				seed := time.Now().UTC().UnixNano()
				infoLog.Printf("leased thread with seed %d\n", seed)
				var fuzzed []string
				for i := 0; i < len(sents); i++ {
					var err error
					thisRet, fuzz, err := fuzzer.MutateSelection(request, sents[i], seed, itersFlag)
					if err != nil {
						break
					}
					fuzzed = append(fuzzed, fuzz)
					request = thisRet
				}
				resp, err := doRequest(request)
				if err != nil {
					errorLog.Printf("%v\n", err)
					responses <- nil
				}
				ret := respAndStrings{resp: resp, strings: fuzzed}
				responses <- &ret
				wg.Done()
			}(ret)
		}

		infoLog.Println("waiting for printing to finish")
		<-donePrinting
		infoLog.Println("shutting down")
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

		count, err := util.CountPerms(fps)
		if count == 0 {
			errorLog.Println("found an empty file?")
			printUsage()
			os.Exit(1)
		}
		if err != nil {
			errorLog.Printf("%v\n", err)
			os.Exit(1)
		}
		perms := make(chan []string, count)
		responses := make(chan *respAndStrings, count)
		done := make(chan struct{})
		donePrinting := make(chan struct{})
		var wg sync.WaitGroup

		go fillPerms(perms, fps, offs, done)
		go printLoop(count, responses, donePrinting)

		// wait until the perms channel is filled
		<-done
		for len(perms) != 0 { // while there's still info to read from perms
			for i := 0; i < threadsFlag; i++ { // lease threadsFlag threads
				go func() { //each of those threads do this
					// add one to the waitgroup
					wg.Add(1)
					//check to see if we're going to sleep before we send
					if randomWaitFlag {
						scale := rand.Intn(100) / 100
						infoLog.Printf("Sleeping for %d seconds\n", scale)
						time.Sleep(time.Duration(time.Duration(scale) * time.Second))
					}
					if waitFlag != time.Duration(0) {
						time.Sleep(waitFlag)
					}
					// get a permutation from the channel
					bruteThis, ok := <-perms
					if ok { // if we were able to read one
						// call the brute forcer to get the new request
						do, err := brute.Forcer(ret, bruteThis, sents)
						if err != nil {
							responses <- nil
							errorLog.Printf("%v", err)
						}
						// go and do that response (this is SLOW)
						resp, err := doRequest(do)
						if err != nil { // if we can't do the response
							errorLog.Printf("%v", err)
							// we need to give the printer channel SOMETHING
							responses <- nil
						} else { // give this to the printer channel
							ret := respAndStrings{resp: resp, strings: bruteThis}
							responses <- &ret
						}
					}
					// mark this thread done
					wg.Done()
				}()
			}
			// main waits here for all those threads to return
			infoLog.Printf("leased %d routines, waiting for them to return", threadsFlag)
			wg.Wait()
		}
		// main waits here until there's nothing more to print
		<-donePrinting
		infoLog.Println("shutting down")
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

func fillPerms(perms chan []string, fps []*os.File, offs []int64, done chan struct{}) {
	var exhausted = false
	var out []string
	var err error
	out, exhausted, err = util.ReadInputs(fps, offs)
	for !exhausted {
		if err != nil {
			return
		}
		perms <- out
		out, exhausted, err = util.ReadInputs(fps, offs)
	}
	close(perms)
	close(done)
}

func printLoop(n int, printChan chan *respAndStrings, doneChan chan struct{}) {
	for i := 0; i < n; i++ {
		m, ok := <-printChan
		if ok {
			if m != nil {
				resStr, err := decoder.Decode(m.resp)
				if err != nil {
					errorLog.Println("error decoding response")
				}
				fmt.Println("#####################################")
				fmt.Printf("SENT\n%q\nRESPONSE:\n", m.strings)
				fmt.Println("#####################################")
				dump, err := httputil.DumpResponse(m.resp, false)
				if err != nil {
					errorLog.Printf("error dumping response header: %v", err)
				}
				fmt.Printf("%s", dump)

				fmt.Printf("%s\n", resStr)
			}
		}
	}
	close(doneChan)
}
