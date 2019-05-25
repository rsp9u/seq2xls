package seq2xls

import (
	"strings"

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
		if msg.Type != model.SelfReference {
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
		} else {
			w := spanX / 3
			h := spanY / 3
			line1 := shape.NewLine()
			line2 := shape.NewLine()
			line3 := shape.NewLine()
			line1.SetStartPos(calcLifelineCenterX(msg.From), y)
			line1.SetEndPos(calcLifelineCenterX(msg.From)+w, y)
			line2.SetStartPos(calcLifelineCenterX(msg.From)+w, y)
			line2.SetEndPos(calcLifelineCenterX(msg.From)+w, y+h)
			line3.SetStartPos(calcLifelineCenterX(msg.From)+w, y+h)
			line3.SetEndPos(calcLifelineCenterX(msg.From), y+h)
			line3.SetTailType("triangle")
			ss.AddShape(line1)
			ss.AddShape(line2)
			ss.AddShape(line3)
		}

		if msg.Text != "" {
			var c int
			if msg.From.Index < msg.To.Index {
				c = calcLifelineCenterX(msg.From)
			} else {
				c = calcLifelineCenterX(msg.To)
			}
			textbox := shape.NewRectangle()
			textbox.SetNoFill(true)
			textbox.SetNoLine(true)
			textbox.SetText(msg.Text, "en-US")
			textbox.SetLeftTop(c, y-20)
			textbox.SetSize(spanX, spanY)
			ss.AddShape(textbox)
		}
	}
}

// DrawNotes adds 'Note' into the spreadsheet.
func DrawNotes(ss *oxml.Spreadsheet, notes []*model.Note) {
	baseY := marginY + sizeY
	for _, note := range notes {
		i := note.Assoc.Index
		y := baseY + spanY*(i+1)
		w := maxLine(note.Text) * 8
		h := (len(strings.Split(note.Text, "\n"))+1)*15 + 8

		rect := shape.NewRectangle()
		rect.SetSize(w, h)
		rect.SetText(note.Text, "en-US")
		rect.SetFillColor(note.ColorHex)
		if note.OnLeft {
			rect.SetLeftTop(calcLifelineCenterX(note.Assoc.From)-12-w, y)
		} else {
			rect.SetLeftTop(calcLifelineCenterX(note.Assoc.To)+12, y)
		}
		ss.AddShape(rect)
	}
}

func maxLine(text string) int {
	max := 0
	for _, line := range strings.Split(text, "\n") {
		if len(line) > max {
			max = len(line)
		}
	}
	return max
}
