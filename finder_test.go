package htmlx

import (
	"bytes"
	"io"
	"os"
	"path"
	"reflect"
	"regexp"
	"testing"

	p "github.com/wkhere/htmlx/pred"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestEmpty(t *testing.T) {
	var empty Finder

	if !empty.IsEmpty() {
		t.Error("expected empty finder to be, well, empty")
	}
	if empty.String() != "" {
		t.Error("expected empty.String to return empty string")
	}
	if empty.InnerText() != "" {
		t.Error("expected empty.InnerText to return empty string")
	}
	if empty.Parent() != empty {
		t.Error("expected empty.Parent to return empty finder")
	}
	if empty.FirstChild() != empty {
		t.Error("expected empty.FirstChild to return empty finder")
	}
	if empty.LastChild() != empty {
		t.Error("expected empty.LastChild to return empty finder")
	}
	if empty.PrevSibling() != empty {
		t.Error("expected empty.PrevSibling to return empty finder")
	}
	if empty.NextSibling() != empty {
		t.Error("expected empty.NextSibling to return empty finder")
	}

	if empty.Find(p.True()) != empty {
		t.Error("expected empty.Find to return empty finder")
	}
	if empty.FindSibling(p.True()) != empty {
		t.Error("expected empty.FindSibling to return empty finder")
	}
	if empty.FindPrevSibling(p.True()) != empty {
		t.Error("expected empty.FindPrevSibling to return empty finder")
	}
	if empty.FindAll(p.True()).Collect() != nil {
		t.Error("expected empty.FindAll to return empty stream")
	}
	if empty.FindSiblings(p.True()).Collect() != nil {
		t.Error("expected empty.FindSiblings to return empty stream")
	}
	if empty.FindPrevSiblings(p.True()).Collect() != nil {
		t.Error("expected empty.FindSiblings to return empty stream")
	}
	if empty.FindWithSiblings(p.True()).Collect() != nil {
		t.Error("expected empty.FindWithSiblings to return empty stream")
	}

	if ok := empty.Attr().Exists("href"); ok {
		t.Error("expected empty finder to not have href attr")
	}
	if _, ok := empty.Attr().Val("href"); ok {
		t.Error("expected empty finder to not return href attr")
	}

	if r := empty.StreamSelf().First(); !r.IsEmpty() {
		t.Error("expected empty finder streamed and extracted to be empty")
	}
	if rr := empty.StreamSelf().Collect(); len(rr) > 0 {
		t.Error("expected empty finder streamed and collected to be empty")
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

func TestFromNode(t *testing.T) {
	f := testdata("simple.html")
	node, _ := html.Parse(f)
	f.Close()

	top := FinderFromNode(node)
	div0 := top.Find(p.Element(atom.Div))

	if div0.Find(p.Class("bar")) != div0.Find(p.Class("other")) {
		t.Error("mismatch")
	}
}

func TestFromString(t *testing.T) {
	top, _ := FinderFromString(`<div id="1"></div>`)

	div := top.Find(p.Element(atom.Div))

	if res := top.Find(p.ID("1")); res != div {
		t.Error("mismatch")
	}
}

func TestFind(t *testing.T) {
	f := testdata("simple.html")
	top, _ := FinderFromData(f)
	f.Close()

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
		id1.LastChild().Find(p.ID("2")),
		span1.FindSibling(p.Class("xyz")),
		span1.FindSibling(p.Class("xyz")),
		span1.FindPrevSibling(p.Element(atom.Span)),
		spanInner.FindSibling(p.Element(atom.Span)),
		spanInner.FindPrevSibling(p.Element(atom.Span)),
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

func TestFindElementWith1Pred(t *testing.T) {
	top, _ := FinderFromString(`
		<div id="1">
		  <span id="inner"></span>
		</div>
		<span id="2"></span>
	`)

	e := top.Find(p.Element(atom.Span, p.ID("inner")))
	if v, _ := e.Attr().ID(); v != "inner" {
		t.Errorf("expected to find inner element, got: id=`%v`", v)
	}
}

func TestFindElementWith2Preds(t *testing.T) {
	top, _ := FinderFromString(`
		<div id="1">
		  <span id="inner" class="xyz foo"></span>
		</div>
		<span id="2"></span>
	`)

	e := top.Find(p.Element(atom.Span, p.Class("foo"), p.ID("inner")))
	if v, _ := e.Attr().ID(); v != "inner" {
		t.Errorf("expected to find inner element, got: id=`%v`", v)
	}
}

func TestFindAllPrinted(t *testing.T) {
	f := testdata("simple.html")
	top, _ := FinderFromData(f)
	f.Close()

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
	f := testdata("simple.html")
	top, _ := FinderFromData(f)
	f.Close()

	ff := top.FindAll(p.Element(atom.Span)).Collect()

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

func TestCollectInject(t *testing.T) {
	f := testdata("simple.html")
	top, _ := FinderFromData(f)
	f.Close()

	a := top.FindAll(p.Element(atom.Span)).Collect()
	b := Inject(a).Collect()
	if !reflect.DeepEqual(b, a) {
		t.Errorf("mismatch:\n`%s`\nand\n`%s`", b, a)
	}
}

func TestFinderStreamSelectors(t *testing.T) {
	f := testdata("simple.html")
	top, _ := FinderFromData(f)
	f.Close()

	t.Run("First", func(t *testing.T) {
		t.Parallel()
		ff := top.FindAll(p.Element(atom.Span))
		e := ff.First()
		if s, res := "1st", e.InnerText(); res != s {
			t.Errorf("got `%s`, exp `%s`", res, s)
		}
		e2 := ff.First()
		if s, res := "2nd", e2.InnerText(); res != s {
			t.Errorf("got `%s`, exp `%s`", res, s)
		}
	})

	t.Run("Last", func(t *testing.T) {
		t.Parallel()
		ff := top.FindAll(p.Element(atom.Span))
		e := ff.Last()
		if s, res := "inner", e.InnerText(); res != s {
			t.Errorf("got `%s`, exp `%s`", res, s)
		}
		if e2 := ff.First(); !e2.IsEmpty() {
			t.Error("expected FinderStream.Last then First to be empty")
		}
		if e2 := ff.Last(); !e2.IsEmpty() {
			t.Error("expected FinderStream.Last then Last to be empty")
		}
	})

	t.Run("Select", func(t *testing.T) {
		t.Parallel()
		ff := top.FindAll(p.Element(atom.Span))
		e := ff.Select(p.ID("2"))
		if s, res := "2nd", e.InnerText(); res != s {
			t.Errorf("got `%s`, exp `%s`", res, s)
		}
		e2 := ff.First()
		if s, res := "3rd", e2.InnerText(); res != s {
			t.Errorf("got `%s`, exp `%s`", res, s)
		}
		e3 := ff.Select(p.ID("nonexistent"))
		if !e3.IsEmpty() {
			t.Errorf("expected empty finder, got:\n%v", e3)
		}
	})

	t.Run("Filter", func(t *testing.T) {
		t.Parallel()
		ff := top.FindAll(p.Element(atom.Span))
		r := regexp.MustCompile("another$")
		ff2 := ff.Filter(p.ClassCond(r.MatchString))
		ee := ff2.Collect()
		if len(ee) != 2 {
			t.Errorf("got %d, exp %d", len(ee), 2)
		}
		for i, f := range ee {
			if s, res := atom.Span, f.Node.DataAtom; res != s {
				t.Errorf("ee[%d]: got `%s`, exp `%s`", i, res, s)
			}
		}
		if s, res := "3rd", ee[0].InnerText(); res != s {
			t.Errorf("got `%s`, exp `%s`", res, s)
		}
		if s, res := "inner", ee[1].InnerText(); res != s {
			t.Errorf("got `%s`, exp `%s`", res, s)
		}
	})

	t.Run("TakeN", func(t *testing.T) {
		t.Parallel()

		ff := top.FindAll(p.Element(atom.Span))
		ee := ff.TakeN(2).Collect()

		if len(ee) != 2 {
			t.Errorf("got %d, exp %d", len(ee), 2)
		}
		for i, f := range ee {
			if s, res := atom.Span, f.Node.DataAtom; res != s {
				t.Errorf("ee[%d]: got `%s`, exp `%s`", i, res, s)
			}
		}
		if s, res := "1st", ee[0].InnerText(); res != s {
			t.Errorf("got `%s`, exp `%s`", res, s)
		}
		if s, res := "2nd", ee[1].InnerText(); res != s {
			t.Errorf("got `%s`, exp `%s`", res, s)
		}
	})

	t.Run("DropN", func(t *testing.T) {
		t.Parallel()

		ff := top.FindAll(p.Element(atom.Span))
		ee := ff.DropN(2).Collect()

		if len(ee) != 2 {
			t.Errorf("got %d, exp %d", len(ee), 2)
		}
		for i, f := range ee {
			if s, res := atom.Span, f.Node.DataAtom; res != s {
				t.Errorf("ee[%d]: got `%s`, exp `%s`", i, res, s)
			}
		}
		if s, res := "3rd", ee[0].InnerText(); res != s {
			t.Errorf("got `%s`, exp `%s`", res, s)
		}
		if s, res := "inner", ee[1].InnerText(); res != s {
			t.Errorf("got `%s`, exp `%s`", res, s)
		}
	})
}

func TestFamilySelectors(t *testing.T) {
	file := testdata("simple.html")
	top, _ := FinderFromData(file)
	file.Close()

	if !top.Parent().IsEmpty() {
		t.Error("expected top.Parent to be empty")
	}

	f1 := top.Find(p.Element(atom.Span))
	f0 := f1.Parent()
	if id := "1"; !f0.Attr().HasID(id) {
		t.Errorf("expected f1.Parent to have id=%s", id)
	}

	if f1.NextSibling().PrevSibling() != f1 {
		t.Error("expected NextSibling.PrevSibling to be self")
	}

	var f Finder

	for f = f1.FirstChild(); !f.NextSibling().IsEmpty(); f = f.NextSibling() {
	}
	if f != f1.LastChild() {
		t.Error("expected loop of NextSibling to go from FirstChild to LastChild")
	}

	for f = f1.LastChild(); !f.PrevSibling().IsEmpty(); f = f.PrevSibling() {
	}
	if f != f1.FirstChild() {
		t.Error("expected loop of PrevSibling to go from LastChild to FirstChild")
	}
}

func TestStreamSelf(t *testing.T) {
	top, _ := FinderFromString(`<span id="1">`)

	e := top.Find(p.ID("1"))

	if e2 := e.StreamSelf().First(); e2 != e {
		t.Error("expected finder streamed and extracted to be the same")
	}
}

func TestFindSiblings(t *testing.T) {
	f := testdata("simple.html")
	top, _ := FinderFromData(f)
	f.Close()

	e0 := top.Find(p.Element(atom.Span))

	ff := e0.FindSiblings(p.Element(atom.Span)).Collect()

	if len(ff) != 2 {
		t.Errorf("len: got %d, exp %d", len(ff), 2)
	}
	for i, f := range ff {
		if s, res := atom.Span, f.Node.DataAtom; res != s {
			t.Errorf("ff[%d]: got `%s`, exp `%s`", i, res, s)
		}
	}
	if s, res := "2nd", ff[0].InnerText(); res != s {
		t.Errorf("got `%s`, exp `%s`", res, s)
	}
	if s, res := "3rd", ff[1].InnerText(); res != s {
		t.Errorf("got `%s`, exp `%s`", res, s)
	}

	if e := ff[0].FindPrevSibling(p.Element(atom.Span)); e != e0 {
		t.Errorf("mismatch:\ngot:\n%v\nexp:\n%v", e, e0)
	}

	if p := ff[1].FindSiblings(p.Element(atom.Span)).Collect(); p != nil {
		t.Errorf("mismatch:\ngot:\n%v\nexp: nil", p)
	}
}

func TestFindPrevSiblings(t *testing.T) {
	f := testdata("simple.html")
	top, _ := FinderFromData(f)
	f.Close()

	e0 := top.Find(p.Element(atom.Span)).
		FindSiblings(p.Element(atom.Span)).Last()

	ff := e0.FindPrevSiblings(p.Element(atom.Span)).Collect()

	if len(ff) != 2 {
		t.Errorf("got %d, exp %d", len(ff), 2)
	}
	for i, f := range ff {
		if s, res := atom.Span, f.Node.DataAtom; res != s {
			t.Errorf("ff[%d]: got `%s`, exp `%s`", i, res, s)
		}
	}
	if s, res := "2nd", ff[0].InnerText(); res != s {
		t.Errorf("got `%s`, exp `%s`", res, s)
	}
	if s, res := "1st", ff[1].InnerText(); res != s {
		t.Errorf("got `%s`, exp `%s`", res, s)
	}

	if e := ff[0].FindSibling(p.Element(atom.Span)); e != e0 {
		t.Errorf("mismatch:\ngot:\n%v\nexp:\n%v", e, e0)
	}

	if p := ff[1].FindPrevSiblings(p.Element(atom.Span)).Collect(); p != nil {
		t.Errorf("mismatch:\ngot:\n%v\nexp: nil", p)
	}
}

func TestFindWithSiblings(t *testing.T) {
	f := testdata("simple.html")
	top, _ := FinderFromData(f)
	f.Close()

	ff := top.FindWithSiblings(p.Element(atom.Span)).Collect()

	if len(ff) != 3 {
		t.Errorf("got %d, exp %d", len(ff), 3)
	}
	for i, f := range ff {
		if s, res := atom.Span, f.Node.DataAtom; res != s {
			t.Errorf("ff[%d]: got `%s`, exp `%s`", i, res, s)
		}
	}

	expText := []string{"1st", "2nd", "3rd"}

	for i, s := range expText {
		if res := ff[i].InnerText(); res != s {
			t.Errorf("got `%s`, exp `%s`", res, s)
		}
	}

}

func TestAttrShortcuts(t *testing.T) {
	f := testdata("simple.html")
	top, _ := FinderFromData(f)
	f.Close()

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
	f := testdata("simple.html")
	top, _ := FinderFromData(f)
	f.Close()

	var r *regexp.Regexp
	var e Finder

	e = top.Find(p.InnerText("2nd"))
	if !e.Attr().HasID("2") {
		t.Error(`failed to find inner text "2nd"`)
	}

	r = regexp.MustCompile("3rd")
	e = top.Find(p.InnerTextCond(r.MatchString))
	if !e.Attr().HasClass("bar") {
		t.Error(`failed to find inner text r"3rd"`)
	}

	r = regexp.MustCompile("[oO]ther")
	e = top.Find(p.ClassCond(r.MatchString))
	if !e.Attr().HasID("2") {
		t.Error(`failed to find class r"[oO]ther"`)
	}

	r = regexp.MustCompile("2")
	e = top.Find(p.IDCond(r.MatchString))
	if !e.Attr().HasID("2") {
		t.Error(`failed to find ID r"2"`)
	}

	r = regexp.MustCompile("bo.m")
	e = top.Find(p.AttrCond("attr2", r.MatchString))
	if !e.Attr().HasID("2") {
		t.Error(`failed to find attr in "attr2" by r"bo.m"`)
	}

	r = regexp.MustCompile("xyz")
	e = top.Find(p.AttrWordCond("class", r.MatchString))
	if !e.Attr().HasClass("xyz") {
		t.Error(`failed to find attr word in "class" by r"xyz"`)
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
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		id1 := f.Find(p.ID("id1"))
		id1.Find(p.Attr("class", "baz"))
	}
}

func BenchmarkGoV(b *testing.B) {
	f := testdata("gatesofvienna.html")
	top, _ := FinderFromData(f)
	f.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		top.Find(p.Class("html-end-of-file"))
	}
}

func BenchmarkTRSiblings(b *testing.B) {
	f := testdata("gatesofvienna.html")
	top, _ := FinderFromData(f)
	f.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		trs := top.FindWithSiblings(p.Element(atom.Tr))
		trs.Collect()
	}
}

func testdata(filename string) (f *os.File) {
	f, err := os.Open(path.Join("testdata", filename))
	if err != nil {
		panic(err)
	}
	return
}
