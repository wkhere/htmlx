package find

import (
	"testing"

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
}

func TestFind(t *testing.T) {
	var s string
	top, _ := FinderFromString(`
		<div>
			<div id="id1">
				<span class="foo">1st</span>
				<span id="id2" class="bar">2nd</span>
				<span class="bar">3rd</span>
				<div>
					<span class="xyz">inner</span>
				</div>
			</div>
		</div>
	`)
	id1 := top.FindById("id1")

	span1 := id1.FindElement(atom.Span)
	span1content := span1.FirstChild()
	s = "1st"
	if res := span1content.String(); res != s {
		t.Errorf("mismatch:\ngot `%s`\nexp `%s`", res, s)
	}

	span2 := id1.FindByAttr("class", "bar")
	span2content := span2.FirstChild()
	s = "2nd"
	if res := span2content.String(); res != s {
		t.Errorf("mismatch:\ngot `%s`\nexp `%s`", res, s)
	}

	if res := span1.FindSiblingElement(atom.Span); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span2)
	}

	if res := span1.FindSiblingByAttr("class", "bar"); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`",
			res.String(), span2.String())
	}

	if res := span1.FindSiblingById("id2"); res != span2 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span2)
	}

	span3 := span2.FindSiblingElement(atom.Span)
	if res := span2.FindSiblingByAttr("class", "bar"); res != span3 {
		t.Errorf("mismatch:\ngot `%v`\nexp `%v`", res, span3)
	}

	if bad := top.FindById("bad"); !bad.IsEmpty() {
		t.Errorf("mismatch:\ngot `%v`\nexp empty find", bad)
	}
	if bad := id1.FindByAttr("class", "bad"); !bad.IsEmpty() {
		t.Errorf("mismatch:\ngot `%v`\nexp empty find", bad)
	}
	if bad := id1.NextSibling().FindSiblingById("any"); !bad.IsEmpty() {
		t.Errorf("mismatch:\ngot `%v`\nexp empty find", bad)
	}
	if bad := id1.FirstChild().FindSiblingById("bad"); !bad.IsEmpty() {
		t.Errorf("mismatch:\ngot `%v`\nexp empty find", bad)
	}
	if bad := span1.FindSiblingByAttr("class", "xyz"); !bad.IsEmpty() {
		t.Errorf("mismatch:\ngot `%s`\nexp empty find", bad)
	}
}

// todo: discover why it got slower.
// maybe (f Finder) -> (f *Finder) ?
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
