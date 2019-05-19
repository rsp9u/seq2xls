package seqdiag

import (
	"testing"

	"github.com/rsp9u/seq2xls/seqdiag/ast"
	"github.com/rsp9u/seq2xls/seqdiag/lexer"
	"github.com/rsp9u/seq2xls/seqdiag/parser"
)

const testDataLifeline = `
seqdiag {
  webserver; database; browser;
  browser  -> webserver [label = "GET /index.html"];
  browser <-- webserver;
  browser  -> webserver [label = "POST /blog/comment"];
              webserver  -> database [label = "INSERT comment"];
              webserver <-- database;
  browser <-- webserver;
}
`

func TestExtractLifelines(t *testing.T) {
	input := []byte(testDataLifeline)
	lex := lexer.NewLexer(input)
	p := parser.NewParser()
	st, err := p.Parse(lex)
	if err != nil {
		panic(err)
	}

	d, ok := st.(*ast.Diagram)
	if !ok {
		t.Fatalf("This is not a seqdiag")
	}
	lls, err := ExtractLifelines(d)
	if err != nil {
		t.Fatalf("Extract error %v", err)
	}

	if len(lls) != 3 {
		t.Fatalf("Too many lifelines %d", len(lls))
	}

	for _, ll := range lls {
		switch ll.Name {
		case "browser":
			if ll.Index != 2 {
				t.Fatalf("Invalid lifeline index %s[%d]", ll.Name, ll.Index)
			}
		case "webserver":
			if ll.Index != 0 {
				t.Fatalf("Invalid lifeline index %s[%d]", ll.Name, ll.Index)
			}
		case "database":
			if ll.Index != 1 {
				t.Fatalf("Invalid lifeline index %s[%d]", ll.Name, ll.Index)
			}
		default:
			t.Fatalf("Unknown lifeline %s", ll.Name)
		}
	}
}
