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
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// DefalutFont 默认字体
var DefalutFont *truetype.Font

func init() {
	// 加载字体
	DefalutFont, _ = truetype.Parse(Roboto)
}

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
	fmt.Println("初始画布", op.Width, op.Height)

	imgrect := image.Rect(0, 0, op.Width, op.Height)
	rgba := image.NewRGBA(imgrect)

	fmt.Println("绘制底色", op.BackgroundColor)
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

	fmt.Println("绘制标题", op.Title)
	titlepoint := image.Point{0, op.TitleFontSize}
	drawtext(rgba, op.FontColor, op.TitleFontSize, f, titlepoint.X, titlepoint.Y, op.Title)

	// 绘制六边形的地方
	edge := op.Height - titlepoint.Y - 2*op.FontSize
	yoffset := titlepoint.Y + op.FontSize
	xoffset := op.Width/2 - edge/2

	if len(op.DrawDatas) < 3 {
		fmt.Println("不支持3点以下")
		return nil
	}

	fmt.Println("画边线")
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

		for i := 0; i < len(pts); i++ {
			if i == len(pts)-1 {
				drawline(rgba, op.BacklineColor, pts[0].X, pts[0].Y, pts[len(pts)-1].X, pts[len(pts)-1].Y, 1)
			} else {
				drawline(rgba, op.BacklineColor, pts[i].X, pts[i].Y, pts[i+1].X, pts[i+1].Y, 1)
			}

			if j == 0 {
				drawtext(rgba, op.FontColor, op.FontSize, f, pts[i].X, pts[i].Y, op.DrawDatas[i].Name)
			}

		}
		//标尺文字
		tmpx := pts[len(pts)-1].X/2 + pts[0].X/2
		tmpy := pts[len(pts)-1].Y/2 + pts[0].Y/2
		drawtext(rgba, op.FontColor, op.FontSize, f, tmpx, tmpy, fmt.Sprintf("%d", op.DataValues[j]))
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

	dx := math.Abs(float64(x0 - x1))
	dy := math.Abs(float64(y0 - y1))
	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy

	for {
		rgba.Set(x0, y0, linecolor)

		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
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
	c.SetDPI(72)
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
