package protocol

import (
	"fmt"
)

type (
	Domain struct {
		Schema string
		Domain string
		Port   int
	}
)

func (svc *Domain) Formatter() string {
	svc.getDomain().getPort()

	url := fmt.Sprintf("%s://%s:%d", svc.Schema, svc.Domain, svc.Port)
	if svc.Port <= 0 {
		url = fmt.Sprintf("%s://%s", svc.Schema, svc.Domain)
	}

	return url
}

func (svc *Domain) getDomain() *Domain {
	if svc.Domain == "" {
		svc.Domain = defaultDomain
	}

	return svc
}

func (svc *Domain) getPort() *Domain {
	if svc.Port <= 0 {
		svc.Port = defaultPort
	}

	return svc
}
