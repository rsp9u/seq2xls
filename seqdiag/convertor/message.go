package convertor

import (
	"strings"

	"github.com/golang-collections/collections/stack"
	"github.com/rsp9u/seq2xls/model"
	"github.com/rsp9u/seq2xls/seqdiag/ast"
)

// ScanTimeline extracts all time series elements from the diagram AST and puts them into the given diagram model.
func ScanTimeline(d *ast.Diagram, seq *model.SequenceDiagram) error {
	seq.Messages = []*model.Message{}
	seq.Notes = []*model.Note{}

	_, err := scanTimelineFromStmts(d.Stmts.Items, seq, 0)
	if err != nil {
		return err
	}
	return nil
}

func scanTimelineFromStmts(stmts []ast.Stmt, seq *model.SequenceDiagram, index int) (int, error) {
	var (
		indexCnt, indexPlus int
		err                 error
	)
	indexCnt = 0

	for _, stmt := range stmts {
		switch v := stmt.(type) {
		case ast.ContainerStmt:
			indexPlus, err = scanTimelineFromStmts(v.GetItems(), seq, index)
			if err != nil {
				return 0, err
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
					From:     getLifeline(seq.Lifelines, getFromNode(sgmt).Value),
					To:       getLifeline(seq.Lifelines, getToNode(sgmt).Value),
					Type:     edgeType,
					ColorHex: "000000",
					Text:     text,
				}
				seq.Messages = append(seq.Messages, msg)
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
					seq.Notes = append(seq.Notes, &model.Note{
						Assoc:    msg,
						OnLeft:   lnote.OnLeft,
						Text:     lnote.Text,
						ColorHex: lnote.ColorHex,
					})
				}
				if rnote != nil {
					seq.Notes = append(seq.Notes, &model.Note{
						Assoc:    msg,
						OnLeft:   rnote.OnLeft,
						Text:     rnote.Text,
						ColorHex: rnote.ColorHex,
					})
				}
			}

			if v.EdgeBlock != nil {
				indexPlus, err = scanTimelineFromStmts(v.EdgeBlock.Items, seq, index)
				if err != nil {
					return 0, err
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
						From:     getLifeline(seq.Lifelines, getFromNode(sgmt).Value),
						To:       getLifeline(seq.Lifelines, getToNode(sgmt).Value),
						Type:     getMessageType(sgmt),
						ColorHex: "000000",
					}
					seq.Messages = append(seq.Messages, msg)
					index++
					indexCnt++
				}
			}
		}
	}

	return indexCnt, nil
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
