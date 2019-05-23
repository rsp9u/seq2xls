package seqdiag

import (
	"strings"

	"github.com/golang-collections/collections/stack"
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
		case ast.ContainerStmt:
			lls, indexPlus, err = extractLifelinesFromStmts(v.GetItems(), lls, index)
			if err != nil {
				return nil, 0, err
			}
			index += indexPlus
			indexCnt += indexPlus

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
				indexCnt += indexPlus
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
		case ast.ContainerStmt:
			msgs, indexPlus, err = extractMessagesFromStmts(v.GetItems(), lls, msgs, index)
			if err != nil {
				return nil, 0, err
			}
			index += indexPlus
			indexCnt += indexPlus

		case *ast.EdgeStmt:
			tripReplySgmts := stack.New()
			for _, sgmt := range v.EdgeSegments.Items {
				edgeType := getMessageType(sgmt)
				msg := &model.Message{
					Index:    index,
					From:     getLifeline(lls, getFromNode(sgmt).Value),
					To:       getLifeline(lls, getToNode(sgmt).Value),
					Type:     edgeType,
					ColorHex: "000000",
				}
				msgs = append(msgs, msg)
				index++
				indexCnt++

				if edgeType != model.SelfReference && isTripMessage(sgmt) {
					tripReplySgmts.Push(&ast.EdgeSegment{
						LeftNode:  sgmt.LeftNode,
						RightNode: sgmt.RightNode,
						Edge:      "<-",
					})
				}
			}

			if v.EdgeBlock != nil {
				msgs, indexPlus, err = extractMessagesFromStmts(v.EdgeBlock.Items, lls, msgs, index)
				if err != nil {
					return nil, 0, err
				}
				index += indexPlus
				indexCnt += indexPlus
			}

			for tripReplySgmts.Len() != 0 {
				s := tripReplySgmts.Pop()
				sgmt, ok := s.(*ast.EdgeSegment)
				if ok {
					msg := &model.Message{
						Index:    index,
						From:     getLifeline(lls, getFromNode(sgmt).Value),
						To:       getLifeline(lls, getToNode(sgmt).Value),
						Type:     getMessageType(sgmt),
						ColorHex: "000000",
					}
					msgs = append(msgs, msg)
					index++
					indexCnt++
				}
			}
		}
	}

	return msgs, indexCnt, nil
}

func getLifeline(lls []*model.Lifeline, name string) *model.Lifeline {
	for _, ll := range lls {
		if ll.Name == name {
			return ll
		}
	}
	return nil
}

func getFromNode(sgmt *ast.EdgeSegment) *ast.ID {
	if strings.HasSuffix(sgmt.Edge, ">") {
		return sgmt.LeftNode
	}
	return sgmt.RightNode
}

func getToNode(sgmt *ast.EdgeSegment) *ast.ID {
	if strings.HasSuffix(sgmt.Edge, ">") {
		return sgmt.RightNode
	}
	return sgmt.LeftNode
}

func getMessageType(sgmt *ast.EdgeSegment) model.MessageType {
	if sgmt.LeftNode.Value == sgmt.RightNode.Value {
		return model.SelfReference
	}
	switch sgmt.Edge {
	case "->", "->>", "=>":
		return model.Synchronous
	case "-->", "-->>":
		return model.Asynchronous
	case "<-", "<--", "<<-", "<<--":
		return model.Reply
	default:
		return model.Synchronous
	}
}

func isTripMessage(sgmt *ast.EdgeSegment) bool {
	return sgmt.Edge == "=>"
}
