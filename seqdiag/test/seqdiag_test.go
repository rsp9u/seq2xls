package test

import (
	"testing"

	"github.com/rsp9u/seq2xls/seqdiag/ast"
	"github.com/rsp9u/seq2xls/seqdiag/lexer"
	"github.com/rsp9u/seq2xls/seqdiag/parser"
)

const testDataSimple = `
seqdiag {
  browser  -> webserver [label = "GET /index.html"];
  browser <-- webserver;
  browser  -> webserver [label = "POST /blog/comment"];
              webserver  -> database [label = "INSERT comment"];
              webserver <-- database;
  browser <-- webserver;
}
`

func checkEqual(t *testing.T, act, exp, errfmt string) {
	if act != exp {
		t.Fatalf(errfmt, act)
	}
}

func checkEqualInt(t *testing.T, act, exp int, errfmt string) {
	if act != exp {
		t.Fatalf(errfmt, act)
	}
}

func checkEdgeSgmt(t *testing.T, sgmt *ast.EdgeSegment, l, r, e string) {
	checkEqual(t, sgmt.LeftNode.Value, l, "Wrong node name %v")
	checkEqual(t, sgmt.RightNode.Value, r, "Wrong node name %v")
	checkEqual(t, sgmt.Edge, e, "Wrong edge %v")
}

func TestSimple(t *testing.T) {
	input := []byte(testDataSimple)
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

	var e *ast.EdgeStmt

	e = d.Stmts.Items[0].(*ast.EdgeStmt)
	checkEdgeSgmt(t, e.EdgeSegments.Items[0], "browser", "webserver", "->")
	checkEqual(t, e.Options.Items[0].Type.Value, `label`, "Wrong option type %v")
	checkEqual(t, e.Options.Items[0].Value.Value, `GET /index.html`, "Wrong option value %v")

	e = d.Stmts.Items[1].(*ast.EdgeStmt)
	checkEdgeSgmt(t, e.EdgeSegments.Items[0], "browser", "webserver", "<--")
	checkEqualInt(t, len(e.Options.Items), 0, "Wrong option size %v")

	e = d.Stmts.Items[2].(*ast.EdgeStmt)
	checkEdgeSgmt(t, e.EdgeSegments.Items[0], "browser", "webserver", "->")
	checkEqual(t, e.Options.Items[0].Type.Value, `label`, "Wrong option type %v")
	checkEqual(t, e.Options.Items[0].Value.Value, `POST /blog/comment`, "Wrong option value %v")

	e = d.Stmts.Items[3].(*ast.EdgeStmt)
	checkEdgeSgmt(t, e.EdgeSegments.Items[0], "webserver", "database", "->")
	checkEqual(t, e.Options.Items[0].Type.Value, `label`, "Wrong option type %v")
	checkEqual(t, e.Options.Items[0].Value.Value, `INSERT comment`, "Wrong option value %v")

	e = d.Stmts.Items[4].(*ast.EdgeStmt)
	checkEdgeSgmt(t, e.EdgeSegments.Items[0], "webserver", "database", "<--")
	checkEqualInt(t, len(e.Options.Items), 0, "Wrong option size %v")

	e = d.Stmts.Items[5].(*ast.EdgeStmt)
	checkEdgeSgmt(t, e.EdgeSegments.Items[0], "browser", "webserver", "<--")
	checkEqualInt(t, len(e.Options.Items), 0, "Wrong option size %v")
}
