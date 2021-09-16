package main

import (
	"flag"
	"fmt"

	klog "k8s.io/klog/v2"

	"golang-learning/practise/klog-practise/test"
)

// https://github.com/kubernetes/klog/blob/master/examples/coexist_klog_v1_and_v2/coexist_klog_v1_and_v2.go

func main() {
	// initialize klog/v2, can also bind to a local flagset if desired
	klog.InitFlags(nil)

	//flag.Set("logtostderr", "false") // By default klog logs to stderr, switch that off
	//flag.Set("alsologtostderr", "false") // false is default, but this is informative
	//flag.Set("stderrthreshold", "FATAL") // stderrthreshold defaults to ERROR, we don't want anything in stderr
	//flag.Set("log_file", "test.log") // log to a file

	// parse klog/v2 flags
	flag.Parse()
	// make sure we flush before exiting
	defer klog.Flush()

	test.Test1()
	klog.Warningf("aaa %v", "dddd")
	klog.InfoS("hello from klog (v2)!", "aaa", "bbb")
	klog.ErrorS(fmt.Errorf("error info"), "a", "aaaa", "bbb")
}
