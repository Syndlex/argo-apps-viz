package appsofapps

import (
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/syndlex/argo-apps-viz/pkg/model"
)

func AppsOfAppsRenderGraph(applicationSetList *v1alpha1.ApplicationSetList, applicationList *v1alpha1.ApplicationList) *charts.Graph {
	var c = model.Model(applicationList, applicationSetList)
	node, links, categories := renderGraphFromCluster(c)

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

func getGraphNodesFromApps(r *model.Apps, category string) []opts.GraphNode {
	var nodes []opts.GraphNode
	for _, value := range r.Apps {
		nodes = append(nodes, opts.GraphNode{
			Name:       value.Name,
			Category:   category,
			SymbolSize: 12.5,
		})
	}
	return nodes
}

func getGraphLinksFromApps(r *model.Apps) []opts.GraphLink {
	var nodes []opts.GraphLink
	for _, value := range r.Apps {
		nodes = append(nodes, opts.GraphLink{
			Source: value.Name,
			Target: value.ManagedBy,
		})
	}
	return nodes
}

func getGraphNodesFromDomains(d *model.Domains, category string) []opts.GraphNode {
	var nodes []opts.GraphNode
	for _, values := range d.Repos {
		nodes = append(nodes, getGraphNodesFromApps(values, category)...)
	}
	return nodes
}

func getLinksFromDomains(d *model.Domains) []opts.GraphLink {
	var nodes []opts.GraphLink
	for _, values := range d.Repos {
		nodes = append(nodes, getGraphLinksFromApps(values)...)
	}
	return nodes
}

func renderGraphFromCluster(c model.Clusters) ([]opts.GraphNode, []opts.GraphLink, []*opts.GraphCategory) {
	var categories []*opts.GraphCategory
	var nodes []opts.GraphNode
	var links []opts.GraphLink
	for cluster, domains := range c.Domains {
		categories = append(categories, &opts.GraphCategory{Name: cluster})
		nodes = append(nodes, getGraphNodesFromDomains(domains, cluster)...)
		links = append(links, getLinksFromDomains(domains)...)
	}
	return nodes, links, categories
}
