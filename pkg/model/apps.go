package model

import (
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"strings"
)

type App struct {
	Name      string
	ManagedBy string
}

type Apps struct {
	Apps []App
}

func (r *Apps) AddApps(item v1alpha1.Application) {
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
	r.Apps = append(r.Apps, App{
		Name:      item.Name,
		ManagedBy: tracking,
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
	r.Apps = append(r.Apps, App{
		Name:      item.Name,
		ManagedBy: tracking,
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
