package main

import (
	"image/color"
	"image/jpeg"
	"os"

	chart "github.com/chenxiao1990/radar-chart"
)

func main() {

	more()
}

func normal() {
	op := chart.NewOption()

	// 可以指定字体  这个msyh.ttf是微软雅黑，可以显示中文
	// op.FontFile = "./msyh.ttf"
	// 设置文字颜色
	// op.FontColor = color.RGBA{255, 0, 0, 255}
	// 设置连线颜色
	// op.BacklineColor = color.RGBA{255, 255, 0, 255}

	op.Width = 1024
	op.Height = 600
	op.Title = "Title"
	op.DataValues = []int{100, 80, 60, 40, 20, 0}
	op.DrawDatas = []chart.DrawData{
		chart.DrawData{Name: "one", Value: 90},
		chart.DrawData{Name: "two", Value: 74},
		chart.DrawData{Name: "three", Value: 68},
		chart.DrawData{Name: "four", Value: 60},
		chart.DrawData{Name: "five", Value: 77},
		chart.DrawData{Name: "six", Value: 88},
	}
	img := chart.DrawRadar(op)

	// 输出jpg图片
	file, _ := os.Create("dst.jpg")
	jpeg.Encode(file, img, nil)
	defer file.Close()

	// file, _ := os.Create("dst.png")
	// png.Encode(file, img)
	// defer file.Close()

	// 输出到二进制
	// out := bytes.NewBuffer([]byte{})
	// jpeg.Encode(out, img, nil)
	// fmt.Println(out.Bytes())
}

func more() {
	op := chart.NewMoreOption()

	// 可以指定字体  这个msyh.ttf是微软雅黑，可以显示中文
	// op.FontFile = "./msyh.ttf"
	// 设置文字颜色
	// op.FontColor = color.RGBA{255, 0, 0, 255}
	// 设置连线颜色
	// op.BacklineColor = color.RGBA{255, 255, 0, 255}

	op.Width = 1024
	op.Height = 600
	op.Title = "窝气"
	op.DataValues = []int{100, 80, 60, 40, 20, 0}
	op.DataKeys = []string{"身高", "体重", "BMI", "肺活量", "长跑", "无语", "解决"}
	op.DrawDatas = []chart.DrawMoreData{

		{Name: "班级平均值", FaceColor: color.RGBA{12, 14, 255, 100}, LinklineColor: color.RGBA{12, 14, 255, 255}, Values: []int{80, 66, 77, 44, 76, 100, 88}},
		{Name: "年级平均值", FaceColor: color.RGBA{4 * 16, 9*16 + 14, 255, 100}, LinklineColor: color.RGBA{4 * 16, 9*16 + 14, 255, 255}, Values: []int{58, 0, 92, 70, 76, 57, 66}},
		{Name: "我的成绩", LinklineColor: color.RGBA{34, 139, 34, 255}, Values: []int{98, 85, 92, 80, 80, 60, 55}},
	}
	img := chart.DrawMoreRadar(op)

	// 输出jpg图片
	file, _ := os.Create("dst1.jpg")
	jpeg.Encode(file, img, nil)
	defer file.Close()
}
