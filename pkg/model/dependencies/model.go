package dependencies

import (
	alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

func CreatModel(applist *alpha1.ApplicationList, setList *alpha1.ApplicationSetList) clusters {
	c := clusters{
		Domains: make(map[string]*domains),
	}
	for _, item := range applist.Items {
		c.AddDomainFromApp(item)
	}
	for _, item := range setList.Items {
		c.AddDomainFromAppSet(item)
	}

	return c
}
