# frute
an http fuzzer / brute forcer

## requires:

`go get github.com/kortschak/ct`

## example

go run main.go -r testfiles/fuzzme3 testfiles/testFile0 testfiles/testFile1
go run main.go -r testfiles/fuzzme2 -f -times 2
## usage
```
  usage:
      frute [-h] [--url url --method method [--body body] [-header header_pair]] request_output_file
  or  frute -f -r request_file [-s sentinel_file] [-times number_of_times
  or  frute -r request_file [-s sentinel_file] [-R random_wait] [-RT random_time_range] [-W wait time] [-T threads]

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

    -r --request
      specify path to text file which contains a request to send
    -s --sentinel
      specify path to text file which contains a list of sentinels to replace while bruting/fuzzing
    -f
      fuzz string at sentinels
    -times
      times to repeat the fuzzing and sending

    -r --request
      specify path to text file which contains a request to send
    -R --random-wait
      wait a random fraction of a second
    -RT --random-wait-range
      define a scalar for the random wait time (ie randomWaitTime * range)
    -W --wait
      set a static time to wait
    -T
      number of threads (requests in air at any time) (10 by default)
```
## notes
### Threads?
Frute is written in go, and uses goroutines to handle sending requests, printing, and filling a channel with the strings to use while brute forcing.  Because of that, frute defaults to using only 10 requests "in flight" at once.  This can be changed (see usage)
#### Performance concerns
In order to speed up execution, frute prefills a channel with all the possible combinations of the provided files.  Mathematically, what this means is that there is a channel of strings which contains the number of lines in your *multiplied together*.  This value can get very large very fast, but in most cases you should only be limited by your machine's memory.
### Sentinels?
Frute uses a dead simple method of "mark up" of requests: you can define the sentinels that it searches for in the request by simply adding them to the request, then informing frute of the sentinels being used.  Good examples are strings not likely to show up in the request: two or three (or four or more!) exclamation points, commas, periods, basically anything.  It then finds the marked areas of the request, and fuzzes or replaces them with files from the brute forcer.
There's a problem to be aware of, however: frute uses regexes to find the locations to change in the request.  Because of this, you can't use any character that has meaning in a regex as a sentinel.  These characters are `"]", "^", "\", "[", ".", "(", ")", "-"`.
By default, frute uses `"~~~", "!!!", "@@@", "###", "&&&", "<<<", ">>>", "___", ",,,", "'''"`, but these can redefined as you like (see usage).