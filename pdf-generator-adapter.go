package pdfgen

import (
	"io"
)

// ------ A4PDFTableGenerator -----
type A4PDFTableGenerator struct {}
func (a *A4PDFTableGenerator) Init() {}
func (a *A4PDFTableGenerator) Cleanup() {}
func (a *A4PDFTableGenerator) GetPageSize() (pageSize *Rect) { return PageSizeA4 }
func (a *A4PDFTableGenerator) GetLineAttr() (drawLine bool, width float64, lineType string) {return true, 0.5, "normal"}
func (a *A4PDFTableGenerator) GetFonts() (<-chan *Font) {return nil}
func (a *A4PDFTableGenerator) GetImages() (<-chan *ImageAttr) {return nil}
func (a *A4PDFTableGenerator) GetWriter() io.Writer { return nil }
func (a *A4PDFTableGenerator) GetTitle() (title string, y float64, fontFamily string, fontSize float64) { return }
func (a *A4PDFTableGenerator) GetColumnTitles() (x, y float64, height float64, fontFamily string, fontSize float64, titles []*Title) { return }
func (a *A4PDFTableGenerator) GetRows() (firstY float64, height float64, fontFamily string, fontSize float64, rows <-chan map[string]string) { return }

