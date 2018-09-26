# frute
an http fuzzer / brute forcer

## requires:

`go get github.com/kortschak/ct`

## usage

usage:
    frute [-h] [--url url --method method [--body body] [-header header_pair]] request_output_file
or  frute -r request_file ...

    -h --help
      print this usage message, then quit

    -u --url
      specified url to use while generating a new request text file
    -m --method
      the method to use on the new request
    -b --body
      the body of the request to use
    -H --header
      headers to include while generating a request

    -r
      specify path to text file which contains a request to send
    -sent --sentinels

## notes
### Sentinels?
Frute uses a dead simple method of "mark up" of requests: you can define the sentinels that it searches for in the request by simply adding them to the request, then informing frute of the sentinels being used.  Good examples are strings not likely to show up in the request: two or three (or four or more!) exclamation points, commas, periods, basically anything.  It then finds the marked areas of the request, and fuzzes or replaces them with files from the brute forcer.