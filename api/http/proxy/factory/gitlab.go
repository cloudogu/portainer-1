package factory

import (
	"net/http"
	"net/url"

	"github.com/cloudogu/portainer-ce/api/http/proxy/factory/gitlab"
)

func newGitlabProxy(uri string) (http.Handler, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	proxy := newSingleHostReverseProxyWithHostHeader(url)
	proxy.Transport = gitlab.NewTransport()
	return proxy, nil
}
