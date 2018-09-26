package requester

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

// ReadReqAndDo reads a request in from a file path, then performs the request,
// returning any errors
func ReadReqAndDo(c chan *http.Response, filepath string) error {
	fmt.Println("open file")
	fp, err := os.Open(filepath)
	if err != nil {
		c <- nil
		return err
	}
	defer fp.Close()

	fmt.Println("new reader")
	reader := bufio.NewReader(fp)
	req, err := http.ReadRequest(reader)
	if err != nil {
		c <- nil
		return err
	}

	fmt.Println("client do")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c <- nil
		return err
	}

	fmt.Println("pack channel")
	c <- resp
	return nil
}
