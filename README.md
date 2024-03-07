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

