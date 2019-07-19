package main

import (
	"github.com/kazukousen/forgetptr/analysis/passes/forgetptr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(
		forgetptr.Analyzer,
	)
}
