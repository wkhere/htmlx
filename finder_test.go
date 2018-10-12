package htmlx

import (
	"os"
	"path"
	"testing"

	p "github.com/wkhere/htmlx/pred"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestEmpty(t *testing.T) {
	var empty Finder

	if !empty.IsEmpty() {
		t.Errorf("expected finder to be empty")
	}
	if empty.Find(p.ID("whatever")) != empty {
		t.Errorf("expected empty.Find to return also empty finder")
	}
	if empty.FindSibling(p.ID("whatever")) != empty {
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
	div0 := top.Find(p.Element(atom.Div))

	if div0.Find(p.Class("bar")) != div0.Find(p.Class("other")) {
		t.Errorf("mismatch")
	}
}

func TestFromString(t *testing.T) {
	top, _ := FinderFromString(`<div id="1"></div>`)

	div := top.Find(p.Element(atom.Div))

	if res := top.Find(p.ID("1")); res != div {
		t.Errorf("mismatch")
	}
}

func TestFind(t *testing.T) {
	top, _ := FinderFromData(testdata("simple.html"))

	id1 := top.Find(p.ID("1"))

	span1 := id1.Find(p.Element(atom.Span))
	span1text := span1.FirstChild().String()
	if s := "1st"; span1text != s {
		t.Errorf("mismatch:\ngot `%s`\nexp `%s`", span1text, s)
	}

	span2 := id1.Find(p.Class("bar"))
	span2text := span2.FirstChild().String()
	if s := "2nd"; span2text != s {
		t.Errorf("mismatch:\ngot `%s`\nexp `%s`", span2text, s)
	}

	if res := span1.FindSibling(p.Element(atom.Span)); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span2)
	}

	if res := span1.FindSibling(p.Class("bar")); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`",
			res.String(), span2.String())
	}

	if res := span1.FindSibling(p.ID("2")); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span2)
	}

	if res := id1.Find(p.Attr("attr2", "boom")); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span2)
	}

	if res := span1.FindSibling(p.Attr("attr2", "boom")); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span2)
	}

	span3 := span2.FindSibling(p.Element(atom.Span))
	span3text := span3.FirstChild().String()
	if s := "3rd"; span3text != s {
		t.Errorf("mismatch:\ngot `%s`\nexp `%s`", span3text, s)
	}

	if res := span2.FindSibling(p.Class("bar")); res != span3 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span3)
	}

	spanInner := id1.Find(p.Class("xyz"))
	spanInnerText := spanInner.FirstChild().String()
	if s := "inner"; spanInnerText != s {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", spanInnerText, s)
	}

	div4 := span3.FindSibling(p.Element(atom.Div))

	if res := div4.Find(p.Class("xyz")); res != spanInner {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, spanInner)
	}

	badFinds := []Finder{
		top.Find(p.ID("bad")),
		id1.Find(p.Class("bad")),
		id1.NextSibling().FindSibling(p.ID("any")),
		id1.FirstChild().FindSibling(p.ID("bad")),
		span1.FindSibling(p.Class("xyz")),
		div4.Find(p.Class("another")),
	}

	for i, bad := range badFinds {
		if !bad.IsEmpty() {
			t.Errorf("mismatch in tc[%d]:\ngot `%v`\nexp empty find", i, bad)
		}
	}
}

func TestDepthFind(t *testing.T) {
	top, _ := FinderFromString(`
		<div id="1">
		  <span id="inner"></span>
		</div>
		<span id="2"></span>
	`)

	e := top.Find(p.Element(atom.Span))
	if v, _ := e.Attr().ID(); v != "inner" {
		t.Errorf("expected to find inner element, got: id=`%v`", v)
	}
}

func TestAttrShortcuts(t *testing.T) {
	top, _ := FinderFromData(testdata("simple.html"))

	e := top.Find(p.ID("2"))

	if v, _ := e.Attr().Val("attr2"); v != "boom" {
		t.Errorf("mismatch:\ngot `%v`\nexp `boom`", v)
	}
	if _, ok := e.Attr().Val("attr_nonexistent"); ok {
		t.Errorf("mismatch:\ngot ok=true\nexp ok=false")
	}
	if v, _ := e.Attr().ID(); v != "2" {
		t.Errorf("mismatch:\ngot id=`%v`\nexp id=`2`", v)
	}
	if cs, _ := e.Attr().ClassList(); cs[0] != "bar" {
		t.Errorf("e.ClassList()[0] should be `bar`")
	}
	if !e.Attr().Exists("attr2") {
		t.Errorf(`e.HasAttr("attr2") should be true`)
	}
	if e.Attr().Exists("attr_nonexistent") {
		t.Errorf(`e.HasAttr("attr_nonexistent") should be false`)
	}
	if !e.Attr().HasVal("attr2", "boom") {
		t.Errorf(`e.HasAttrVal("attr2", "boom") should be true`)
	}
	if !e.Attr().HasID("2") {
		t.Errorf(`e.HasId("2") should be true`)
	}
	if e.Attr().HasID("bad_id") {
		t.Errorf(`e.HasId("bad_id") should be false`)
	}
	if !e.Attr().HasClass("bar") {
		t.Errorf(`e.HasClass("bar") should be true`)
	}
	if e.Attr().HasClass("class_nonexistent") {
		t.Errorf(`e.HasClass("class_nonexistent") should be false`)
	}
}

func TestEmptyFinderAttr(t *testing.T) {
	f := Finder{}

	tab := []bool{
		func() bool { _, ok := f.Attr().ID(); return ok }(),
		func() bool { _, ok := f.Attr().ClassList(); return ok }(),
		f.Attr().Exists("any"),
		f.Attr().HasID("any"),
		f.Attr().HasClass("any"),
	}

	for i, tc := range tab {
		if tc {
			t.Errorf("tc[%d] should be false", i)
		}
	}
}

// Measurements were done on `metal` machine, perf mode.
// Adding closures raises execution time from 210 ns/op to 231 ns/op.
// Now with Finder struct & methods it's 328ns/op.
// Todo: discover why it got slower.
// Maybe (f Finder) -> (f *Finder) ?
func BenchmarkBasic(b *testing.B) {
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
		id1 := f.Find(p.ID("id1"))
		id1.Find(p.Attr("class", "baz"))
	}
}

func BenchmarkGoV(b *testing.B) {
	f, _ := FinderFromData(testdata("gatesofvienna.html"))
	for n := 0; n < b.N; n++ {
		f.Find(p.Class("html-end-of-file"))
	}
}

func testdata(filename string) (f *os.File) {
	f, err := os.Open(path.Join("testdata", filename))
	if err != nil {
		panic(err)
	}
	return
}
