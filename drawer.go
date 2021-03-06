package seq2xls

import (
	"math"
	"strings"

	"github.com/golang-collections/collections/stack"
	"github.com/rsp9u/go-xlsshape/oxml"
	"github.com/rsp9u/go-xlsshape/oxml/shape"
	"github.com/rsp9u/seq2xls/model"
)

type fragmentReserve struct {
	left, right, top, bottom int
	leftLifeline             *model.Lifeline
	rightLifeline            *model.Lifeline
	body                     *model.Fragment
}

const (
	marginX     = 20
	marginY     = 20
	sizeX       = 6 * 20
	sizeY       = 3 * 20
	spanX       = 192
	spanY       = 40
	tailY       = spanY * 3 / 2
	fragOffsetX = spanX / 3
	fragMarginX = 8
	fragMarginY = 24
	fragGuardX  = 48
	fragGuardY  = 24
)

// DrawSequenceDiagram draws a sequence diagram into the given spreadsheet.
func DrawSequenceDiagram(ss *oxml.Spreadsheet, seq *model.SequenceDiagram) {
	bottom := drawTimeline(ss, seq)
	drawLifelines(ss, seq.Lifelines, bottom)
}

// drawLifelines adds the shapes which composes 'Lifeline' into the spreadsheet.
//
// 'Lifeline' is composed of a rectangle and a dashed line.
func drawLifelines(ss *oxml.Spreadsheet, lls []*model.Lifeline, bottom int) {
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
		line.SetEndPos(rectXCenter, bottom+tailY)
		line.SetDashType("dash")
		ss.UnshiftShape(line)
	}
}

func calcLifelineCenterX(ll *model.Lifeline) int {
	return marginX + spanX*ll.Index + sizeX/2
}

// drawTimeline adds the shapes of the time series elements into the spreadsheet.
func drawTimeline(ss *oxml.Spreadsheet, seq *model.SequenceDiagram) (y int) {
	y = marginY + sizeY + spanY
	fragRsvs := stack.New()
	fragLimitLeft := 0
	fragLimitRight := math.MaxInt32

	for _, sep := range seq.Separators {
		if sep.Before == nil {
			deltaY := drawSeparator(ss, sep, y, len(seq.Lifelines))
			y += deltaY
		}
	}

	for _, msg := range seq.Messages {
		// fragment opening
		for _, frag := range seq.Fragments {
			if frag.Begin == msg {
				leftll, rightll := getBothEndsLifeline(frag, seq.Messages)

				left := calcLifelineCenterX(leftll) - fragOffsetX
				if left <= fragLimitLeft {
					left = fragLimitLeft + fragMarginX
				}
				fragLimitLeft = left

				right := calcLifelineCenterX(rightll) + fragOffsetX
				if right >= fragLimitRight {
					right = fragLimitRight - fragMarginX
				}
				fragLimitRight = right

				fragRsvs.Push(&fragmentReserve{
					top:           y,
					left:          left,
					right:         right,
					leftLifeline:  leftll,
					rightLifeline: rightll,
					body:          frag,
				})
				y += fragMarginY
			}
		}

		// proceed a message
		deltaY := 0
		deltaY += drawMessage(ss, msg, y)
		for _, note := range seq.Notes {
			if note.Assoc == msg {
				deltaY += drawNote(ss, note, y)
			}
		}

		// fragment closing
		for fragRsvs.Len() != 0 {
			frag, ok := fragRsvs.Peek().(*fragmentReserve)
			if !ok || frag.body.End != msg {
				break
			}
			fragRsvs.Pop()
			y += fragMarginY
			frag.bottom = y

			drawFragment(ss, frag)
		}
		if fragRsvs.Len() != 0 {
			frag, ok := fragRsvs.Peek().(*fragmentReserve)
			if ok {
				fragLimitLeft = frag.left
				fragLimitRight = frag.right
			}
		} else {
			fragLimitLeft = 0
			fragLimitRight = math.MaxInt32
		}

		y += deltaY

		for _, sep := range seq.Separators {
			if sep.Before == msg {
				deltaY := drawSeparator(ss, sep, y, len(seq.Lifelines))
				y += deltaY
			}
		}
	}

	return
}

func drawMessage(ss *oxml.Spreadsheet, msg *model.Message, y int) (deltaY int) {
	y += spanY / 2
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

	if msg.Type == model.SelfReference {
		return spanY + spanY/3
	}
	return spanY
}

func drawNote(ss *oxml.Spreadsheet, note *model.Note, y int) (deltaY int) {
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

	return 0
}

func drawFragment(ss *oxml.Spreadsheet, frag *fragmentReserve) {
	rect := shape.NewRectangle()
	rect.SetLeftTop(frag.left, frag.top)
	rect.SetSize(frag.right-frag.left, frag.bottom-frag.top)
	rect.SetNoFill(true)
	rect.SetText(frag.body.Type.String(), "en-US")
	ss.AddShape(rect)

	line1 := shape.NewLine()
	line2 := shape.NewLine()
	line1.SetStartPos(frag.left, frag.top+fragGuardY)
	line1.SetEndPos(frag.left+fragGuardX, frag.top+fragGuardY)
	line2.SetStartPos(frag.left+fragGuardX, frag.top+fragGuardY)
	line2.SetEndPos(frag.left+fragGuardX+12, frag.top)
	ss.AddShape(line1)
	ss.AddShape(line2)
}

func drawSeparator(ss *oxml.Spreadsheet, sep *model.Separator, y, nLls int) (deltaY int) {
	left := marginX
	right := marginX + spanX*(nLls-1) + sizeX
	center := (right-left)/2 + left

	line1 := shape.NewLine()
	line2 := shape.NewLine()
	line1.SetStartPos(left, y+12)
	line1.SetEndPos(right, y+12)
	line2.SetStartPos(left, y+18)
	line2.SetEndPos(right, y+18)
	ss.AddShape(line1)
	ss.AddShape(line2)

	w := len(sep.Text) * 12
	h := 20
	rect := shape.NewRectangle()
	rect.SetLeftTop(center-w/2, y+15-h/2)
	rect.SetSize(w, h)
	rect.SetText(sep.Text, "en-US")
	rect.SetHAlign("ctr")
	rect.SetVAlign("ctr")
	ss.AddShape(rect)

	return 12 + 6 + 12
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

func getBothEndsLifeline(frag *model.Fragment, msgs []*model.Message) (mostLeft, mostRight *model.Lifeline) {
	for i := frag.Begin.Index; i <= frag.End.Index; i++ {
		if mostLeft == nil || msgs[i].From.Index < mostLeft.Index {
			mostLeft = msgs[i].From
		}
		if mostLeft == nil || msgs[i].To.Index < mostLeft.Index {
			mostLeft = msgs[i].To
		}
		if mostRight == nil || msgs[i].From.Index > mostRight.Index {
			mostRight = msgs[i].From
		}
		if mostRight == nil || msgs[i].To.Index > mostRight.Index {
			mostRight = msgs[i].To
		}
	}

	return
}
