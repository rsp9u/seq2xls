package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/rsp9u/go-xlsshape/oxml"
	"github.com/rsp9u/seq2xls"
	"github.com/rsp9u/seq2xls/seqdiag"
	"github.com/rsp9u/seq2xls/seqdiag/convertor"
)

func main() {
	switch runtime.GOOS {
	case "windows":
		runOnWindows()
	default:
		runOnLinux()
	}
}

func runOnLinux() {
	var inpath, outpath string
	flag.StringVar(&inpath, "i", "-", "input file path")
	flag.StringVar(&outpath, "o", "", "output file path")
	flag.Parse()
	if outpath == "" {
		fmt.Printf("missing output file path\n\n")
		flag.Usage()
		os.Exit(1)
	}
	convert(inpath, outpath)
}

func runOnWindows() {
	flag.Parse()
	for _, inpath := range flag.Args() {
		ext := filepath.Ext(inpath)
		outpath := inpath[0:len(inpath)-len(ext)] + ".xlsx"
		convert(inpath, outpath)
	}
}

func convert(inpath, outpath string) {
	var (
		b   []byte
		err error
	)

	if inpath == "-" {
		stdin := bufio.NewScanner(os.Stdin)
		buf := new(bytes.Buffer)
		buf.Grow(1024)
		// read until eof
		for stdin.Scan() {
			buf.Write(stdin.Bytes())
			buf.WriteString("\n")
		}
		b = buf.Bytes()
	} else {
		b, err = ioutil.ReadFile(inpath)
		if err != nil {
			log.Fatal(err)
		}
	}
	d := seqdiag.ParseSeqdiag(b)

	ss := oxml.NewSpreadsheet()
	seq, err := convertor.AstToModel(d)
	if err != nil {
		log.Fatal(err)
	}

	seq2xls.DrawLifelines(ss, seq.Lifelines, len(seq.Messages))
	seq2xls.DrawMessages(ss, seq.Messages)
	seq2xls.DrawNotes(ss, seq.Notes)
	ss.Dump(outpath)
}
