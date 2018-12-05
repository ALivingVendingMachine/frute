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
  or  frute -f -r request_file [-s sentinel_file] [-times number_of_times] [-I number_of_iterations]
  or  frute -r request_file [-s sentinel_file] [-R random_wait] [-RT random_time_range] [-W wait time] [-T threads]

    -h --help
      print this usage message, then quit

    -u --url <url>
      specified url to use while generating a new request text file
    -m --method <method>
      the method to use on the new request
    -b --body <body>
      the body of the request to use
    -H --header <header>
      headers to include while generating a request

    -r --request <request_file>
      specify path to text file which contains a request to send
    -s --sentinel <sentinel_file>
      specify path to text file which contains a list of sentinels to replace while bruting/fuzzing
    -f 
      fuzz string at sentinels
    -times <number_of_times>
      times to repeat the fuzzing and sending
    -I <number_of_interations>
      times to iterate at each fuzzing location


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
### Iterations?
At every single sentinel point, frute cuts the string that's there out, and then runs it through a mutator.  It can add and subtract characters, randomly mutate them, all sorts of fun stuff.  However, every single time you do this, it can add more and more characters.  Generally speaking, more iterations means a longer string.  So if you're fuzzing, AND testing for buffer overflows, numbers of iterators larger than 10 are advisable.  For example, with a seed of 1234, "hello world" (where only world is fuzzed) becomes:

1 iteration: "hello o⟭d"
3 iterations: "hello o⟭doldo갶d"
5 iterations: "hello o⟭doldo갶d∻l푌ᢪ睁d"
25 iterations: "hello o⟭doldo갶d∻l푌ᢪ睁dol癈d骲l�䩠\uf26ado\ue14cdoldo\u0560dold첛ldoldoldᦉldol돷oldoldoldor\uf576o舀᭤dolﾆoldwrd"
250 iterations: "hello o⟭doldo갶d∻l푌ᢪ睁dol癈d骲l�䩠\uf26ado\ue14cdoldo\u0560dold첛ldoldoldᦉldol돷oldoldoldor\uf576o舀᭤dolﾆoldwrdol纒揟\uf5a2dol苉oldwrdwᦱr閧olꉛdol恿緑ldol\uf142㛞l̓do\ua632樝doldo뮊ld噭ቼ黁d㨦ldᗊr╬oldoldwrdol槣ol齉矶l岬dꅏl统oldold鼓ldw⬨dol�o◠d酣ldo躞doldolⷃol�d䐆ld⨢⼖d\ue79fr糘酩ldwr뎪\uee81驈d끓鿔doldold牫ldᧀ巤do簰d駎l깤doldo\ue0aad患ⓜ덜old娠ldoldo恠do\uf8ecd�␢doldoאָ郀ｩldo䥔d癥狛恀or쳙㕭숦dw硔rӘoldo\uf894�d힡ꢨldoldold븇齬doldo椌dol뺈姲鷘doldoldo꘤doldo쓟l突oldol䬢Ưldoldoláol❽졆쪽doldo⒔ldo\uf3efdoldュl갢old쭸덡ldoldold\u0c5bldoldoldo삽dol\ue202old玐ld쟣ldo쭅do盫깊dolゥoldordoldol섚dᑸldol塜\ue396rdold쁶꽿d㧚ldoldo酣릍〓l\uebbbo匛d䤥ldoꌑdoˡdol\uf087䱭l\ua6f9o镵doldo\ueab2뒒dol\ue948䱟ldo咚dퟴld᮱l뷻dol〲톎⟙doldo\uf13f嫫״rdoldo괟洸do㨎doldoldoldoldol꺺dol\ue3c7dol\uecf7o\ue9ead掎ldolᯒdold＄l樳doldoldoldoldol鳔dol쩣old霡嚆浌oldoldol풽old\ue077ldo캨椕doldﺑldold㍭ldol撿\u175cldo復d魷lⷦ\uf7b0oldo兣ldold\uf87e\uf828ldo‡倶緂ldord�ldwr浸긪l愕dol䱰歮l些oldoꦉd9ldol룵ꋗ灏rdold첫ldold횴ld厄ldor艍툟l喙wrdoldo隑裶ol\uf6f0d瑜ldol⦧ol딶Ꚑ脴\ue074䊗oldoldoldoldo嵻doldoldşldoldoldold쳨壄doldol䀂oldoꠥd\ue235ld"