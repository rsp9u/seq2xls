package convertor

import (
	"fmt"
	"strings"

	"github.com/golang-collections/collections/stack"
	"github.com/rsp9u/seq2xls/model"
	"github.com/rsp9u/seq2xls/seqdiag/ast"
)

// ScanTimeline extracts all time series elements from the diagram AST and puts them into the given diagram model.
func ScanTimeline(d *ast.Diagram, seq *model.SequenceDiagram) error {
	seq.Messages = []*model.Message{}
	seq.Fragments = []*model.Fragment{}
	seq.Notes = []*model.Note{}

	err := scanTimelineInStmts(d.Stmts.Items, seq)
	if err != nil {
		return err
	}
	return nil
}

func scanTimelineInStmts(stmts []ast.Stmt, seq *model.SequenceDiagram) error {
	for _, stmt := range stmts {
		switch v := stmt.(type) {
		case *ast.FragmentStmt:
			frag := &model.Fragment{
				Index: len(seq.Fragments),
				Type:  getFragmentType(v),
			}
			seq.Fragments = append(seq.Fragments, frag)

			beginIndex := len(seq.Messages)
			err := scanTimelineInStmts(v.GetItems(), seq)
			endIndex := len(seq.Messages) - 1
			if err != nil {
				return err
			}
			if endIndex < beginIndex {
				return fmt.Errorf("empty fragment is not allowed")
			}

			frag.Begin = seq.Messages[beginIndex]
			frag.End = seq.Messages[endIndex]

		case *ast.GroupStmt:
			err := scanTimelineInStmts(v.GetItems(), seq)
			if err != nil {
				return err
			}

		case *ast.EdgeStmt:
			tripReplySgmts := stack.New()
			text := getMessageLabel(v)
			lnote := getMessageLeftNote(v)
			rnote := getMessageRightNote(v)

			for _, sgmt := range v.EdgeSegments.Items {
				edgeType := getMessageType(sgmt)
				msg := &model.Message{
					Index:    len(seq.Messages),
					From:     getLifeline(seq.Lifelines, getFromNode(sgmt).Value),
					To:       getLifeline(seq.Lifelines, getToNode(sgmt).Value),
					Type:     edgeType,
					ColorHex: "000000",
					Text:     text,
				}
				seq.Messages = append(seq.Messages, msg)

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
				err := scanTimelineInStmts(v.EdgeBlock.Items, seq)
				if err != nil {
					return err
				}
			}

			for tripReplySgmts.Len() != 0 {
				s := tripReplySgmts.Pop()
				sgmt, ok := s.(*ast.EdgeSegment)
				if ok {
					msg := &model.Message{
						Index:    len(seq.Messages),
						From:     getLifeline(seq.Lifelines, getFromNode(sgmt).Value),
						To:       getLifeline(seq.Lifelines, getToNode(sgmt).Value),
						Type:     getMessageType(sgmt),
						ColorHex: "000000",
					}
					seq.Messages = append(seq.Messages, msg)
				}
			}
		case *ast.SeparatorStmt:
			var beforeMsg *model.Message
			if len(seq.Messages) > 0 {
				beforeMsg = seq.Messages[len(seq.Messages)-1]
			}
			sep := &model.Separator{
				Text:   v.Value,
				Before: beforeMsg,
			}
			seq.Separators = append(seq.Separators, sep)
		}
	}

	return nil
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

func getFragmentType(stmt *ast.FragmentStmt) model.FragmentType {
	switch stmt.Type {
	case "loop":
		return model.Loop
	case "alt":
		return model.Alt
	default:
		return model.UnknownFragment
	}
}
