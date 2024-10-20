package cli

import (
	"github.com/go-echarts/go-echarts/v2/components"
	"io"
	"os"
)

func CreateFile(filename string, charter ...components.Charter) error {
	page := components.NewPage()
	page.SetLayout(components.PageFullLayout)
	page.SetPageTitle("Cluster Representation")
	page.AddCharts(charter...)
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	return page.Render(io.MultiWriter(f))
}
