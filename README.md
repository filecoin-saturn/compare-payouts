# compare-payouts
Compare two Saturn payouts CSV files

## Install 

```shell
go install github.com/filecoin-saturn/compare-payouts@latest
```
## Run

### Usage
```shell
$ compare-payouts -help
Usage of compare-payouts:
  -f1 string
        first payouts csv file
  -f2 string
        second payouts csv file
  -top1 int
        limit file 1 to N records with highest FIL
  -top2 int
        limit file 2 to N records with highest FIL
```

### Compare payouts CSV files
```shell
compare-payouts -f1 Saturn-FVM-Payouts-2024-02.csv -f2 Saturn-FVM-Payouts-2024-01.csv 
Payouts 1 stats:
  Filename: Saturn-FVM-Payouts-2024-02.csv
  Number of payouts: <redacted>
  Highest FIL: <redacted>
  Average FIL: <redacted>
  Total FIL: <redacted>
  Payouts in file 1 only: <redacted>

Payouts 2 stats:
  Filename: Saturn-FVM-Payouts-2024-01.csv
  Number of payouts: <redacted>
  Highest FIL: <redacted>
  Average FIL: <redacted>
  Total FIL: <redacted>
  Payouts in file 2 only: <redacted>

Payouts in both files:  <redacted>
```

### Deduct overpayments from payouts
```shell
compare-payouts -deduct -f1 Saturn-FVM-Payouts-2024-03.csv -f2 Saturn-FVM-overpaid-2024-02.csv

Reducing current payout for <redacted>
  payout before deduction: 19.888213373800003
  overpayment amount: 24.084187046100002
  payout after deduction: 0
  remaining overpayment: 4.195973672299999

Reducing current payout for <redacted>
  payout before deduction: 35.134526091699996
  overpayment amount: 50.437306431100005
  payout after deduction: 0
  remaining overpayment: 15.302780339400009

Overpaid address <redacted> not found in payouts
Reducing current payout for <redacted>
  payout before deduction: 29.7999360835
  overpayment amount: 37.3769563082
  payout after deduction: 0
  remaining overpayment: 7.5770202247

...

--------------------------------
Total payouts before deductions: 29388.1022690668 FIL
Total payout adjustment:         -9507.476867757496 FIL
Total payouts after deductions:  19880.625401309302 FIL
Total leftover overpayments:     1276.6633122924998 FIL

Wrote adjusted payouts to: payouts/Saturn-FVM-Payouts-2024-03-adjusted.csv
Wrote leftover overpayments to: payouts/Saturn-FVM-overpaid-2024-02-leftover.csv
```
