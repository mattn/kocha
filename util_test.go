package kocha

import (
	"go/format"
	"html/template"
	"reflect"
	"regexp"
	"testing"
)

func TestToCamelCase(t *testing.T) {
	for v, expected := range map[string]string{
		"kocha":   "Kocha",
		"KochA":   "KochA",
		"koch_a":  "KochA",
		"k_oc_ha": "KOcHa",
		"k_Oc_hA": "KOcHA",
	} {
		actual := ToCamelCase(v)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("%v: Expect %v, but %v", v, expected, actual)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	for v, expected := range map[string]string{
		"kocha":  "kocha",
		"Kocha":  "kocha",
		"kochA":  "koch_a",
		"kOcHa":  "k_oc_ha",
		"ko_cha": "ko_cha",
	} {
		actual := ToSnakeCase(v)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("%v: Expect %v, but %v", v, expected, actual)
		}
	}
}

func Test_normPath(t *testing.T) {
	for v, expected := range map[string]string{
		"/":           "/",
		"/path":       "/path",
		"/path/":      "/path/",
		"//path//":    "/path/",
		"/path/to":    "/path/to",
		"/path/to///": "/path/to/",
	} {
		actual := normPath(v)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("%v: Expect %v, but %v", v, expected, actual)
		}
	}
}

func TestGoString(t *testing.T) {
	re := regexp.MustCompile(`^/path/to/([^/]+)(?:\.html)?$`)
	actual := GoString(re)
	expected := `regexp.MustCompile("^/path/to/([^/]+)(?:\\.html)?$")`
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}

	tmpl := template.Must(template.New("test").Parse(`foo{{.name}}bar`))
	actual = GoString(tmpl)
	expected = `template.Must(template.New("test").Funcs(kocha.TemplateFuncs).Parse(kocha.Gunzip("\x1f\x8b\b\x00\x00\tn\x88\x02\xffJ\xcbϯ\xae\xd6\xcbK\xccM\xad\xadMJ,\x02\x04\x00\x00\xff\xff4%\x83\xb6\x0f\x00\x00\x00")))`
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}

	actual = GoString(testGoString{})
	expected = "gostring"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}

	actual = GoString(nil)
	expected = "nil"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}

	var ptr *int
	actual = GoString(ptr)
	expected = "nil"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}

	n := 10
	nptr := &n
	actual = GoString(nptr)
	expected = "10"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}

	aBuf, err := format.Source([]byte(GoString(struct {
		Name, path string
		Route      map[string]interface{}
		G          *testGoString
	}{
		Name: "foo",
		path: "path",
		Route: map[string]interface{}{
			"first":  "Tokyo",
			"second": "Kyoto",
			"third":  []int{10, 11, 20},
		},
		G: &testGoString{},
	})))
	if err != nil {
		t.Fatal(err)
	}
	eBuf, err := format.Source([]byte(`
struct {
	Name string
	path string
	Route map[string]interface{}
	G *kocha.testGoString
}{

	G: gostring,

	Name: "foo",

	Route: map[string]interface{}{

		"first": "Tokyo",

		"second": "Kyoto",

		"third": []int{

			10,

			11,

			20,
		},
	},
}`))
	if err != nil {
		t.Fatal(err)
	}
	actual = string(aBuf)
	expected = string(eBuf)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %q, but %q", expected, actual)
	}
}

type testGoString struct{}

func (g testGoString) GoString() string { return "gostring" }

func Test_Gzip(t *testing.T) {
	actual := Gzip("kocha")
	expected := "\x1f\x8b\b\x00\x00\tn\x88\x02\xff\xca\xceO\xceH\x04\x04\x00\x00\xff\xff\f\x93\x85\x96\x05\x00\x00\x00"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}

	// reversibility test
	gz := Gzip("kocha")
	actual = Gunzip(gz)
	expected = "kocha"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}
}

func TestGunzip(t *testing.T) {
	actual := Gunzip("\x1f\x8b\b\x00\x00\tn\x88\x02\xff\xca\xceO\xceH\x04\x04\x00\x00\xff\xff\f\x93\x85\x96\x05\x00\x00\x00")
	expected := "kocha"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}

	// reversibility test
	raw := Gunzip("\x1f\x8b\b\x00\x00\tn\x88\x02\xff\xca\xceO\xceH\x04\x04\x00\x00\xff\xff\f\x93\x85\x96\x05\x00\x00\x00")
	actual = Gzip(raw)
	expected = "\x1f\x8b\b\x00\x00\tn\x88\x02\xff\xca\xceO\xceH\x04\x04\x00\x00\xff\xff\f\x93\x85\x96\x05\x00\x00\x00"
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expect %v, but %v", expected, actual)
	}
}
