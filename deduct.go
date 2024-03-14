package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	adjustedSuffix = "-adjusted"
	leftoverSuffix = "-leftover"
)

func deduct(csv1Path, csv2Path string, top2 int) error {
	payouts, err := readPayoutsCSV(csv1Path)
	if err != nil {
		return err
	}

	overpaid, err := readPayoutsCSV(csv2Path)
	if err != nil {
		return err
	}
	if top2 != 0 {
		overpaid, _ = topPayouts(overpaid, top2)
	}

	sumBefore, _ := statsFIL(payouts)
	var leftover, totalAdjusted float64
	var applied int

	for addr, over := range overpaid {
		payout, ok := payouts[addr]
		if !ok {
			fmt.Println("Overpaid address", addr, "not found in payouts")
			continue
		}
		fmt.Println("Reducing current payout for", addr)
		fmt.Println("  payout before deduction:", payout.fil)
		fmt.Println("  overpayment amount:", over.fil)

		if over.fil > payout.fil {
			over.fil -= payout.fil
			totalAdjusted += payout.fil
			payout.fil = 0
			delete(payouts, addr)
			leftover += over.fil
		} else if over.fil < payout.fil {
			payout.fil -= over.fil
			totalAdjusted += over.fil
			over.fil = 0
			delete(overpaid, addr)
		} else {
			totalAdjusted += over.fil
			delete(payouts, addr)
			delete(overpaid, addr)
			payout.fil = 0
			over.fil = 0
		}
		fmt.Println("  payout after deduction:", payout.fil)
		fmt.Println("  remaining overpayment:", over.fil)
		fmt.Println()
		applied++
	}
	sumAfter, _ := statsFIL(payouts)
	fmt.Println("--------------------------------")
	fmt.Println("Applied", applied, "deductions from f2 to payouts in f1")
	fmt.Println("Total payouts before deductions:", sumBefore, "FIL")
	fmt.Println("Total payout adjustment:        ", -totalAdjusted, "FIL")
	fmt.Println("Total payouts after deductions: ", sumAfter, "FIL")
	fmt.Println("Total leftover overpayments:    ", leftover, "FIL")

	fmt.Println()

	adjustedName := strings.TrimSuffix(csv1Path, ".csv") + adjustedSuffix + ".csv"
	err = writePayoutsCSV(adjustedName, payouts)
	if err != nil {
		return fmt.Errorf("failed to write adjusted payouts csv file: %w", err)
	}
	fmt.Println("Wrote adjusted payouts to:", adjustedName)

	if len(overpaid) != 0 {
		leftoverName := strings.TrimSuffix(csv2Path, ".csv") + leftoverSuffix + ".csv"
		err = writePayoutsCSV(leftoverName, overpaid)
		if err != nil {
			return fmt.Errorf("failed to write loetover overpayments csv file: %w", err)
		}
		fmt.Println("Wrote leftover overpayments to:", leftoverName)
	}

	return nil
}

func writePayoutsCSV(filePath string, records map[string]*record) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("opening csv file for writing: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	err = w.Write([]string{"Recipient", "FIL", "Method", "Params"})
	if err != nil {
		return err
	}
	recStrs := make([]string, 4)
	for addr, rec := range records {
		recStrs[0] = addr
		recStrs[1] = strconv.FormatFloat(rec.fil, 'f', -1, 64)
		recStrs[2] = rec.method
		recStrs[3] = rec.params
		if err = w.Write(recStrs); err != nil {
			return fmt.Errorf("writing record: %w", err)
		}
	}
	return nil
}
