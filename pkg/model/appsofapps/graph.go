package appsofapps

import (
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/syndlex/argo-apps-viz/pkg/model"
)

func AppsOfAppsRenderGraph(applicationSetList *v1alpha1.ApplicationSetList, applicationList *v1alpha1.ApplicationList, roots []string, stops []string) *charts.Graph {
	var apps = createModel(applicationList, applicationSetList)
	var node []opts.GraphNode
	var links []opts.GraphLink
	var categories []*opts.GraphCategory
	if len(roots) == 0 {
		node, links, categories = renderGraphFromCluster(apps)
	} else {
		node, links, categories = renderGraphFromAppsRoot(apps, roots, stops)
	}
	var legend []string
	for _, category := range categories {
		legend = append(legend, category.Name)
	}

	graph := charts.NewGraph()
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Show ArgoCd Representation of a Cluster"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:    "100%",
			Height:   "100%",
			Renderer: "svg",
		}),
		charts.WithLegendOpts(opts.Legend{Data: legend}),
	)
	graph.AddSeries("graph", node, links,
		charts.WithGraphChartOpts(
			opts.GraphChart{
				Force:              &opts.GraphForce{Repulsion: 200},
				Roam:               opts.Bool(true),
				Draggable:          opts.Bool(true),
				FocusNodeAdjacency: opts.Bool(true),
				Categories:         categories,
			}),
		charts.WithLabelOpts(opts.Label{Show: opts.Bool(true), Position: "top", Color: "Black"}),

		charts.WithLineStyleOpts(opts.LineStyle{
			Width: 2.5,
			Type:  "dotted",
		}),
	)
	return graph
}

// region Render Graph with root
func renderGraphFromAppsRoot(apps model.Apps, roots []string, stops []string) ([]opts.GraphNode, []opts.GraphLink, []*opts.GraphCategory) {
	var categories []*opts.GraphCategory
	var nodes []opts.GraphNode
	var links []opts.GraphLink
	if len(roots) != 0 {
		for _, root := range roots {
			categories = append(categories, getGraphCategoriesFromAppsRoot(apps, root, stops)...)
			nodes = append(nodes, getGraphNodesFromAppsRoot(apps, root, stops)...)
			links = append(links, getGraphLinksFromAppsRoot(apps, root, stops)...)
		}
	}
	return nodes, links, categories
}

func isAppInStops(stops []string, app model.App) bool {
	for _, stop := range stops {
		if app.Name == stop {
			return true
		}
	}
	return false
}

func getGraphNodesFromAppsRoot(apps model.Apps, appName string, stops []string) []opts.GraphNode {
	var nodes []opts.GraphNode
	for _, app := range apps.Apps {
		if isAppInStops(stops, app) {
			continue
		}
		if app.Name == appName {
			nodes = append(nodes, opts.GraphNode{
				Name:     app.Name,
				Category: app.Cluster,
				Tooltip: &opts.Tooltip{
					Formatter: app.Tooltip,
				},
				SymbolSize: 12.5,
			})
		}
		if app.ManagedBy == appName {
			nodes = append(nodes, getGraphNodesFromAppsRoot(apps, app.Name, stops)...)
		}
	}
	nodes = append(nodes)
	return nodes
}

func getGraphLinksFromAppsRoot(apps model.Apps, appName string, stops []string) []opts.GraphLink {
	var nodes []opts.GraphLink
	for _, app := range apps.Apps {
		if isAppInStops(stops, app) {
			continue
		}
		if app.Name == appName {
			nodes = append(nodes, opts.GraphLink{
				Source: app.Name,
				Target: app.ManagedBy,
			})
		}
		if app.ManagedBy == appName {
			nodes = append(nodes, getGraphLinksFromAppsRoot(apps, app.Name, stops)...)
		}
	}
	return nodes
}

func getGraphCategoriesFromAppsRoot(apps model.Apps, appName string, stops []string) []*opts.GraphCategory {
	var nodes []*opts.GraphCategory
	for _, app := range apps.Apps {
		if isAppInStops(stops, app) {
			continue
		}
		if app.Name == appName {
			nodes = append(nodes, &opts.GraphCategory{
				Name: app.Cluster,
			})
		}
		if app.ManagedBy == appName {
			nodes = append(nodes, getGraphCategoriesFromAppsRoot(apps, app.Name, stops)...)
		}
	}
	return nodes
}

//endregion

// region WithoutRoot
func renderGraphFromCluster(apps model.Apps) ([]opts.GraphNode, []opts.GraphLink, []*opts.GraphCategory) {
	var categories []*opts.GraphCategory
	var nodes []opts.GraphNode
	var links []opts.GraphLink
	for _, app := range apps.Apps {
		categories = append(categories, &opts.GraphCategory{
			Name: app.Cluster,
		})
		nodes = append(nodes, opts.GraphNode{
			Name:       app.Name,
			Category:   app.Cluster,
			SymbolSize: 12.5,
			Tooltip: &opts.Tooltip{
				Formatter: app.Tooltip,
			},
		})
		links = append(links, opts.GraphLink{
			Source: app.Name,
			Target: app.ManagedBy,
		})
	}
	return nodes, links, categories
}

//endregion
