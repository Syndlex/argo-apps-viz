package dependencies

import (
	model "argo-apps-viz/pkg/model"

	alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
)

type clusters struct {
	Domains map[string]*domains
}

func (c clusters) AddDomainFromApp(item alpha1.Application) {
	destName := item.Spec.Destination.Name
	//Use Server when name is not used
	if destName == "" {
		destName = item.Spec.Destination.Server
	}

	domain := c.Domains[destName]
	if domain == nil {
		domain = &domains{
			Repos: make(map[string]*model.Apps),
		}
		c.Domains[destName] = domain
	}
	domain.addDomainApp(item)
}

func (c clusters) AddDomainFromAppSet(item alpha1.ApplicationSet) {
	destName := item.Spec.Template.Spec.Destination.Name
	//Use Server when name is not used
	if destName == "" {
		destName = item.Spec.Template.Spec.Destination.Server
	}

	domain := c.Domains[destName]
	if domain == nil {
		domain = &domains{
			Repos: make(map[string]*model.Apps),
		}
		c.Domains[destName] = domain
	}
	domain.addDomainAppSet(item)
}
