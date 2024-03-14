package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type record struct {
	fil    *big.Float
	method string
	params string
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
	fmt.Println("  Average FIL:", mean1.String())
	fmt.Println("  Total FIL:", sum1.String())
	fmt.Println("  Payouts in file 1 only:", setDiff1)
	fmt.Println()

	fmt.Println("Payouts 2 stats:")
	sum2, mean2 := statsFIL(payouts2)
	fmt.Println("  Filename:", filepath.Base(csv2Path))
	fmt.Println("  Number of payouts:", len(payouts2))
	fmt.Println("  Highest FIL:", payouts2[sorted2[0]])
	fmt.Println("  Average FIL:", mean2.String())
	fmt.Println("  Total FIL:", sum2.String())
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

func readPayoutsCSV(fileName string) (map[string]*record, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	records := make(map[string]*record)
	rdr := csv.NewReader(f)

	for {
		rec, err := rdr.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		addr := rec[0]
		if len(addr) < 32 || !strings.HasPrefix(addr, "f") {
			continue
		}
		var fil big.Float
		_, _, err = fil.Parse(rec[1], 10)
		if err != nil {
			return nil, fmt.Errorf("cannot parse fil value %q: %s", rec[1], err)
		}
		records[rec[0]] = &record{
			fil:    &fil,
			method: rec[2],
			params: rec[3],
		}
	}

	return records, nil
}

func topPayouts(recs map[string]*record, top int) (map[string]*record, []string) {
	sorted := sortFIL(recs)
	if top > 0 && len(recs) > top {
		topRecs := make(map[string]*record, top)
		sorted = sorted[:top]
		for _, addr := range sorted {
			topRecs[addr] = recs[addr]
		}
		recs = topRecs
	}

	return recs, sorted
}

func sortFIL(recs map[string]*record) []string {
	type kv struct {
		Key string
		Val *big.Float
	}

	sorted := make([]kv, len(recs))
	var i int
	for k, v := range recs {
		sorted[i] = kv{k, v.fil}
		i++
	}
	slices.SortFunc(sorted, func(a, b kv) int {
		return b.Val.Cmp(a.Val)
	})

	keys := make([]string, len(sorted))
	for i := range sorted {
		keys[i] = sorted[i].Key
	}

	return keys
}

func statsFIL(recs map[string]*record) (*big.Float, *big.Float) {
	sum := new(big.Float)
	for _, v := range recs {
		sum.Add(sum, v.fil)
	}

	mean := new(big.Float).Quo(sum, new(big.Float).SetInt64(int64(len(recs))))
	return sum, mean
}
