package pdfgen

import (
	"testing"
	"os"
	"io"
	"fmt"
	"log"
)

func TestPDFAdapter(t *testing.T) {
	pg := &pdfTest{}
	GeneratePDFTable(pg)
	log.Printf("testing TestPDFAdapter done!\n")
}

// --------- PDFTableGenerator implementation -------
type pdfTest struct {
	A4PDFTableGenerator
	o *os.File
}

func (t *pdfTest) Init() {
	o, e := os.OpenFile("./a.pdf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if e != nil {
		log.Printf("%v\n", e)
		return
	}
	t.o = o
}

func (t *pdfTest) Cleanup() {
	if t.o != nil {
		t.o.Close()
	}
}

func (t *pdfTest) GetFonts() (<-chan *Font) {
	c := make(chan *Font)
	go func() {
		c <- &Font{
			Family: "msyh",
			FontPath: "./msyh.ttf",
		}
		close(c)
	}()

	return c
}

func (t *pdfTest) GetImages() (<-chan *ImageAttr) {
	c := make(chan *ImageAttr)
	go func(){
		c <- &ImageAttr{
			Point: Point{
				X: 10,
				Y: 10,
			},
			Rect: Rect{
				W: 40,
				H: 40,
			},
			ImagePath: "./logo.png",
		}
		close(c)
	}()

	return c
}

func (t *pdfTest) GetWriter() io.Writer {
	return t.o
}

func (t *pdfTest) GetTitle() (title string, y float64, fontFamily string, fontSize float64) {
	title = "看看pdf标题是否居中"
	y = 40
	fontFamily = "msyh"
	fontSize = 20.0
	return
}

func (t *pdfTest) GetColumnTitles() (x, y float64, height float64, fontFamily string, fontSize float64, titles []*Title) {
	x, y = 30.0, 70.0
	height = 30.0
	fontFamily = "msyh"
	fontSize = 12.0
	titles = []*Title{
		&Title{
			Name: "序号",
			Width: 50,
		},
		&Title{
			Name: "封号",
			Width: 70,
		},
		&Title{
			Name: "外号",
			Width: 70,
		},
		&Title{
			Name: "姓名",
			Width: 40,
		},
		&Title{
			Name: "性别",
			Width: 40,
		},
		&Title{
			Name: "签名",
			Width: 260,
			ColumnValueAlignLeft: true,
		},
	}

	return
}

var data = []map[string]string{
	map[string]string{
		"封号": "千手斗罗",
		"外号": "千手斗罗",
		"姓名": "唐三",
		"性别": "男",
		"签名": "过去别再遗憾；未来无须忧虑；现在加倍珍惜。",
	},
	map[string]string{
		"封号": "柔骨斗罗",
		"外号": "柔骨斗罗",
		"姓名": "小舞",
		"性别": "女",
		"签名": "不乱于心，不困于情。不畏将来，不念过往。",
	},
}

func (t *pdfTest) GetRows() (firstY float64, height float64, fontFamily string, fontSize float64, rows <-chan map[string]string) {
	firstY = 100
	height = 29.0
	fontFamily = "msyh"
	fontSize = 12.0

	c := make(chan map[string]string)
	go func() {
		j := 0
		dataCount := len(data)
		var src map[string]string

		for i:=0; i<20; i++ {
			d := map[string]string{}

			if j >= dataCount {
				j = 0
			}
			src = data[j]
			j += 1

			for k, v := range src {
				d[k] = v
			}
			d["序号"] = fmt.Sprintf("%d", i+1)

			c <- d
		}
		close(c)
	}()

	rows = c
	return
}

