package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var (
		csv1Path   string
		csv2Path   string
		top1, top2 int
	)
	flag.StringVar(&csv1Path, "f1", "", "first payouts csv file")
	flag.StringVar(&csv2Path, "f2", "", "second payouts csv file")
	flag.IntVar(&top1, "top1", 0, "limit file 1 to N records with highest FIL")
	flag.IntVar(&top2, "top2", 0, "limit file 2 to N records with highest FIL")
	flag.Parse()

	if csv1Path == "" {
		fmt.Fprintln(os.Stderr, "missing value for -f1")
		os.Exit(1)
	}

	if csv1Path == "" {
		fmt.Fprintln(os.Stderr, "missing value for -f2")
		os.Exit(1)
	}

	err := compare(csv1Path, csv2Path, top1, top2)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func compare(csv1Path, csv2Path string, top1, top2 int) error {
	payouts1, err := readPayoutsCSV(csv1Path)
	if err != nil {
		return err
	}
	var sorted1 []string
	payouts1, sorted1 = topPayouts(payouts1, top1)

	payouts2, err := readPayoutsCSV(csv2Path)
	if err != nil {
		return err
	}
	var sorted2 []string
	payouts2, sorted2 = topPayouts(payouts2, top2)

	var setDiff1, setDiff2 int
	var setInter int

	for addr := range payouts1 {
		_, ok := payouts2[addr]
		if ok {
			setInter++
		} else {
			setDiff1++
		}
	}
	for addr := range payouts2 {
		_, ok := payouts1[addr]
		if !ok {
			setDiff2++
		}
	}

	fmt.Println("Payouts 1 stats:")
	sum1, mean1 := statsFIL(payouts1)
	fmt.Println("  Filename:", filepath.Base(csv1Path))
	fmt.Println("  Number of payouts:", len(payouts1))
	fmt.Println("  Highest FIL:", payouts1[sorted1[0]])
	fmt.Println("  Average FIL:", mean1)
	fmt.Println("  Total FIL:", sum1)
	fmt.Println("  Payouts in file 1 only:", setDiff1)
	fmt.Println()

	fmt.Println("Payouts 2 stats:")
	sum2, mean2 := statsFIL(payouts2)
	fmt.Println("  Filename:", filepath.Base(csv2Path))
	fmt.Println("  Number of payouts:", len(payouts2))
	fmt.Println("  Highest FIL:", payouts2[sorted2[0]])
	fmt.Println("  Average FIL:", mean2)
	fmt.Println("  Total FIL:", sum2)
	fmt.Println("  Payouts in file 2 only:", setDiff2)
	fmt.Println()

	fmt.Println("Payouts in both files: ", setInter)
	return nil
}

func printPayouts(records map[string]float64) {
	for addr, fil := range records {
		fmt.Println("Addr:", addr, "FIL:", fil)
	}
}

func readPayoutsCSV(fileName string) (map[string]float64, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	recs := make(map[string]float64)
	rdr := csv.NewReader(f)

	for {
		record, err := rdr.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		addr := record[0]
		if !strings.HasPrefix(addr, "f1") {
			continue
		}
		fil, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse fil value %q: %s", record[1], err)
		}
		recs[record[0]] = fil
	}

	return recs, nil
}

func topPayouts(recs map[string]float64, top int) (map[string]float64, []string) {
	sorted := sortFIL(recs)
	if top > 0 && len(recs) > top {
		topRecs := make(map[string]float64, top)
		sorted = sorted[:top]
		for _, addr := range sorted {
			topRecs[addr] = recs[addr]
		}
		recs = topRecs
	}

	return recs, sorted
}

func sortFIL(recs map[string]float64) []string {
	type kv struct {
		Key string
		Val float64
	}

	sorted := make([]kv, len(recs))
	var i int
	for k, v := range recs {
		sorted[i] = kv{k, v}
		i++
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Val > sorted[j].Val
	})

	keys := make([]string, len(sorted))
	for i := range sorted {
		keys[i] = sorted[i].Key
	}

	return keys
}

func statsFIL(recs map[string]float64) (sum float64, mean float64) {
	for _, v := range recs {
		sum += v
	}
	mean = sum / float64(len(recs))
	return
}
