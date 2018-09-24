package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/alivingvendingmachine/frute/util"
	"github.com/kortschak/ct"
)

//f "github.com/alivingvendingmachine/frute/fuzzer"
//r "github.crm/alivingvendingmachine/frute/replacer"

var (
	//traceLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger

	traceColor = ct.Fg(ct.Cyan).Paint
	infoColor  = ct.Fg(ct.Green).Paint
	warnColor  = ct.Fg(ct.Yellow).Paint
	errorColor = ct.Fg(ct.Red).Paint

	helpFlag   bool
	urlFlag    string
	methodFlag string
	bodyFlag   string
)

func printUsage() {
	fmt.Println("usage: frute [-h] [--url url [--method method [--body body]]] request_output_file")
	fmt.Println("\t-h --help\n\tprint this usage message, then quit\n")
	fmt.Println("\t-u --url\n\tspecified url to use while generating a new request text file")
	fmt.Println("\t-m --method\n\tthe method to use on the new request")
	fmt.Println("\t-b --body\n\tthe body of the request to use")
}

func init() {
	flag.Usage = printUsage
	//traceLog = log.New(os.Stderr, fmt.Sprint(traceColor("TRACE: ")), log.Ldate|log.Ltime|log.Lshortfile)
	infoLog = log.New(os.Stderr, fmt.Sprint(infoColor("INFO: ")), log.Ldate|log.Ltime)
	warnLog = log.New(os.Stderr, fmt.Sprint(warnColor("WARN: ")), log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, fmt.Sprint(errorColor("ERROR: ")), log.Ldate|log.Ltime)

	const (
		helpUsage   = "prints usage, then exits"
		urlUsage    = "url to generate a request to"
		methodUsage = "method to use while generating request"
		bodyUsage   = "request body to use while generating request"
	)
	flag.BoolVar(&helpFlag, "h", false, helpUsage+" (shorthand)")
	flag.BoolVar(&helpFlag, "help", false, helpUsage)

	flag.StringVar(&urlFlag, "u", "", urlUsage+" (shorthand)")
	flag.StringVar(&urlFlag, "url", "", urlUsage)

	flag.StringVar(&methodFlag, "m", "", methodUsage+" (shorthand)")
	flag.StringVar(&methodFlag, "method", "", methodUsage)

	flag.StringVar(&bodyFlag, "b", "", bodyUsage+" (shorthand)")
	flag.StringVar(&bodyFlag, "body", "", bodyUsage)
}

func main() {
	flag.Parse()

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if len(flag.Args()) == 0 && urlFlag == "" {
		flag.Usage()
		os.Exit(0)
	}
	if urlFlag != "" && methodFlag == "" {
		errorLog.Println("method cannot be blank!\n")
		flag.Usage()
		os.Exit(1)
	}
	if urlFlag != "" && len(flag.Args()) == 0 {
		errorLog.Println("no output file for request\n")
		flag.Usage()
		os.Exit(1)
	}
	if urlFlag != "" && bodyFlag == "" {
		warnLog.Println("request body is blank")
	}

	if urlFlag != "" {
		infoLog.Println("generating request")
		err := util.GenerateRequest(methodFlag, urlFlag, bodyFlag, flag.Args()[0])
		if err != nil {
			errorLog.Printf("%v\n", err)
		}
		os.Exit(0)
	}
}
