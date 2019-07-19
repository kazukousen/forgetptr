package forgetptr_test

import (
	"testing"

	"github.com/kazukousen/forgetptr/analysis/passes/forgetptr"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, forgetptr.Analyzer, "a")
}
