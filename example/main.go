package main

import (
	"image/jpeg"
	"os"

	chart "github.com/chenxiao1990/radar-chart"
)

func main() {

	op := chart.NewOption()

	// 可以指定字体  这个msyh.ttf是微软雅黑，可以显示中文
	// op.FontFile = "./msyh.ttf"
	// 设置文字颜色
	// op.FontColor = color.RGBA{255, 0, 0, 255}
	// 设置连线颜色
	// op.BacklineColor = color.RGBA{255, 255, 0, 255}

	op.Title = "kjkj"
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

	// 输出到二进制
	// out := bytes.NewBuffer([]byte{})
	// jpeg.Encode(out, img, nil)
	// fmt.Println(out.Bytes())
}
