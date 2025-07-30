package model

import (
	"strings"

	"github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	"github.com/go-echarts/go-echarts/v2/types"
)

type App struct {
	Name      string
	ManagedBy string
	Cluster   string
	Tooltip   types.FuncStr
}

type Apps struct {
	Apps []App
}

func (r *Apps) AddApps(item v1alpha1.Application) {
	tracking := item.Annotations["argocd.argoproj.io/tracking-id"]

	if tracking != "" {
		tracking = strings.Split(tracking, ":")[0]
	}
	if tracking == "" {
		tracking = item.Labels["argocd.argoproj.io/instance"]
	}
	if tracking == "" {
		references := item.OwnerReferences
		if len(references) != 0 {
			tracking = references[0].Name
		} else {
			println(item.GetName(), " has no reference with Label (argocd.argoproj.io/instance), OwnerReferences or Annotation (argocd.argoproj.io/tracking-id)")
		}
	}
	destName := item.Spec.Destination.Name
	//Use Server when name is not used
	if destName == "" {
		destName = item.Spec.Destination.Server
	}

	r.Apps = append(r.Apps, App{
		Name:      item.Name,
		ManagedBy: tracking,
		Cluster:   destName,
		Tooltip: types.FuncStr(
			"Application:<br/>" +
				item.Name + " - " + destName +
				"<br/>Created: " + item.CreationTimestamp.String() +
				"<br/>Project: " + item.Spec.Project),
	})
}

func (r *Apps) AddAppsFromSet(item v1alpha1.ApplicationSet) {
	tracking := item.Annotations["argocd.argoproj.io/tracking-id"]
	if tracking != "" {
		tracking = strings.Split(tracking, ":")[0]
	} else {
		references := item.OwnerReferences
		if len(references) != 0 {
			tracking = references[0].Name
		} else {
			println(item.GetName(), " has no refrence")
		}
	}
	destName := item.Spec.Template.Spec.Destination.Name
	//Use Server when name is not used
	if destName == "" {
		destName = item.Spec.Template.Spec.Destination.Server
	}

	r.Apps = append(r.Apps, App{
		Name:      item.Name,
		ManagedBy: tracking,
		Cluster:   destName,
		Tooltip: types.FuncStr("ApplicationSet:<br/>" +
			item.Name + " - " + destName +
			"<br/>Created: " + item.CreationTimestamp.String()),
	})
}

func (r *Apps) GetRoots(baseNodes []string) *[]App {
	var list []App
	for _, app := range r.Apps {
		if app.ManagedBy == "" {
			list = append(list, app)
		}
		//add base Node if Provided
		for _, baseNode := range baseNodes {
			if app.Name == baseNode {
				list = append(list, app)
			}
		}
	}
	return &list
}

func (r *Apps) GetManagedBy(name string) []App {
	var list []App
	for _, app := range r.Apps {
		if app.ManagedBy == name {
			list = append(list, app)
		}
	}
	return list
}
