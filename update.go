package main

import (
	"os"
	"io/ioutil"
	"http"
	"fmt"
)

func CheckForUpdates() (err os.Error) {
	res, _, err := http.Get("http://github.com/mjm/Go2Lunch/raw/master/README.md")
	if res.StatusCode != 200 {
		err = os.NewError(fmt.Sprint("Response had status code ", res.StatusCode))
		return
	}
	buf, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	fmt.Printf("%s\n", buf)
	return
}

func main() {
	err := CheckForUpdates()
	if err != nil {
		panic(err)
	}
}
