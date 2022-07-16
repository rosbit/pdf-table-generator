package pdfgen

import (
	"github.com/signintech/gopdf"
	"fmt"
	"log"
)

func GeneratePDFTable(pg PDFTableGenerator) {
	pg.Init()
	defer pg.Cleanup()

	rowsHandled := false
	rowFirstY, rowHeight, rowFontFamily, rowFontSize, rows := pg.GetRows()
	if rows == nil {
		return
	}

	defer func() {
		if rowsHandled {
			return
		}
		// to make sure all the rows will be consumed
		for _ = range rows {}
	}()

	writer := pg.GetWriter()
	if writer == nil {
		return
	}

	titleX, titleY, titleHeight, titleFontFamily, titleFontSize, titles := pg.GetColumnTitles()
	if len(titles) == 0 {
		return
	}

	pageSize := pg.GetPageSize()
	if pageSize == nil {
		pageSize = PageSizeA4
	}

	pdf := &gopdf.GoPdf{}
	defer pdf.Close()
	pdf.Start(gopdf.Config{PageSize: *pageSize})

	drawLine, lineWidth, lineType := pg.GetLineAttr()
	if drawLine {
		pdf.SetLineWidth(lineWidth)
		pdf.SetLineType(lineType)
	}
	addFonts(pdf, pg)

	rowsHandled = true
	bottomMargin := pageSize.H - rowHeight * 1.5

	newPage := func(page int) (rowX, rowY float64) {
		pdf.AddPage()
		drawImages(pdf, pg)
		outputTitle(pdf, pg, pageSize)
		outputColumnTitles(pdf, titleX, titleY, titleHeight, titleFontFamily, titleFontSize, titles)
		if drawLine {
			drawLines(pdf, pageSize, titleX, rowFirstY, rowHeight, rowFontSize, titleHeight, bottomMargin, titles)
		}

		pdf.SetFont(rowFontFamily, "", rowFontSize/2)
		drawPageNo(pdf, pageSize, page, bottomMargin+rowHeight)
		pdf.SetFont(rowFontFamily, "", rowFontSize)
		return titleX, rowFirstY
	}

	addNewPage, page := true, 1
	var rowX, rowY float64
	for row := range rows {
		if addNewPage {
			addNewPage = false
			rowX, rowY = newPage(page)
		}

		outputRow(pdf, row, titles, rowX, rowY)
		rowY += rowHeight

		if rowY >= bottomMargin {
			addNewPage, page = true, page+1
		}
	}

	pdf.Write(writer)
}

func drawLines(pdf *gopdf.GoPdf, pageSize *Rect, firstX, firstY, rowHeight, fontSize, titleHeight, bottomMargin float64, titles []*Title) {
	rowX, rowY := firstX, firstY - (rowHeight - fontSize) / 2.0
	rowEndX := pageSize.W - firstX
	titleLineY := rowY - titleHeight
	pdf.Line(rowX, titleLineY, rowEndX, titleLineY)
	for rowY < bottomMargin {
		pdf.Line(rowX, rowY, rowEndX, rowY)
		rowY += rowHeight
	}
	pdf.Line(rowX, rowY, rowEndX, rowY)

	rowEndY := rowY
	rowY = titleLineY
	pdf.Line(rowX, rowY, rowX, rowEndY)
	lastTitleIdx := len(titles) - 1
	for i, title := range titles {
		if i == lastTitleIdx {
			rowX = rowEndX
		} else {
			rowX += title.Width
			if rowX > rowEndX {
				rowX = rowEndX
			}
		}
		pdf.Line(rowX, rowY, rowX, rowEndY)
	}
}

func addFonts(pdf *gopdf.GoPdf, pg PDFTableGenerator) {
	fonts := pg.GetFonts()
	if fonts != nil {
		for font := range fonts {
			if font == nil {
				continue
			}
			err := pdf.AddTTFFont(font.Family, font.FontPath)
			if err != nil {
				log.Printf("failed to read font %s: %v\n", font.Family, err)
			}
		}
	}
}

func drawImages(pdf *gopdf.GoPdf, pg PDFTableGenerator) {
	images := pg.GetImages()
	if images != nil {
		for image := range images {
			if image == nil {
				continue
			}
			imgPath, pos, rect := image.ImagePath, &image.Point, &image.Rect
			err := pdf.Image(imgPath, pos.X, pos.Y, rect)
			if err != nil {
				log.Printf("failed to draw image %s: %v\n", imgPath, err)
			}
		}
	}
}

func outputTitle(pdf *gopdf.GoPdf, pg PDFTableGenerator, pageSize *Rect) {
	title, y, fontFamily, fontSize := pg.GetTitle()
	if len(title) == 0 {
		return
	}
	pdf.SetFont(fontFamily, "", fontSize)
	pdf.SetY(y)
	textWidth, _ := pdf.MeasureTextWidth(title)
	realX := (pageSize.W - textWidth) / 2.0
	pdf.SetX(realX)
	pdf.Text(title)
}

func drawPageNo(pdf *gopdf.GoPdf, pageSize *Rect, page int, bottomMargin float64) {
	p := fmt.Sprintf("%d", page)
	pdf.SetY(bottomMargin+5)
	textWidth, _ := pdf.MeasureTextWidth(p)
	realX := (pageSize.W - textWidth) / 2.0
	pdf.SetX(realX)
	pdf.Text(p)
}

func outputColumnTitles(pdf *gopdf.GoPdf, x, y float64, height float64, fontFamily string, fontSize float64, titles []*Title) {
	pdf.SetFont(fontFamily, "", fontSize)
	pdf.SetY(y)
	titleX := x
	for i, _ := range titles {
		title := titles[i]
		textWidth, _ := pdf.MeasureTextWidth(title.Name)
		realX := titleX + (title.Width - textWidth) / 2.0
		pdf.SetX(realX)
		pdf.Cell(nil, title.Name)

		titleX += title.Width
	}
}

func outputRow(pdf *gopdf.GoPdf, row map[string]string, titles []*Title, x, y float64) {
	pdf.SetY(y)
	rowX := x

	var realX float64
	for i, _ := range titles {
		title := titles[i]
		if col, ok := row[title.Name]; ok {
			// output column
			if len(col) > 0 {
				if title.ColumnValueAlignLeft {
					realX = rowX + 1.5
				} else {
					textWidth, _ := pdf.MeasureTextWidth(col)
					realX = rowX + (title.Width - textWidth) / 2.0
				}
				pdf.SetX(realX)
				pdf.Cell(nil, col)
			}
		}
		rowX += title.Width
	}
}
