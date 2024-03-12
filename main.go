package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		csv1Path   string
		csv2Path   string
		opDeduct   bool
		top1, top2 int
	)
	flag.BoolVar(&opDeduct, "deduct", false, "Deduce overpayments in file2 from payouts in file1")
	flag.StringVar(&csv1Path, "f1", "", "first payouts csv file")
	flag.StringVar(&csv2Path, "f2", "", "second payouts csv file")
	flag.IntVar(&top1, "top1", 0, "limit file 1 to N records with highest FIL")
	flag.IntVar(&top2, "top2", 0, "limit file 2 to N records with highest FIL")
	flag.Parse()

	if csv1Path == "" {
		fmt.Fprintln(os.Stderr, "missing value for -f1")
		os.Exit(1)
	}

	if csv2Path == "" {
		fmt.Fprintln(os.Stderr, "missing value for -f2")
		os.Exit(1)
	}

	var err error
	if opDeduct {
		if top1 != 0 || top2 != 0 {
			fmt.Fprintln(os.Stderr, "-top1 and -top2 are not available with -deduct")
			os.Exit(1)
		}
		err = deduct(csv1Path, csv2Path)
	} else {
		err = compare(csv1Path, csv2Path, top1, top2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
