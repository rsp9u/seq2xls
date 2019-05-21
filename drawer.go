package seq2xls

import (
	"github.com/rsp9u/go-xlsshape/oxml"
	"github.com/rsp9u/go-xlsshape/oxml/shape"
	"github.com/rsp9u/seq2xls/model"
)

const (
	marginX = 20
	marginY = 20
	sizeX   = 6 * 20
	sizeY   = 3 * 20
	spanX   = 192
	spanY   = 40
	tailY   = 60
)

// DrawLifelines adds the shapes composes 'Lifeline' into the spreadsheet.
//
// 'Lifeline' is composed of a rectangle and a dashed line.
func DrawLifelines(ss *oxml.Spreadsheet, lls []*model.Lifeline, nMsg int) {

	for _, ll := range lls {
		i := ll.Index
		rect := shape.NewRectangle()
		rect.SetLeftTop(marginX+spanX*i, marginY)
		rect.SetSize(sizeX, sizeY)
		rect.SetText(ll.Name, "en-US")
		rect.SetHAlign("ctr")
		rect.SetVAlign("ctr")
		ss.AddShape(rect)

		rectXCenter := calcLifelineCenterX(ll)
		rectBottom := marginY + sizeY
		line := shape.NewLine()
		line.SetStartPos(rectXCenter, rectBottom)
		line.SetEndPos(rectXCenter, rectBottom+spanY*(nMsg+1)+tailY)
		line.SetDashType("dash")
		ss.AddShape(line)
	}
}

func calcLifelineCenterX(ll *model.Lifeline) int {
	return marginX + spanX*ll.Index + sizeX/2
}

// DrawMessages adds the shapes of 'Message' into the spreadsheet.
//
// 'Message' is an arrow line.
func DrawMessages(ss *oxml.Spreadsheet, msgs []*model.Message) {
	baseY := marginY + sizeY
	for _, msg := range msgs {
		i := msg.Index
		y := baseY + spanY*(i+1)
		line := shape.NewLine()
		line.SetStartPos(calcLifelineCenterX(msg.From), y)
		line.SetEndPos(calcLifelineCenterX(msg.To), y)
		switch msg.Type {
		case model.Asynchronous:
			line.SetTailType("arrow")
		case model.Reply:
			line.SetTailType("arrow")
			line.SetDashType("dash")
		default:
			line.SetTailType("triangle")
		}
		ss.AddShape(line)
	}
}
