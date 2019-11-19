package gloov1

import (
	"fmt"
	"strconv"
	"strings"

	glooV0 "github.com/leosunmo/gloo-vs-upgrader/internal/gloov0"
)

const (
	// V1GlooAPIVersion is the V1 Gloo VirtualService API version
	V1GlooAPIVersion = "gateway.solo.io/v1"

	// V1GlooKind is the V1 virtual service Kind
	V1GlooKind = "VirtualService"
)

// ConvertVirtualService converts provided VirtualService to a Mapping
func ConvertVirtualService(v0Vs glooV0.VirtualService, kSvc bool) (VirtualService, error) {
	v1Vs := VirtualService{}
	var err error
	err = v1Vs.buildMetadata(v0Vs)
	if err != nil {
		return v1Vs, err
	}
	err = v1Vs.buildVirtualHost(v0Vs)
	if err != nil {
		return v1Vs, err
	}
	err = v1Vs.buildRoutes(v0Vs, kSvc)
	if err != nil {
		return v1Vs, err
	}
	return v1Vs, nil
}

func (v1Vs *VirtualService) buildMetadata(v0Vs glooV0.VirtualService) error {
	v1Vs.APIVersion = V1GlooAPIVersion
	v1Vs.Kind = V1GlooKind
	v1Vs.Metadata = Metadata{
		Name:      v0Vs.Metadata.Name,
		Namespace: v0Vs.Metadata.Namespace,
		Labels:    v0Vs.Metadata.Labels,
	}
	return nil
}

func (v1Vs *VirtualService) buildVirtualHost(v0Vs glooV0.VirtualService) error {
	// Copy the things that haven't changed
	v1Vs.Spec.DisplayName = v0Vs.Spec.DisplayName
	v1Vs.Spec.VirtualHost.Domains = v0Vs.Spec.VirtualHost.Domains

	// VirtualHostPlugins to VirtualHostOptions
	if v0Vs.Spec.VirtualHost.VirtualHostPlugins != nil {
		v1Vs.Spec.VirtualHost.Options = &VirtualHostOptions{}
		v1Opts := v1Vs.Spec.VirtualHost.Options
		// ExtAuth
		if v0Vs.Spec.VirtualHost.VirtualHostPlugins.Extensions.Configs.Extauth != nil {
			// Make empty map of string to string as we might use this setting
			// in the future
			v1Opts.Extauth = &ExtAuthExtension{}
			v1Opts.Extauth.CustomAuth = map[string]string{}
		}
		// Ratelimit, pretty sure this path is Enterprise only now.
		if v0Vs.Spec.VirtualHost.VirtualHostPlugins.Extensions.Configs.RateLimit != nil {
			v1Opts.RateLimit = v0Vs.Spec.VirtualHost.VirtualHostPlugins.Extensions.Configs.RateLimit
		}
	} // VirtualHostOptions

	// Cors
	// Cors is moving from high level VirtualHost to VirtualHostOptions
	if v0Vs.Spec.VirtualHost.CorsPolicy != nil {
		if v1Vs.Spec.VirtualHost.Options == nil {
			v1Vs.Spec.VirtualHost.Options = &VirtualHostOptions{}
		}
		v1Opts := v1Vs.Spec.VirtualHost.Options
		v1Cors := CorsPolicy(*v0Vs.Spec.VirtualHost.CorsPolicy)
		v1Opts.Cors = &v1Cors
	}

	return nil
}

func (v1Vs *VirtualService) buildRoutes(v0Vs glooV0.VirtualService, kubeSvc bool) error {
	var err error
	// Routes
	v1Routes := []Route{}
	for _, v0Route := range v0Vs.Spec.VirtualHost.Routes {
		v1Route := Route{}
		v1Route.Matchers = make([]Matchers, 1)

		// One to one mapping of matchers right now as we don't
		// really route multiple paths to the same upstream very often
		v1Route.Matchers[0].Methods = v0Route.Matcher.Methods
		v1Route.Matchers[0].Prefix = v0Route.Matcher.Prefix

		// RouteAction
		if kubeSvc {
			v1Route.RouteAction, err = convertToKubeRoute(v0Route.RouteAction)
			if err != nil {
				return fmt.Errorf("failed to convert to Kube Svc route, err %s", err)
			}
		} else {
			v1Route.RouteAction, err = convertRouteAction(v0Route.RouteAction)
			if err != nil {
				return fmt.Errorf(err.Error())
			}
		}
		// Use RouteOptions instead of RoutePlugins
		// RoutePlugins/RouteOptions
		if v0Route.RoutePlugins != nil {
			v1Route.RouteOptions = &RouteOptions{}
			// RoutePlugin Extensions
			if v0Route.RoutePlugins.Extensions != nil {
				v1Route.RouteOptions.Extensions = &Extensions{
					Configs: Configs{
						Extauth: Extauth{
							Disable: v0Route.RoutePlugins.Extensions.Configs.Extauth.Disable,
						},
					},
				}
			} // RoutePlugin Extensions

			// Header Manipulation
			if v0Route.RoutePlugins.HeaderManipulation != nil {
				v1rhs := []RequestHeadersToAdd{}
				for _, h := range *v0Route.RoutePlugins.HeaderManipulation.RequestHeadersToAdd {
					v1rh := RequestHeadersToAdd{Append: h.Append, Header: Header(h.Header)}
					v1rhs = append(v1rhs, v1rh)
				}
				v1Route.RouteOptions.HeaderManipulation = &HeaderManipulation{
					RequestHeadersToAdd: &v1rhs,
				}
			} // Header Manipulation

			// Retries
			if v0Route.RoutePlugins.Retries != nil {
				v1Retries := Retries(*v0Route.RoutePlugins.Retries)
				v1Route.RouteOptions.Retries = &v1Retries
			} // Retriues

			// PrefixRewrite
			if v0Route.RoutePlugins.PrefixRewrite != nil {
				v1Route.RouteOptions.PrefixRewrite = v0Route.RoutePlugins.PrefixRewrite.PrefixRewrite
			} // PrefixRewrite

			// Timeout
			v1Route.RouteOptions.Timeout = v0Route.RoutePlugins.Timeout

		} // RoutePlugins/RouteOptions

		v1Routes = append(v1Routes, v1Route)
	} // Routes
	v1Vs.Spec.VirtualHost.Routes = v1Routes
	return nil
}

func convertRouteAction(v0RA glooV0.RouteAction) (RouteAction, error) {
	v1RA := RouteAction{}
	if v0RA.Single.Upstream != nil {
		v1RA.Single.Upstream = &ResourceRef{
			Name:      v0RA.Single.Upstream.Name,
			Namespace: v0RA.Single.Upstream.Namespace,
		}
	} else if v0RA.Single.Kube != nil {
		p, err := strconv.Atoi(v0RA.Single.Kube.Port)
		if err != nil {
			return RouteAction{}, fmt.Errorf("port for %s isn't an int", v0RA.Single.Kube.Ref.Name)
		}
		v1RA.Single.Kube = &RouteKube{
			Port: p,
			Ref:  ResourceRef(v0RA.Single.Kube.Ref),
		}
	} else {
		return v1RA, fmt.Errorf("unable to find any RouteAction")
	}
	return v1RA, nil
}

// convertToKubeRoute takes a RouteAction with an Upstream destination
// and returns a RouteAction with a Kube svc destination
func convertToKubeRoute(v0RA glooV0.RouteAction) (RouteAction, error) {
	us := v0RA.Single.Upstream
	if us == nil {
		// The Upstream is empty, it might already have KubeSvc route,
		// let's try to convert that instead
		v1RA, err := convertRouteAction(v0RA)
		if err != nil {
			return RouteAction{}, fmt.Errorf("no upstream route destination found")
		}
		return v1RA, nil
	}
	var serviceName string
	var servicePort string
	var namespace string
	// Upstream Ref format: {namespace}-{serviceName}-{port}
	// OR if it's too long: {namespace}-{truncServiceName}{randomHash}

	// Find port
	if l := strings.LastIndex(us.Name, "-"); l != -1 {
		serviceName = us.Name[:l]
		servicePort = us.Name[l+1:]
	} else {
		return RouteAction{}, fmt.Errorf("no port found in Gloo Upstream name")
	}
	// Find Namespace through Gloo US and remove it from serviceName
	if l := strings.Index(us.Name, "-"); l != -1 {
		namespace = us.Name[:l]
		serviceName = serviceName[l+1:]
	} else {
		return RouteAction{}, fmt.Errorf("no dash found in Gloo Upstream name")
	}
	p, err := strconv.Atoi(servicePort)
	if err != nil {
		return RouteAction{}, fmt.Errorf("port for %s isn't an int", serviceName)
	}
	if !(p >= 1) || !(p <= 65535) {
		return RouteAction{}, fmt.Errorf("port for %s isn't a valid port number", serviceName)
	}

	var v1RA RouteAction
	v1RA.Single.Kube = &RouteKube{
		Port: p,
		Ref: ResourceRef{
			Name:      serviceName,
			Namespace: namespace,
		},
	}
	return v1RA, nil
}
