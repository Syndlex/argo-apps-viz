package model

import (
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"net/url"
)

type Domains struct {
	Repos map[string]*Apps
}

func (d Domains) addDomainApp(item v1alpha1.Application) {
	repoURL := item.Spec.GetSource().RepoURL
	pars, err := url.Parse(repoURL)
	if err != nil {
		panic(err)
	}

	host := pars.Host
	app := d.Repos[host]
	if app == nil {
		app = &Apps{}
		d.Repos[host] = app
	}
	app.AddApps(item)
}

func (d Domains) addDomainAppSet(item v1alpha1.ApplicationSet) {
	repoURL := item.Spec.Template.Spec.GetSource().RepoURL
	pars, err := url.Parse(repoURL)
	if err != nil {
		panic(err)
	}

	host := pars.Host
	app := d.Repos[host]
	if app == nil {
		app = &Apps{}
		d.Repos[host] = app
	}
	app.AddAppsFromSet(item)
}
