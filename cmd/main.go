package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/rsp9u/go-xlsshape/oxml"
	"github.com/rsp9u/seq2xls"
	"github.com/rsp9u/seq2xls/seqdiag"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	d := seqdiag.ParseSeqdiag(b)

	ss := oxml.NewSpreadsheet()
	lls, err := seqdiag.ExtractLifelines(d)
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := seqdiag.ExtractMessages(d, lls)
	if err != nil {
		log.Fatal(err)
	}

	seq2xls.DrawLifelines(ss, lls, len(msgs))
	seq2xls.DrawMessages(ss, msgs)
	ss.Dump("example1.xlsx")
}
