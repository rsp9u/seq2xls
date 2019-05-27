package extract

import (
	"strings"

	"github.com/golang-collections/collections/stack"
	"github.com/rsp9u/seq2xls/model"
	"github.com/rsp9u/seq2xls/seqdiag/ast"
)

// ExtractMessages extracts message elements from the diagram.
func ExtractMessages(d *ast.Diagram, lls []*model.Lifeline) ([]*model.Message, []*model.Note, error) {
	msgs := []*model.Message{}
	notes := []*model.Note{}
	msgs, notes, _, err := extractMessagesFromStmts(d.Stmts.Items, lls, msgs, notes, 0)
	if err != nil {
		return nil, nil, err
	}
	return msgs, notes, nil
}

func extractMessagesFromStmts(stmts []ast.Stmt, lls []*model.Lifeline, msgs []*model.Message, notes []*model.Note, index int) ([]*model.Message, []*model.Note, int, error) {
	var (
		indexCnt, indexPlus int
		err                 error
	)
	indexCnt = 0

	for _, stmt := range stmts {
		switch v := stmt.(type) {
		case ast.ContainerStmt:
			msgs, notes, indexPlus, err = extractMessagesFromStmts(v.GetItems(), lls, msgs, notes, index)
			if err != nil {
				return nil, nil, 0, err
			}
			index += indexPlus
			indexCnt += indexPlus

		case *ast.EdgeStmt:
			tripReplySgmts := stack.New()
			text := getMessageLabel(v)
			lnote := getMessageLeftNote(v)
			rnote := getMessageRightNote(v)

			for _, sgmt := range v.EdgeSegments.Items {
				edgeType := getMessageType(sgmt)
				msg := &model.Message{
					Index:    index,
					From:     getLifeline(lls, getFromNode(sgmt).Value),
					To:       getLifeline(lls, getToNode(sgmt).Value),
					Type:     edgeType,
					ColorHex: "000000",
					Text:     text,
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

				if lnote != nil {
					notes = append(notes, &model.Note{
						Assoc:    msg,
						OnLeft:   lnote.OnLeft,
						Text:     lnote.Text,
						ColorHex: lnote.ColorHex,
					})
				}
				if rnote != nil {
					notes = append(notes, &model.Note{
						Assoc:    msg,
						OnLeft:   rnote.OnLeft,
						Text:     rnote.Text,
						ColorHex: rnote.ColorHex,
					})
				}
			}

			if v.EdgeBlock != nil {
				msgs, notes, indexPlus, err = extractMessagesFromStmts(v.EdgeBlock.Items, lls, msgs, notes, index)
				if err != nil {
					return nil, nil, 0, err
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

	return msgs, notes, indexCnt, nil
}

func getLifeline(lls []*model.Lifeline, name string) *model.Lifeline {
	for _, ll := range lls {
		if ll.Name == name {
			return ll
		}
	}
	return nil
}

func getMessageLabel(stmt *ast.EdgeStmt) string {
	for _, opt := range stmt.Options.Items {
		if opt.Type.String() == "label" {
			return opt.Value.String()
		}
	}
	return ""
}

func getMessageLeftNote(stmt *ast.EdgeStmt) *model.Note {
	for _, opt := range stmt.Options.Items {
		if opt.Type.String() == "leftnote" {
			return &model.Note{
				Assoc:    nil,
				OnLeft:   true,
				Text:     opt.Value.String(),
				ColorHex: "ffb6c1",
			}
		}
	}
	return nil
}

func getMessageRightNote(stmt *ast.EdgeStmt) *model.Note {
	for _, opt := range stmt.Options.Items {
		if opt.Type.String() == "note" || opt.Type.String() == "rightnote" {
			return &model.Note{
				Assoc:    nil,
				OnLeft:   false,
				Text:     opt.Value.String(),
				ColorHex: "ffb6c1",
			}
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
