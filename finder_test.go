package htmlx

import (
	"os"
	"path"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestEmpty(t *testing.T) {
	var empty Finder

	if !empty.IsEmpty() {
		t.Errorf("expected finder to be empty")
	}
	if empty.FindById("whatever") != empty {
		t.Errorf("expected empty.Find to return also empty finder")
	}
	if empty.FindSiblingById("whatever") != empty {
		t.Errorf("expected empty.FindSibling to return also empty finder")
	}
	if empty.FirstChild() != empty {
		t.Errorf("expected empty.FirstChild to return also empty finder")
	}
	if empty.NextSibling() != empty {
		t.Errorf("expected empty.NextSibling to return also empty finder")
	}
	if empty.String() != "" {
		t.Errorf("expected empty.String to return empty string")
	}
}

func TestFromNode(t *testing.T) {
	node, _ := html.Parse(testdata("simple.html"))
	top := FinderFromNode(node)
	div0 := top.FindElement(atom.Div)

	if div0.FindByClass("bar") != div0.FindByClass("other") {
		t.Errorf("mismatch")
	}
}

func TestFromString(t *testing.T) {
	top, _ := FinderFromString(`<div id="1"></div>`)

	div := top.FindElement(atom.Div)

	if res := top.FindById("1"); res != div {
		t.Errorf("mismatch")
	}
}

func TestFind(t *testing.T) {
	var s string

	top, _ := FinderFromData(testdata("simple.html"))

	id1 := top.FindById("id1")

	span1 := id1.FindElement(atom.Span)
	span1content := span1.FirstChild()
	s = "1st"
	if res := span1content.String(); res != s {
		t.Errorf("mismatch:\ngot `%s`\nexp `%s`", res, s)
	}

	span2 := id1.FindByClass("bar")

	span2content := span2.FirstChild()
	s = "2nd"
	if res := span2content.String(); res != s {
		t.Errorf("mismatch:\ngot `%s`\nexp `%s`", res, s)
	}

	if res := span1.FindSiblingElement(atom.Span); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span2)
	}

	if res := span1.FindSiblingByClass("bar"); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`",
			res.String(), span2.String())
	}

	if res := span1.FindSiblingById("id2"); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span2)
	}

	if res := id1.FindByAttr("attr2", "boom"); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span2)
	}

	if res := span1.FindSiblingByAttr("attr2", "boom"); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span2)
	}

	span3 := span2.FindSiblingElement(atom.Span)
	if res := span2.FindSiblingByClass("bar"); res != span3 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span3)
	}

	badFinds := []Finder{
		top.FindById("bad"),
		id1.FindByClass("bad"),
		id1.NextSibling().FindSiblingById("any"),
		id1.FirstChild().FindSiblingById("bad"),
		span1.FindSiblingByClass("xyz"),
	}

	for i, bad := range badFinds {
		if !bad.IsEmpty() {
			t.Errorf("mismatch in tc[%d]:\ngot `%v`\nexp empty find", i, bad)
		}
	}
}

// BenchmarkFind measurements were done on `metal` machine, perf mode.
// Adding closures raises execution time from 210 ns/op to 231 ns/op.
// Now with Finder struct & methods it's 328ns/op.
// Todo: discover why it got slower.
// Maybe (f Finder) -> (f *Finder) ?
func BenchmarkFind(b *testing.B) {
	f, _ := FinderFromString(`
		<div>
			<div id="id1">
				<span class="foo">1st</span>
				<span class="bar">2nd</span>
				<span class="baz">3rd</span>
			</div>
		</div>
	`)
	for n := 0; n < b.N; n++ {
		id1 := f.FindById("id1")
		id1.FindByAttr("class", "baz")
	}
}

func testdata(filename string) (f *os.File) {
	f, err := os.Open(path.Join("testdata", filename))
	if err != nil {
		panic(err)
	}
	return
}
