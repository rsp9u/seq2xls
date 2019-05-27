package extract

import (
	"testing"

	"github.com/rsp9u/seq2xls/model"
)

const testDataLifeline = `
seqdiag {
  webserver; database; browser;
  browser  -> webserver [label = "GET /index.html"];
  browser <-- webserver;
  group {
    foo; bar;
  }
  loop {
    baz => qux; 
  }
}
`

func checkLifeline(t *testing.T, ll *model.Lifeline, idx int, name string) {
	if ll.Index != idx {
		t.Fatalf("Invalid lifeline index %s[%d]", ll.Name, ll.Index)
	}
	if ll.Name != name {
		t.Fatalf("Mismatches lifeline name [expect: %s, actual: %s]", name, ll.Name)
	}
}

func TestExtractLifelines(t *testing.T) {
	d := ParseSeqdiag([]byte(testDataLifeline))
	lls, err := ExtractLifelines(d)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}

	if len(lls) > 7 {
		t.Fatalf("Too many lifelines %d", len(lls))
	}

	checkLifeline(t, lls[0], 0, "webserver")
	checkLifeline(t, lls[1], 1, "database")
	checkLifeline(t, lls[2], 2, "browser")
	checkLifeline(t, lls[3], 3, "foo")
	checkLifeline(t, lls[4], 4, "bar")
	checkLifeline(t, lls[5], 5, "baz")
	checkLifeline(t, lls[6], 6, "qux")
}
