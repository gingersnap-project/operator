package context

import "github.com/engytita/engytita-operator/pkg/reconcile"

type Context struct {
	reconcile.Context
	*ServiceBinding
}

type ServiceBinding struct {
	Host     string
	Port     int
	Username string
	Password string
}
