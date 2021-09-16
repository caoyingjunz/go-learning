package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	h    bool
	s    string
	q    *string
	port *int
)

func init() {
	flag.BoolVar(&h, "h", false, "print the helps")
	flag.StringVar(&s, "s", "/ets/ss", "sss")
	q = flag.String("q", "9999", "suppress non-error messages during configuration testing")
	port = flag.Int("p", 8080, "server Port")

	// 改变默认的 Usage
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `caoyingjun version: v1
Usage: nginx [-hvVtTq] [-s signal] [-c filename] [-p prefix] [-g directives]
Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	fmt.Println(s)
	fmt.Println(*q)
}
