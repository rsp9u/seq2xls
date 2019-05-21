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

// ExtractMessages extracts message elements from the diagram.
func ExtractMessages(d *ast.Diagram, lls []*model.Lifeline) ([]*model.Message, error) {
	msgs := []*model.Message{}
	msgs, _, err := extractMessagesFromStmts(d.Stmts.Items, lls, msgs, 0)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func extractMessagesFromStmts(stmts []ast.Stmt, lls []*model.Lifeline, msgs []*model.Message, index int) ([]*model.Message, int, error) {
	var (
		indexCnt, indexPlus int
		err                 error
	)
	indexCnt = 0

	for _, stmt := range stmts {
		switch v := stmt.(type) {
		case *ast.FragmentStmt:
		case *ast.GroupStmt:
			msgs, indexPlus, err = extractMessagesFromStmts(v.Stmts.Items, lls, msgs, index)
			if err != nil {
				return nil, 0, err
			}
			index += indexPlus

		case *ast.EdgeStmt:
			for _, sgmt := range v.EdgeSegments.Items {
				msg := &model.Message{
					Index: index,
					From: getLifeline(lls, sgmt.LeftNode.Value),
					To: getLifeline(lls, sgmt.RightNode.Value),
					Type: getMessageType(sgmt.Edge),
					ColorHex: "000000",
				}
				msgs = append(msgs, msg)
				index++
				indexCnt++
			}

			if v.EdgeBlock != nil {
				lls, indexPlus, err = extractLifelinesFromStmts(v.EdgeBlock.Items, lls, index)
				if err != nil {
					return nil, 0, err
				}
				index += indexPlus
			}
		}
	}

	return lls, indexCnt, nil
}
