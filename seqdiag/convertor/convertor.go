package convertor

import (
	"github.com/rsp9u/seq2xls/model"
	"github.com/rsp9u/seq2xls/seqdiag/ast"
)

// AstToModel converts from the sequence diagram AST of 'seqdiag' to the drawable model.
func AstToModel(d *ast.Diagram) (*model.SequenceDiagram, error) {
	lls, err := ExtractLifelines(d)
	if err != nil {
		return nil, err
	}

	msgs, notes, err := ExtractMessages(d, lls)
	if err != nil {
		return nil, err
	}

	return &model.SequenceDiagram{
		Lifelines: lls,
		Messages:  msgs,
		Notes:     notes,
	}, nil
}
