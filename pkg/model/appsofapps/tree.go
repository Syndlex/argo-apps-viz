package appsofapps

import (
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/syndlex/argo-apps-viz/pkg/model"
)

func RenderTree(applicationSetList *v1alpha1.ApplicationSetList, applicationList *v1alpha1.ApplicationList, baseNode []string) *charts.Tree {
	apps := extractModels(applicationSetList, applicationList)
	nodes := renderNodes(apps, baseNode)

	tree := charts.NewTree()
	tree.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "100%",
		}),
		charts.WithTitleOpts(opts.Title{Title: "Cluster Dependency Provider"}),
	)
	tree.AddSeries("tree", nodes).
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

func renderNodes(apps model.Apps, baseNode []string) []opts.TreeData {
	roots := apps.GetRoots(baseNode)
	var treeRoots []opts.TreeData
	for _, root := range *roots {
		treeRoots = append(treeRoots, opts.TreeData{
			Name:     root.Name,
			Children: fillNodes(apps, root.Name, baseNode),
		})
	}
	return treeRoots
}

func fillNodes(apps model.Apps, name string, baseNodes []string) []*opts.TreeData {
	roots := apps.GetManagedBy(name)
	var treeRoots []*opts.TreeData
	for _, root := range roots {
		// Do not add this node because it would lead to a dependency loop
		skip := false
		for _, baseNode := range baseNodes {
			if root.Name == baseNode {
				skip = true
			}
		}
		if skip {
			continue
		}
		treeRoots = append(treeRoots, &opts.TreeData{
			Name:     root.Name,
			Children: fillNodes(apps, root.Name, baseNodes),
		})
	}
	return treeRoots
}

func extractModels(applicationSetList *v1alpha1.ApplicationSetList, applicationList *v1alpha1.ApplicationList) model.Apps {
	var apps model.Apps
	for _, item := range applicationSetList.Items {
		apps.AddAppsFromSet(item)
	}
	for _, item := range applicationList.Items {
		apps.AddApps(item)
	}
	return apps
}
