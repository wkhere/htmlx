package htmlx

import (
	"io/ioutil"
	"testing"

	"github.com/wkhere/htmlx/pp"
)

func BenchmarkPPSimple(b *testing.B) {
	benchmarkPPFile(b, "simple.html")
}

func BenchmarkPPGov(b *testing.B) {
	benchmarkPPFile(b, "gatesofvienna.html")
}

func benchmarkPPFile(b *testing.B, file string) {
	b.Helper()
	p := pp.Printer{CompactSpaces: true, TrimEmptyAttr: true}
	f, _ := FinderFromData(testdata(file))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		p.Print(ioutil.Discard, f.Node)
	}
}
