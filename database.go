package greatspacerace

import (
	"bufio"
	"fmt"
	"github.com/TSavo/go.firebase"
	"os"
)

func init() {
	f, err := os.Open("firebase.secret")
	if err != nil {
		fmt.Printf("error opening firebase.secret: %v\n", err)
	}
	r := bufio.NewReader(f)
	url, e := Readln(r)
	if e != nil {
		panic(e)
	}
	secret, e := Readln(r)
	if e != nil {
		panic(e)
	}
	db = firebase.New(url, secret)
}

var db *firebase.FirebaseRoot;

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
