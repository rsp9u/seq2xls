package seqdiag

import (
	"github.com/rsp9u/seq2xls/model"
	"github.com/rsp9u/seq2xls/seqdiag/ast"
)

// ExtractLifelines extracts lifeline elements from the diagram.
func ExtractLifelines(d *ast.Diagram) ([]*model.Lifeline, error) {
	lls := []*model.Lifeline{}
	lls, _, err := extractLifelinesFromStmts(d.Stmts.Items, lls, 0)
	if err != nil {
		return nil, err
	}
	return lls, nil
}

func extractLifelinesFromStmts(stmts []ast.Stmt, lls []*model.Lifeline, index int) ([]*model.Lifeline, int, error) {
	var (
		indexCnt, indexPlus int
		err                 error
	)
	indexCnt = 0

	for _, stmt := range stmts {
		switch v := stmt.(type) {
		case *ast.FragmentStmt:
		case *ast.GroupStmt:
			lls, indexPlus, err = extractLifelinesFromStmts(v.Stmts.Items, lls, index)
			if err != nil {
				return nil, 0, err
			}
			index += indexPlus

		case *ast.EdgeStmt:
			for _, sgmt := range v.EdgeSegments.Items {
				if !containsLifeline(lls, sgmt.LeftNode.Value) {
					ll := &model.Lifeline{Name: sgmt.LeftNode.Value, Index: index, ColorHex: "FFFFFF"}
					lls = append(lls, ll)
					index++
					indexCnt++
				}

				if !containsLifeline(lls, sgmt.RightNode.Value) {
					ll := &model.Lifeline{Name: sgmt.RightNode.Value, Index: index, ColorHex: "FFFFFF"}
					lls = append(lls, ll)
					index++
					indexCnt++
				}
			}

			if v.EdgeBlock != nil {
				lls, indexPlus, err = extractLifelinesFromStmts(v.EdgeBlock.Items, lls, index)
				if err != nil {
					return nil, 0, err
				}
				index += indexPlus
			}

		case *ast.NodeStmt:
			if !containsLifeline(lls, v.ID.Value) {
				ll := &model.Lifeline{Name: v.ID.Value, Index: index, ColorHex: "FFFFFF"}
				lls = append(lls, ll)
				index++
				indexCnt++
			}
		}
	}

	return lls, indexCnt, nil
}

func containsLifeline(lls []*model.Lifeline, name string) bool {
	for _, ll := range lls {
		if ll.Name == name {
			return true
		}
	}
	return false
}
