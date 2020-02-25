# radar-chart
雷达图 六维图  八维图

![](https://raw.githubusercontent.com/chenxiao1990/radar-chart/master/example/dst.jpg)

# 示例
```
package main

import (
	"image/jpeg"
	"os"

	chart "github.com/chenxiao1990/radar-chart"
)

func main() {

	op := chart.NewOption()
	op.Title = "kjkj"
	// 可以指定字体  这个msyh.ttf是微软雅黑，可以显示中文
	// op.FontFile = "./msyh.ttf"
	// 设置文字颜色
	// op.FontColor = color.RGBA{255, 0, 0, 255}
	// 设置连线颜色
	// op.BacklineColor = color.RGBA{255, 255, 0, 255}
	img := chart.DrawRadar(op)
	file, _ := os.Create("dst.jpg")
	jpeg.Encode(file, img, nil)
	defer file.Close()
}

```
