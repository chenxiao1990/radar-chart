package main

import (
	"testing"

	chart "github.com/chenxiao1990/radar-chart"
)

func Benchmark_Draw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		op := chart.NewOption()

		op.Width = 1024
		op.Height = 768
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
		chart.DrawRadar(op)
	}
}
