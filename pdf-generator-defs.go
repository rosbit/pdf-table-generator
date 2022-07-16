package pdfgen

import (
	"github.com/signintech/gopdf"
	"io"
)

type (
	Point = gopdf.Point
	Rect = gopdf.Rect

	ImageAttr struct {
		Point // Left-upper position
		Rect  // width x height
		ImagePath string
	}

	Font struct {
		Family string
		FontPath string
	}
)

type Title struct {
	Name  string
	Width float64
	ColumnValueAlignLeft bool // whether the value of this column left-alignment or not. false to center-align. the title center-alignment always.
}

var (
	PageSizeA4 = gopdf.PageSizeA4
)

type PDFTableGenerator interface {
	// called first.
	Init()
	// called last.
	Cleanup()

	// PageSize. e.g.: PageSizeA4
	GetPageSize() (pageSize *Rect)

	// line attrubites. e.g.: width=0.5, lineType="normal"
	GetLineAttr() (drawLine bool, width float64, lineType string)

	// get font family name & font file name
	GetFonts() (<-chan *Font)

	// get image file path & left-upper pointer and crop size to show
	GetImages() (<-chan *ImageAttr)

	// where to output pdf
	GetWriter() io.Writer

	// get page title, left-bottom Y, font family name and font size. (title will be displayed center-aligned)
	GetTitle() (title string, y float64, fontFamily string, fontSize float64)

	// get table column titles: (x,y) of fisrt column, title height, font family name, font size and title attributes.
	GetColumnTitles() (x, y float64, height float64, fontFamily string, fontSize float64, titles []*Title)

	// get rows with Y of first row, row height, font family name, font size. rows is a channel containing (title name => value)
	GetRows() (firstY float64, height float64, fontFamily string, fontSize float64, rows <-chan map[string]string)
}
