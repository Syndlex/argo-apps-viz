package model

import (
	alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

type Clusters struct {
	Domains map[string]*Domains
}

func (c Clusters) AddDomainFromApp(item alpha1.Application) {
	destName := item.Spec.Destination.Name

	domain := c.Domains[destName]
	if domain == nil {
		domain = &Domains{
			Repos: make(map[string]*Apps),
		}
		c.Domains[destName] = domain
	}
	domain.addDomainApp(item)
}

func (c Clusters) AddDomainFromAppSet(item alpha1.ApplicationSet) {
	destName := item.Spec.Template.Spec.Destination.Name

	domain := c.Domains[destName]
	if domain == nil {
		domain = &Domains{
			Repos: make(map[string]*Apps),
		}
		c.Domains[destName] = domain
	}
	domain.addDomainAppSet(item)
}
