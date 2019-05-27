package convertor

import (
	"github.com/rsp9u/seq2xls/model"
	"github.com/rsp9u/seq2xls/seqdiag/ast"
)

// AstToModel converts from the sequence diagram AST of 'seqdiag' to the drawable model.
func AstToModel(d *ast.Diagram) (*model.SequenceDiagram, error) {
	seq := &model.SequenceDiagram{}

	lls, err := ExtractLifelines(d)
	if err != nil {
		return nil, err
	}
	seq.Lifelines = lls

	err = ScanTimeline(d, seq)
	if err != nil {
		return nil, err
	}

	return seq, nil
}
