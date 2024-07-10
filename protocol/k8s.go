package protocol

import (
	"fmt"
	"strings"
)

type (
	K8S struct {
		Schema    string
		Domain    string
		Service   string
		Namespace string
		Port      int
	}
)

const (
	defaultDomain    = "svc.cluster.local"
	defaultNamespace = "default"
	defaultPort      = 8080
	defaultSchema    = "http"
)

func (svc *K8S) Formatter() string {
	svc.namespace().domain().port().schema()

	return fmt.Sprintf("%s://%s.%s.%s:%d", svc.Schema, svc.Service, svc.Namespace, svc.Domain, svc.Port)
}

func (svc *K8S) namespace() *K8S {
	if svc.Namespace == "" {
		svc.Namespace = defaultNamespace
	}

	return svc
}

func (svc *K8S) domain() *K8S {
	if svc.Domain == "" {
		svc.Domain = defaultDomain
	}

	return svc
}

func (svc *K8S) port() *K8S {
	if svc.Port <= 0 {
		svc.Port = defaultPort
	}

	return svc
}

func (svc *K8S) schema() *K8S {
	if svc.Schema == "" {
		svc.Schema = defaultSchema
	} else {
		if strings.ToLower(svc.Schema) == "k8ss" {
			svc.Schema = "https"
		} else {
			svc.Schema = "http"
		}
	}

	return svc
}
