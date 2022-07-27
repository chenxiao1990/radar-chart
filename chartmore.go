package chart

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"

	"github.com/golang/freetype"
)

//DrawMoreData 属性详情
type DrawMoreData struct {
	Name          string     //当前类型
	LinklineColor color.RGBA // LinklineColor 属性连接线颜色
	Values        []int
	FaceColor     color.RGBA //面颜色，如果面的alfa为0 那么不绘制面
}

// MoreOptions 多条线绘制参数
type MoreOptions struct {
	// Width 图片宽
	Width int
	// Height 图片高
	Height int

	// BackgroundColor 背景色
	BackgroundColor color.RGBA
	// BacklineColor 背景线颜色
	BacklineColor color.RGBA
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
	//边界线类型
	DataKeys []string
	// Title 标题文字
	Title string
	// DrawDatas 属性
	DrawDatas []DrawMoreData
}

// NewMoreOption 创建一个option 会填充默认参数
func NewMoreOption() MoreOptions {
	return MoreOptions{
		Width:  500,
		Height: 400,
		// BackgroundColor 背景色
		BackgroundColor: color.RGBA{255, 255, 255, 255},
		// BacklineColor 背景线颜色
		BacklineColor: color.RGBA{180, 180, 180, 255},
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
		FontFile: "./msyh.ttf",
		// DataValue 属性边线 从外到内圈
		DataValues: []int{100, 80, 60, 40, 20, 0},
		//各顶点的类型
		DataKeys: []string{"one", "two", "three", "four", "five", "six"},
		// Title 标题文字
		Title: "",
		// DrawDatas 属性
		DrawDatas: []DrawMoreData{
			{Name: "第一条线", LinklineColor: color.RGBA{12, 14, 255, 255}, Values: []int{1, 50, 60, 70, 80, 10}},
			{Name: "第二条线", LinklineColor: color.RGBA{4 * 16, 9*16 + 14, 255, 255}, Values: []int{80, 50, 50, 88, 100, 30}},
		},
	}
}

// DrawMoreRadar 绘制
func DrawMoreRadar(op MoreOptions) *image.RGBA {
	//绘制画布
	imgrect := image.Rect(0, 0, op.Width, op.Height+50)
	rgba := image.NewRGBA(imgrect)
	//绘制底色
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
	//绘制标题
	titlepoint := image.Point{0, op.TitleFontSize}
	drawtext(rgba, op.FontColor, op.TitleFontSize, f, titlepoint.X, titlepoint.Y, op.Title)
	// 绘制六边形的地方
	edge := op.Height - titlepoint.Y - 2*op.FontSize - 10
	yoffset := titlepoint.Y + op.FontSize
	xoffset := op.Width/2 - edge/2

	if len(op.DataKeys) < 3 {
		fmt.Println("不支持3点以下")
		return nil
	}
	//画边界线
	dingdianpts := make([]image.Point, 0)
	for j := 0; j < len(op.DataKeys); j++ {
		//找到6个顶点位置 在直径为width的圆上
		width := edge * op.DataValues[j] / op.DataValues[0]
		pts := make([]image.Point, 0)
		angel := float64(0)
		for i := 0; i < len(op.DataKeys); i++ {
			angel = float64(i) * float64(math.Pi) * 2 / float64(len(op.DataKeys))
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
			wid := getwidth(op.FontSize, f, op.DataKeys[i])

			tmpx -= wid + 5

		}

		drawtext(rgba, op.FontColor, op.FontSize, f, tmpx, tmpy, op.DataKeys[i])
	}

	// 绘制属性得分线
	for _, datas := range op.DrawDatas {
		angel := float64(0)
		pts := make([]image.Point, 0)
		for i := 0; i < len(op.DataKeys); i++ {
			tmpedge := edge * datas.Values[i] / op.DataValues[0]
			angel = float64(i) * float64(math.Pi) * 2 / float64(len(datas.Values))
			newx := math.Sin(angel) * float64(tmpedge/2)
			newy := math.Cos(angel) * float64(tmpedge/2)

			//坐标转换
			newy = float64(-1)*newy + float64(edge/2+yoffset)
			newx += float64(edge/2 + xoffset)
			pts = append(pts, image.Point{int(newx), int(newy)})
		}

		if datas.FaceColor.A > 0 {
			// 绘制面
			drawface(rgba, datas.LinklineColor, datas.FaceColor, pts, 3)
		}

		// 画属性连接线
		for i := 0; i < len(pts); i++ {
			if i == len(pts)-1 {
				drawline(rgba, datas.LinklineColor, pts[0].X, pts[0].Y, pts[len(pts)-1].X, pts[len(pts)-1].Y, 3)
			} else {
				drawline(rgba, datas.LinklineColor, pts[i].X, pts[i].Y, pts[i+1].X, pts[i+1].Y, 3)
			}
		}

	}

	//绘制所以属性代表的意义
	var bzpt = []int{edge / 2, op.Height + 20}
	for _, datas := range op.DrawDatas {
		drawline(rgba, datas.LinklineColor, bzpt[0], bzpt[1], bzpt[0]+40, bzpt[1], 3)
		drawtext(rgba, op.FontColor, op.FontSize, f, bzpt[0]+45, bzpt[1]+5, datas.Name)
		bzpt[0] = bzpt[0] + 150
	}
	return rgba
}
