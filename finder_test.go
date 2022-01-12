package htmlx

import (
	"bytes"
	"io"
	"os"
	"path"
	"regexp"
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

func TestFindAllPrinted(t *testing.T) {
	top, _ := FinderFromData(testdata("simple.html"))

	ff := top.FindAll(p.Element(atom.Span))

	b := new(bytes.Buffer)
	for f := range ff {
		f.Write(b)
		io.WriteString(b, "\n")
	}

	s := `<span class="foo">1st</span>` + "\n" +
		`<span id="2" class="bar other" attr2="boom">2nd</span>` + "\n" +
		`<span class="bar another">3rd</span>` + "\n" +
		`<span class="xyz yet-another">inner</span>` + "\n"

	if res := b.String(); res != s {
		t.Errorf("mismatch:\ngot:\n%v\nexp:\n%v", res, s)
	}
}

func TestFindAll(t *testing.T) {
	top, _ := FinderFromData(testdata("simple.html"))

	ff := top.FindAll(p.Element(atom.Span)).Consume()

	if len(ff) != 4 {
		t.Errorf("got %d, exp %d", len(ff), 4)
	}
	for i, f := range ff {
		if s, res := atom.Span, f.Node.DataAtom; res != s {
			t.Errorf("ff[%d]: got `%s`, exp `%s`", i, res, s)
		}
	}
	if s, res := "1st", ff[0].InnerText(); res != s {
		t.Errorf("got `%s`, exp `%s`", res, s)
	}
	if s, res := "2nd", ff[1].InnerText(); res != s {
		t.Errorf("got `%s`, exp `%s`", res, s)
	}
	if s, res := "3rd", ff[2].InnerText(); res != s {
		t.Errorf("got `%s`, exp `%s`", res, s)
	}
	if s, res := "inner", ff[3].InnerText(); res != s {
		t.Errorf("got `%s`, exp `%s`", res, s)
	}

	if e, e0 := ff[0].FindSibling(p.Element(atom.Span)), ff[1]; e != e0 {
		t.Errorf("mismatch:\ngot:\n%v\nexp:\n%v", e, e0)
	}
	if e, e0 := ff[1].FindPrevSibling(p.Element(atom.Span)), ff[0]; e != e0 {
		t.Errorf("mismatch:\ngot:\n%v\nexp:\n%v", e, e0)
	}

	if e, e0 := ff[1].FindSibling(p.Element(atom.Span)), ff[2]; e != e0 {
		t.Errorf("mismatch:\ngot:\n%v\nexp:\n%v", e, e0)
	}
	if e, e0 := ff[2].FindPrevSibling(p.Element(atom.Span)), ff[1]; e != e0 {
		t.Errorf("mismatch:\ngot:\n%v\nexp:\n%v", e, e0)
	}

	if e, e0 := ff[3].Parent(), ff[2].FindSibling(p.Element(atom.Div)); e != e0 {
		t.Errorf("mismatch:\ngot:\n%v\nexp:\n%v", e, e0)
	}
}

func TestFindAllEmpty(t *testing.T) {
	ff := Finder{}.FindAll(p.Element(atom.Div)).Consume()
	if ff != nil {
		t.Errorf("expected empty slice, got: %v", ff)
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
		t.Errorf("e.a.ClassList()[0] should be `bar`")
	}
	if !e.Attr().Exists("attr2") {
		t.Errorf(`e.a.HasAttr("attr2") should be true`)
	}
	if e.Attr().Exists("attr_nonexistent") {
		t.Errorf(`e.a.HasAttr("attr_nonexistent") should be false`)
	}
	if !e.Attr().HasVal("attr2", "boom") {
		t.Errorf(`e.a.HasAttrVal("attr2", "boom") should be true`)
	}
	if !e.Attr().HasID("2") {
		t.Errorf(`e.a.HasId("2") should be true`)
	}
	if e.Attr().HasID("bad_id") {
		t.Errorf(`e.a.HasId("bad_id") should be false`)
	}
	if !e.Attr().HasClass("bar") {
		t.Errorf(`e.a.HasClass("bar") should be true`)
	}
	if e.Attr().HasClass("class_nonexistent") {
		t.Errorf(`e.a.HasClass("class_nonexistent") should be false`)
	}

	if !e.Attr().HasClassCond(regexp.MustCompile("[oO]ther").MatchString) {
		t.Errorf(`e.a.HasClassCond(r"[oO]ther".MatchString) should be true`)
	}
}

func TestPredShortcuts(t *testing.T) {
	top, _ := FinderFromData(testdata("simple.html"))
	var r *regexp.Regexp
	var e Finder

	e = top.Find(p.InnerText("2nd"))
	if !e.Attr().HasID("2") {
		t.Errorf(`failed to find inner text "2nd"`)
	}

	r = regexp.MustCompile("3rd")
	e = top.Find(p.InnerTextCond(r.MatchString))
	if !e.Attr().HasClass("bar") {
		t.Errorf(`failed to find inner text r"3rd"`)
	}

	r = regexp.MustCompile("[oO]ther")
	e = top.Find(p.ClassCond(r.MatchString))
	if !e.Attr().HasID("2") {
		t.Errorf(`failed to find class r"[oO]ther"`)
	}

	r = regexp.MustCompile("2")
	e = top.Find(p.IDCond(r.MatchString))
	if !e.Attr().HasID("2") {
		t.Errorf(`failed to find ID r"2"`)
	}

	r = regexp.MustCompile("bo.m")
	e = top.Find(p.AttrCond("attr2", r.MatchString))
	if !e.Attr().HasID("2") {
		t.Errorf(`failed to find attr in "attr2" by r"bo.m"`)
	}

	r = regexp.MustCompile("xyz")
	e = top.Find(p.AttrWordCond("class", r.MatchString))
	if !e.Attr().HasClass("xyz") {
		t.Errorf(`failed to find attr word in "class" by r"xyz"`)
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
