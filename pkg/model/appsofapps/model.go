package appsofapps

import (
	"argo-apps-viz/pkg/model" // Using local module

	alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1" // Updated to v3
)

func createModel(applist *alpha1.ApplicationList, setList *alpha1.ApplicationSetList) model.Apps {
	var apps model.Apps
	for _, item := range applist.Items {
		apps.AddApps(item)
	}
	for _, item := range setList.Items {
		apps.AddAppsFromSet(item)
	}
	return apps
}
