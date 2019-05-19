package seq2xls

import (
	"github.com/rsp9u/go-xlsshape/oxml"
	"github.com/rsp9u/go-xlsshape/oxml/shape"
	"github.com/rsp9u/seq2xls/model"
)

func DrawLifelines(ss *oxml.Spreadsheet, lls []*model.Lifeline, nMsg int) {
	const (
		marginX = 20
		marginY = 20
		sizeX   = 6 * 20
		sizeY   = 3 * 20
		spanX   = 192
		spanY   = 40
		tailY   = 60
	)

	for _, ll := range lls {
		i := ll.Index
		rect := shape.NewRectangle()
		rect.SetLeftTop(marginX+spanX*i, marginY)
		rect.SetSize(sizeX, sizeY)
		rect.SetText(ll.Name, "en-US")
		rect.SetHAlign("ctr")
		rect.SetVAlign("ctr")
		ss.AddShape(rect)

		rectXCenter := marginX + spanX*i + sizeX/2
		rectBottom := marginY + sizeY
		line := shape.NewLine()
		line.SetStartPos(rectXCenter, rectBottom)
		line.SetEndPos(rectXCenter, rectBottom+spanY*(nMsg+1)+tailY)
		line.SetDashType("dash")
		ss.AddShape(line)
	}
}
