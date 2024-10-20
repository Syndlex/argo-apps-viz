package model

import alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"

func Model(applist *alpha1.ApplicationList, setList *alpha1.ApplicationSetList) Clusters {
	c := Clusters{
		Domains: make(map[string]*Domains),
	}
	for _, item := range applist.Items {
		c.AddDomainFromApp(item)
	}
	for _, item := range setList.Items {
		c.AddDomainFromAppSet(item)
	}

	return c
}
