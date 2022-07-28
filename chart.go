package chart

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"math"

	"github.com/golang/freetype"
	"github.com/golang/freetype/raster"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// DefalutFont 默认字体
var DefalutFont *truetype.Font

func init() {
	// 加载字体
	DefalutFont, _ = truetype.Parse(Roboto)
}

// DPI 。。。
var DPI float64 = 92

// DrawData 数据
type DrawData struct {
	Name  string
	Value int
}

// Options 绘制参数
type Options struct {
	// Width 图片宽
	Width int
	// Height 图片高
	Height int

	// BackgroundColor 背景色
	BackgroundColor color.RGBA
	// BacklineColor 背景线颜色
	BacklineColor color.RGBA
	// LinklineColor 属性连接线颜色
	LinklineColor color.RGBA
	// TitleColor 标题文字颜色
	TitleColor color.RGBA
	// TitleFontSize 标题文字大小
	TitleFontSize int
	// FontColor 文字颜色
	FontColor color.RGBA
	// FontSize 文字大小
	FontSize int

	// FontFile 字体文件
	// 例如微软雅黑: D:/haha/msyh.ttf
	FontFile string

	// DataValue 属性边线
	DataValues []int
	// Title 标题文字
	Title string
	// DrawDatas 属性
	DrawDatas []DrawData
}

// NewOption 创建一个option 会填充默认参数
func NewOption() Options {
	return Options{
		Width:  500,
		Height: 400,
		// BackgroundColor 背景色
		BackgroundColor: color.RGBA{255, 255, 255, 255},
		// BacklineColor 背景线颜色
		BacklineColor: color.RGBA{180, 180, 180, 255},
		// LinklineColor 属性连接线颜色
		LinklineColor: color.RGBA{4 * 16, 9*16 + 14, 255, 255},
		// TitleColor 标题文字颜色
		TitleColor: color.RGBA{0, 0, 0, 255},
		// TitleFontSize 标题文字大小
		TitleFontSize: 24,
		// FontColor 文字颜色
		FontColor: color.RGBA{80, 80, 80, 255},
		// FontSize 文字大小
		FontSize: 12,

		// FontFile 字体文件
		// 例如微软雅黑: D:/haha/msyh.ttf
		FontFile: "",
		// DataValue 属性边线 从外到内圈
		DataValues: []int{100, 80, 60, 40, 20, 0},
		// Title 标题文字
		Title: "",
		// DrawDatas 属性
		DrawDatas: []DrawData{
			DrawData{Name: "one", Value: 90},
			DrawData{Name: "two", Value: 74},
			DrawData{Name: "three", Value: 68},
			DrawData{Name: "four", Value: 60},
			DrawData{Name: "five", Value: 77},
			DrawData{Name: "six", Value: 88},
		},
	}
}

// DrawRadar 绘制
func DrawRadar(op Options) *image.RGBA {
	//fmt.Println("初始画布", op.Width, op.Height)

	imgrect := image.Rect(0, 0, op.Width, op.Height)
	rgba := image.NewRGBA(imgrect)

	//fmt.Println("绘制底色", op.BackgroundColor)
	draw.Draw(rgba, imgrect, &image.Uniform{op.BackgroundColor}, image.ZP, draw.Src)

	// 加载字体
	f := DefalutFont
	if op.FontFile != "" {
		fontBytes, err := ioutil.ReadFile(op.FontFile)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		f, err = freetype.ParseFont(fontBytes)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	//fmt.Println("绘制标题", op.Title)
	titlepoint := image.Point{0, op.TitleFontSize}
	drawtext(rgba, op.FontColor, op.TitleFontSize, f, titlepoint.X, titlepoint.Y, op.Title)

	// 绘制六边形的地方
	edge := op.Height - titlepoint.Y - 2*op.FontSize - 10
	yoffset := titlepoint.Y + op.FontSize
	xoffset := op.Width/2 - edge/2

	if len(op.DrawDatas) < 3 {
		fmt.Println("不支持3点以下")
		return nil
	}

	//fmt.Println("画边线")
	dingdianpts := make([]image.Point, 0)
	for j := 0; j < len(op.DataValues); j++ {
		//找到6个顶点位置 在直径为width的圆上
		width := edge * op.DataValues[j] / op.DataValues[0]
		pts := make([]image.Point, 0)
		angel := float64(0)
		for i := 0; i < len(op.DrawDatas); i++ {
			angel = float64(i) * float64(math.Pi) * 2 / float64(len(op.DrawDatas))
			newx := math.Sin(angel) * float64(width/2)
			newy := math.Cos(angel) * float64(width/2)

			//坐标转换
			newy = float64(-1)*newy + float64(edge/2+yoffset)
			newx += float64(edge/2 + xoffset)
			pts = append(pts, image.Point{int(newx), int(newy)})
		}
		if j == 0 {
			dingdianpts = pts
		}
		for i := 0; i < len(pts); i++ {
			if i == len(pts)-1 {
				drawline(rgba, op.BacklineColor, pts[0].X, pts[0].Y, pts[len(pts)-1].X, pts[len(pts)-1].Y, 1)
			} else {
				drawline(rgba, op.BacklineColor, pts[i].X, pts[i].Y, pts[i+1].X, pts[i+1].Y, 1)
			}

		}
		//标尺文字
		tmpx := pts[len(pts)-1].X/2 + pts[0].X/2
		tmpy := pts[len(pts)-1].Y/2 + pts[0].Y/2
		drawtext(rgba, op.FontColor, op.FontSize, f, tmpx, tmpy, fmt.Sprintf("%d", op.DataValues[j]))
	}

	// 绘制顶点的文字
	centerx := edge/2 + xoffset
	centery := edge/2 + yoffset
	for i, pt := range dingdianpts {
		tmpy := pt.Y
		tmpx := pt.X
		if centerx == tmpx && tmpy > centery {
			//往下写
			tmpy += op.FontSize
		} else if tmpx > centerx {
			//右侧
			tmpx += 5
		} else if tmpx < centerx {
			//左侧
			wid := getwidth(op.FontSize, f, op.DrawDatas[i].Name)

			tmpx -= wid + 5

		}

		drawtext(rgba, op.FontColor, op.FontSize, f, tmpx, tmpy, op.DrawDatas[i].Name)
	}

	// 绘制属性得分线
	angel := float64(0)
	pts := make([]image.Point, 0)
	for i := 0; i < len(op.DrawDatas); i++ {
		tmpedge := edge * op.DrawDatas[i].Value / op.DataValues[0]
		angel = float64(i) * float64(math.Pi) * 2 / float64(len(op.DrawDatas))
		newx := math.Sin(angel) * float64(tmpedge/2)
		newy := math.Cos(angel) * float64(tmpedge/2)

		//坐标转换
		newy = float64(-1)*newy + float64(edge/2+yoffset)
		newx += float64(edge/2 + xoffset)
		pts = append(pts, image.Point{int(newx), int(newy)})
	}

	// 画属性连接线
	for i := 0; i < len(pts); i++ {
		if i == len(pts)-1 {
			drawline(rgba, op.LinklineColor, pts[0].X, pts[0].Y, pts[len(pts)-1].X, pts[len(pts)-1].Y, 3)
		} else {
			drawline(rgba, op.LinklineColor, pts[i].X, pts[i].Y, pts[i+1].X, pts[i+1].Y, 3)
		}
	}

	return rgba
}

func drawline(rgba *image.RGBA, linecolor color.RGBA, x0, y0, x1, y1 int, linewidth int) {

	r := raster.NewRasterizer(rgba.Rect.Dx(), rgba.Rect.Dy())
	r.UseNonZeroWinding = true

	c := raster.RoundCapper
	j := raster.RoundJoiner

	var path raster.Path
	path.Start(fixed.P(x0, y0))
	path.Add1(fixed.P(x1, y1))
	raster.Stroke(r, path, fixed.I(linewidth), c, j)

	p := raster.NewRGBAPainter(rgba)
	p.SetColor(linecolor)
	r.Rasterize(p)

}

// 多边形蒙版
type DuobianImage struct {
	Pts    []image.Point
	Center image.Point
}

func (s *DuobianImage) At(x, y int) color.Color {

	for i := 0; i < len(s.Pts); i++ {
		var tmp *SanjiaoImage
		if i == len(s.Pts)-1 {
			tmp = &SanjiaoImage{P1: s.Center, P2: s.Pts[i], P3: s.Pts[0]}
		} else {
			tmp = &SanjiaoImage{P1: s.Center, P2: s.Pts[i], P3: s.Pts[1+i]}
		}

		//蒙版 如果x,y在三角形内返回 alfa 255  否则返回alfa 0
		_, _, _, a := tmp.At(x, y).RGBA()
		if a == 0xffff {
			return color.Alpha{255}
		}
	}

	return color.Alpha{}
}

// 三角面蒙版
type SanjiaoImage struct {
	P1 image.Point
	P2 image.Point
	P3 image.Point
}

func (s *SanjiaoImage) At(x, y int) color.Color {
	//蒙版 如果x,y在三角形内返回 alfa 255  否则返回alfa 0
	if isInTriangle(Point{s.P1.X, s.P1.Y}, Point{s.P2.X, s.P2.Y}, Point{s.P3.X, s.P3.Y}, Point{x, y}) {
		return color.Alpha{A: 255}
	}

	return color.Alpha{}
}

type Point struct {
	x int
	y int
}

func product(p1 Point, p2 Point, p3 Point) int {
	//首先根据坐标计算p1p2和p1p3的向量，然后再计算叉乘
	//p1p2 向量表示为 (p2.x-p1.x,p2.y-p1.y)
	//p1p3 向量表示为 (p3.x-p1.x,p3.y-p1.y)
	return (p2.x-p1.x)*(p3.y-p1.y) - (p2.y-p1.y)*(p3.x-p1.x)
}
func isonline(p1 Point, p2 Point, p3 Point) bool {
	tmp := (p2.x-p1.x)*(p3.y-p1.y) - (p2.y-p1.y)*(p3.x-p1.x)
	if tmp == 0 && ((p3.x >= p1.x && p3.x <= p2.x) || (p3.x <= p1.x && p3.x >= p2.x)) && ((p3.y >= p1.y && p3.y <= p2.y) || (p3.y <= p1.y && p3.y >= p2.y)) {
		return true
	}
	return false
}
func isInTriangle(p1, p2, p3, o Point) bool {
	//保证p1，p2，p3是逆时针顺序
	if product(p1, p2, p3) < 0 {
		return isInTriangle(p1, p3, p2, o)
	}
	if product(p1, p2, o) > 0 && product(p2, p3, o) > 0 && product(p3, p1, o) > 0 {
		return true
	}
	if isonline(p1, p2, o) || isonline(p2, p3, o) || isonline(p3, p1, o) {
		return true
	}
	return false
}

// 绘制三角面
func drawface(rgba *image.RGBA, linecolor, bgcolor color.RGBA, center image.Point, pts []image.Point, linewidth int) {

	mask := &DuobianImage{Pts: pts, Center: center}

	for x := 0; x < rgba.Rect.Dx(); x++ {
		for y := 0; y < rgba.Rect.Dy(); y++ {

			_, _, _, a := mask.At(x, y).RGBA()
			if a > 0 {
				// 混合当前颜色和背景色
				cr, cg, cb, _ := rgba.At(x, y).RGBA()
				nr, ng, nb, na := bgcolor.RGBA()

				r := uint8(float32(nr)/65535*float32(na)/65535*255 + float32(cr)/65535*float32(65535-na)/65535*255)
				g := uint8(float32(ng)/65535*float32(na)/65535*255 + float32(cg)/65535*float32(65535-na)/65535*255)
				b := uint8(float32(nb)/65535*float32(na)/65535*255 + float32(cb)/65535*float32(65535-na)/65535*255)

				rgba.Set(x, y, color.RGBA{r, g, b, 255})
			}
		}
	}

}

// x ,y 是文字 左下角的坐标
func drawtext(rgba *image.RGBA,
	fontcolor color.RGBA,
	fontsize int,
	f *truetype.Font,
	x, y int, text string) {

	c := freetype.NewContext()
	c.SetDPI(DPI)
	c.SetFont(f)
	c.SetFontSize(float64(fontsize))
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(&image.Uniform{fontcolor})

	c.SetHinting(font.HintingNone)

	pt := freetype.Pt(x, y)
	_, err := c.DrawString(text, pt)
	if err != nil {
		log.Println(err)
		return
	}
}

func getwidth(fontsize int, f *truetype.Font, text string) int {
	rgba := image.NewRGBA(image.Rect(0, 0, 1000, 200))
	c := freetype.NewContext()
	c.SetDPI(DPI)
	c.SetFont(f)
	c.SetFontSize(float64(fontsize))
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(&image.Uniform{color.RGBA{0, 0, 0, 255}})

	c.SetHinting(font.HintingNone)

	pt := freetype.Pt(0, 0)
	xx, err := c.DrawString(text, pt)
	if err != nil {
		log.Println(err)
		return 0
	}

	var bb = xx.X >> 6
	return int(bb)
}
