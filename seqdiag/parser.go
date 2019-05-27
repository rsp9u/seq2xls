package seqdiag

import (
	"github.com/rsp9u/seq2xls/seqdiag/ast"
	"github.com/rsp9u/seq2xls/seqdiag/lexer"
	"github.com/rsp9u/seq2xls/seqdiag/parser"
)

// ParseSeqdiag parses the given 'seqdiag' text and converts into Go structures.
func ParseSeqdiag(b []byte) *ast.Diagram {
	lex := lexer.NewLexer(b)
	p := parser.NewParser()
	st, err := p.Parse(lex)
	if err != nil {
		panic(err)
	}

	d, ok := st.(*ast.Diagram)
	if !ok {
		panic("This is not a seqdiag")
	}
	return d
}
