package ast

import (
	"strings"

	"github.com/rsp9u/seq2xls/seqdiag/token"
)

type Attr interface{}

/****************
 * Statement
 ****************/
type Stmt interface{}

/****************
 * Diagram
 ****************/
type Diagram struct {
	ID    *ID
	Stmts *DiagramInlineStmtList
}

type DiagramInlineStmtList struct {
	Items []Stmt
}

func NewDiagram(id, stmts Attr) (*Diagram, error) {
	return &Diagram{id.(*ID), stmts.(*DiagramInlineStmtList)}, nil
}

func NewDiagramInlineStmtList(acc, stmt Attr) (*DiagramInlineStmtList, error) {
	if acc == nil {
		acc = &DiagramInlineStmtList{}
	}
	if stmt != nil {
		stmts := acc.(*DiagramInlineStmtList)
		stmts.Items = append(stmts.Items, stmt.(Stmt))
	}
	return acc.(*DiagramInlineStmtList), nil
}

/****************
 * Fragment Statement
 ****************/
type FragmentStmt struct {
	Type  string
	ID    *ID
	Stmts *FragmentInlineStmtList
}

type FragmentInlineStmtList struct {
	Items []Stmt
}

func NewFragmentStmt(t, id, stmts Attr) (*FragmentStmt, error) {
	return &FragmentStmt{t.(string), id.(*ID), stmts.(*FragmentInlineStmtList)}, nil
}

func NewFragmentInlineStmtList(acc, stmt Attr) (*FragmentInlineStmtList, error) {
	if acc == nil {
		acc = &FragmentInlineStmtList{}
	}
	if stmt != nil {
		stmts := acc.(*FragmentInlineStmtList)
		stmts.Items = append(stmts.Items, stmt.(Stmt))
	}
	return acc.(*FragmentInlineStmtList), nil
}

/****************
 * Group Statement
 ****************/
type GroupStmt struct {
	ID    *ID
	Stmts *GroupInineStmtList
}

type GroupInineStmtList struct {
	Items []Stmt
}

func NewGroupStmt(id, stmts Attr) (*GroupStmt, error) {
	return &GroupStmt{id.(*ID), stmts.(*GroupInineStmtList)}, nil
}

func NewGroupInineStmtList(acc, stmt Attr) (*GroupInineStmtList, error) {
	if acc == nil {
		acc = &GroupInineStmtList{}
	}
	if stmt != nil {
		stmts := acc.(*GroupInineStmtList)
		stmts.Items = append(stmts.Items, stmt.(Stmt))
	}
	return acc.(*GroupInineStmtList), nil
}

/****************
 * Edge Statement
 ****************/
type EdgeStmt struct {
	EdgeSegments *EdgeSegmentList
	Options      *OptionList
	EdgeBlock    *EdgeBlockInlineStmtList
}

type EdgeSegmentList struct {
	Items    []*EdgeSegment
	LastNode *ID
}

type EdgeSegment struct {
	LeftNode  *ID
	Edge      string
	RightNode *ID
}

type EdgeBlockInlineStmtList struct {
	Items []Stmt
}

func NewEdgeStmt(sgmts, opt, blk Attr) (*EdgeStmt, error) {
	return &EdgeStmt{
		sgmts.(*EdgeSegmentList),
		opt.(*OptionList),
		blk.(*EdgeBlockInlineStmtList),
	}, nil
}

func NewEdgeSegmentList(l, e, r Attr) (*EdgeSegmentList, error) {
	sgmts := &EdgeSegmentList{}
	sgmt := &EdgeSegment{
		l.(*ID),
		string(e.(*token.Token).Lit),
		r.(*ID),
	}
	sgmts.Items = append(sgmts.Items, sgmt)
	sgmts.LastNode = r.(*ID)
	return sgmts, nil
}

func AppendEdgeSegment(acc, e, r Attr) (*EdgeSegmentList, error) {
	sgmts := acc.(*EdgeSegmentList)
	sgmt := &EdgeSegment{
		sgmts.LastNode,
		string(e.(*token.Token).Lit),
		r.(*ID),
	}
	sgmts.Items = append(sgmts.Items, sgmt)
	sgmts.LastNode = r.(*ID)
	return sgmts, nil
}

func NewEdgeBlockInlineStmtList(acc, stmt Attr) (*EdgeBlockInlineStmtList, error) {
	if acc == nil {
		acc = &EdgeBlockInlineStmtList{}
	}
	if stmt != nil {
		stmts := acc.(*EdgeBlockInlineStmtList)
		stmts.Items = append(stmts.Items, stmt.(Stmt))
	}
	return acc.(*EdgeBlockInlineStmtList), nil
}

/****************
 * Separator Statement
 ****************/
type SeparatorStmt struct {
	Type, Value string
}

func NewSeparatorStmt(attr Attr) (*SeparatorStmt, error) {
	s := attr.(string)
	return &SeparatorStmt{s[0:3], strings.TrimSpace(s[3 : len(s)-3])}, nil
}

/****************
 * Node Statement
 ****************/
type NodeStmt struct {
	ID     *ID
	Option *OptionList
}

func NewNodeStmt(id, opt Attr) (*NodeStmt, error) {
	return &NodeStmt{id.(*ID), opt.(*OptionList)}, nil
}

/****************
 * Attribute Statement
 ****************/
type AttributeStmt struct {
	Type, Value *ID
}

func NewAttributeStmt(t, v Attr) (*AttributeStmt, error) {
	return &AttributeStmt{t.(*ID), v.(*ID)}, nil
}

/****************
 * Option
 ****************/
type OptionList struct {
	Items []*Option
}

type Option struct {
	Type, Value *ID
}

func NewOptionList(acc, opt Attr) (*OptionList, error) {
	if acc == nil {
		acc = &OptionList{[]*Option{}}
	}
	if opt != nil {
		opts := acc.(*OptionList)
		opts.Items = append(opts.Items, opt.(*Option))
	}
	return acc.(*OptionList), nil
}

func NewOption(t, v Attr) (*Option, error) {
	return &Option{t.(*ID), v.(*ID)}, nil
}

/****************
 * ID
 ****************/
type ID struct {
	Value string
}

func NewID(id, t Attr) (*ID, error) {
	switch t.(string) {
	case "string":
		return &ID{TrimQuote(id)}, nil
	default:
		return &ID{TokenToString(id)}, nil
	}
}

func NewEmptyID() *ID {
	return &ID{}
}

func TokenToString(attr Attr) string {
	return string(attr.(*token.Token).Lit)
}

func TrimQuote(attr Attr) string {
	s := string(attr.(*token.Token).Lit)
	return s[1 : len(s)-1]
}

/****************
 * Interfaces
 ****************/
type ContainerStmt interface {
	GetItems() []Stmt
}

func (s *FragmentStmt) GetItems() []Stmt {
	return s.Stmts.Items
}

func (s *GroupStmt) GetItems() []Stmt {
	return s.Stmts.Items
}
