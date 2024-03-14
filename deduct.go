package main

import (
	"encoding/csv"
	"fmt"
	"math/big"
	"os"
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
	sumBefore, _ := statsFIL(payouts)

	overpaid, err := readPayoutsCSV(csv2Path)
	if err != nil {
		return err
	}
	if top2 != 0 {
		overpaid, _ = topPayouts(overpaid, top2)
	}

	var leftover, totalAdjusted big.Float
	var applied int

	for addr, over := range overpaid {
		payout, ok := payouts[addr]
		if !ok {
			fmt.Println("Overpaid address", addr, "not found in payouts")
			continue
		}
		fmt.Println("Reducing current payout for", addr)
		fmt.Println("  payout before deduction:", payout.fil.String())
		fmt.Println("  overpayment amount:", over.fil.String())

		switch over.fil.Cmp(payout.fil) {
		case 1:
			over.fil.Sub(over.fil, payout.fil)
			totalAdjusted.Add(&totalAdjusted, payout.fil)
			payout.fil = new(big.Float)
			delete(payouts, addr)
			leftover.Add(&leftover, over.fil)
		case -1:
			payout.fil.Sub(payout.fil, over.fil)
			totalAdjusted.Add(&totalAdjusted, over.fil)
			over.fil = new(big.Float)
			delete(overpaid, addr)
		case 0:
			totalAdjusted.Add(&totalAdjusted, over.fil)
			delete(payouts, addr)
			delete(overpaid, addr)
			payout.fil = new(big.Float)
			over.fil = new(big.Float)
		}
		fmt.Println("  payout after deduction:", payout.fil)
		fmt.Println("  remaining overpayment:", over.fil)
		fmt.Println()
		applied++
	}
	sumAfter, _ := statsFIL(payouts)
	fmt.Println("--------------------------------")
	fmt.Println("Applied", applied, "deductions from f2 to payouts in f1")
	fmt.Println("Total payouts before deductions:", sumBefore.String(), "FIL")
	fmt.Println("Total payout adjustment:        ", totalAdjusted.Neg(&totalAdjusted).String(), "FIL")
	fmt.Println("Total payouts after deductions: ", sumAfter.String(), "FIL")
	fmt.Println("Total leftover overpayments:    ", leftover.String(), "FIL")

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
		recStrs[1] = rec.fil.String()
		recStrs[2] = rec.method
		recStrs[3] = rec.params
		if err = w.Write(recStrs); err != nil {
			return fmt.Errorf("writing record: %w", err)
		}
	}
	return nil
}
