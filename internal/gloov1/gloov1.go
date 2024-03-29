package gloov1

type VirtualService struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Metadata struct {
	Name      string      `yaml:"name"`
	Namespace string      `yaml:"namespace"`
	Labels    interface{} `yaml:"labels,omitempty"`
}

type Spec struct {
	DisplayName string      `yaml:"displayName"`
	VirtualHost VirtualHost `yaml:"virtualHost"`
}

type VirtualHost struct {
	Domains []string            `yaml:"domains"`
	Routes  []Route             `yaml:"routes"`
	Options *VirtualHostOptions `yaml:"options,omitempty"`
}

type CorsPolicy struct {
	AllowCredentials bool     `yaml:"allowCredentials,omitempty"`
	AllowHeaders     []string `yaml:"allowHeaders,omitempty"`
	AllowMethods     []string `yaml:"allowMethods,omitempty"`
	AllowOrigin      []string `yaml:"allowOrigin,omitempty"`
	ExposeHeaders    []string `yaml:"exposeHeaders,omitempty"`
	MaxAge           string   `yaml:"maxAge,omitempty"`
}
type Matchers struct {
	Methods []string `yaml:"methods"`
	Prefix  string   `yaml:"prefix"`
}

type ExtAuthExtension struct {
	CustomAuth CustomAuth `yaml:"customAuth"`
}

type CustomAuth struct {
	ContextExtensions map[string]string `yaml:"contextExtensions,omitempty"`
}

type Route struct {
	Matchers     []Matchers    `yaml:"matchers"`
	RouteAction  RouteAction   `yaml:"routeAction"`
	RouteOptions *RouteOptions `yaml:"options,omitempty"`
}

type RouteAction struct {
	Single RouteSingle `yaml:"single"`
}

type RouteKube struct {
	Ref  ResourceRef `yaml:"ref"`
	Port int         `yaml:"port"`
}
type RouteSingle struct {
	Upstream *ResourceRef `yaml:"upstream,omitempty"`
	Kube     *RouteKube   `yaml:"kube,omitempty"`
}

type ResourceRef struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type RouteOptions struct {
	Extauth            Extauth             `yaml:"extauth,omitempty"`
	Timeout            string              `yaml:"timeout,omitempty"`
	HeaderManipulation *HeaderManipulation `yaml:"headerManipulation,omitempty"`
	Retries            *Retries            `yaml:"retries,omitempty"`
	PrefixRewrite      string              `yaml:"prefixRewrite,omitempty"`
	Cors               CorsPolicy          `yaml:"cors,omitempty"`
}

type Extauth struct {
	Disable bool `yaml:"disable,omitempty"`
}

type HeaderManipulation struct {
	RequestHeadersToAdd *[]RequestHeadersToAdd `yaml:"requestHeadersToAdd"`
}

type RequestHeadersToAdd struct {
	Header Header `yaml:"header"`
	Append *bool  `yaml:"append"`
}

type Header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Retries struct {
	RetryOn       string `yaml:"retryOn"`
	NumRetries    int    `yaml:"numRetries"`
	PerTryTimeout string `yaml:"perTryTimeout"`
}

type VirtualHostOptions struct {
	Extauth   *ExtAuthExtension `yaml:"extauth,omitempty"`
	Cors      *CorsPolicy       `yaml:"cors,omitempty"`
	RateLimit *struct{}         `yaml:"rate-limit,omitempty"`
}
