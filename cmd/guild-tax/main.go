// Command guild-tax 跑公会税收样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/tax"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("guild-tax", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "round", "round / cap / exempt / idem / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "round":
		ledger := tax.New(30, 1000)
		charged, _ := ledger.Earn("p1", 150, "r1")
		fmt.Fprintf(stdout, "tax=%d\n", charged)
	case "cap":
		ledger := tax.New(30, 10)
		ledger.Earn("p1", 100, "r1")
		charged, _ := ledger.Earn("p1", 100, "r2")
		fmt.Fprintf(stdout, "second=%d\n", charged)
	case "exempt":
		ledger := tax.New(30, 1000)
		ledger.Exempt("p1")
		charged, _ := ledger.Earn("p1", 100, "r1")
		fmt.Fprintf(stdout, "tax=%d\n", charged)
	case "idem":
		ledger := tax.New(30, 1000)
		ledger.Earn("p1", 100, "same")
		ledger.Earn("p1", 100, "same")
		fmt.Fprintf(stdout, "guild=%d\n", ledger.Guild())
	case "scan":
		ledger := tax.New(30, 1000000)
		for index := 0; index < 3000; index++ {
			ledger.Earn(fmt.Sprintf("p-%d", index), 100, fmt.Sprintf("r-%d", index))
		}
		for index := 0; index < 3000; index++ {
			ledger.Today("p-2999")
		}
		fmt.Fprintf(stdout, "scanned=%d\n", ledger.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample round|cap|exempt|idem|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
