package dependencies

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/syndlex/argo-apps-viz/pkg/model"
)

func RenderTree(c clusters) *charts.Tree {
	cluster := renderTreeCluster(c)

	tree := charts.NewTree()
	tree.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "100%",
		}),
		charts.WithTitleOpts(opts.Title{Title: "Cluster Dependency Provider"}),
	)
	start := []opts.TreeData{
		{
			Name:     "Root",
			Children: cluster,
		},
	}
	tree.AddSeries("tree", start).
		SetSeriesOptions(
			charts.WithTreeOpts(
				opts.TreeChart{

					Layout:           "orthogonal",
					Orient:           "LR",
					Roam:             opts.Bool(true),
					InitialTreeDepth: 1,
					Leaves: &opts.TreeLeaves{
						Label: &opts.Label{Show: opts.Bool(true), Position: "right", Color: "Black"},
					},
				},
			),
			charts.WithLabelOpts(opts.Label{Show: opts.Bool(true), Position: "top", Color: "Black"}),
		)

	return tree
}

func renderTreeDomains(d *domains) []*opts.TreeData {
	var clusterNodes []*opts.TreeData
	for key, values := range d.Repos {
		clusterNodes = append(clusterNodes, &opts.TreeData{
			Name:     key,
			Children: renderTreeApps(values),
		})
	}
	return clusterNodes
}

func renderTreeCluster(c clusters) []*opts.TreeData {
	var clusterNodes []*opts.TreeData
	for cluster, repo := range c.Domains {
		clusterNodes = append(clusterNodes, &opts.TreeData{
			Name:     cluster,
			Children: renderTreeDomains(repo),
		})
	}
	return clusterNodes
}

func renderTreeApps(r *model.Apps) []*opts.TreeData {
	var clusterNodes []*opts.TreeData
	for _, value := range r.Apps {
		clusterNodes = append(clusterNodes, &opts.TreeData{
			Name: value.Name,
		})
	}
	return clusterNodes
}
