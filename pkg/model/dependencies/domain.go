package dependencies

import (
	"net/url"

	"argo-apps-viz/pkg/model"

	"github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
)

type domains struct {
	Repos map[string]*model.Apps
}

func (d domains) addDomainApp(item v1alpha1.Application) {
	repoURL := item.Spec.GetSource().RepoURL
	pars, err := url.Parse(repoURL)
	if err != nil {
		panic(err)
	}

	host := pars.Host
	app := d.Repos[host]
	if app == nil {
		app = &model.Apps{}
		d.Repos[host] = app
	}
	app.AddApps(item)
}

func (d domains) addDomainAppSet(item v1alpha1.ApplicationSet) {
	repoURL := item.Spec.Template.Spec.GetSource().RepoURL
	pars, err := url.Parse(repoURL)
	if err != nil {
		panic(err)
	}

	host := pars.Host
	app := d.Repos[host]
	if app == nil {
		app = &model.Apps{}
		d.Repos[host] = app
	}
	app.AddAppsFromSet(item)
}
