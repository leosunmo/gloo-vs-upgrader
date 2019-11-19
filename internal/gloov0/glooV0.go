package gloov0

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
type CorsPolicy struct {
	AllowCredentials bool     `yaml:"allowCredentials,omitempty"`
	AllowHeaders     []string `yaml:"allowHeaders,omitempty"`
	AllowMethods     []string `yaml:"allowMethods,omitempty"`
	AllowOrigin      []string `yaml:"allowOrigin,omitempty"`
	ExposeHeaders    []string `yaml:"exposeHeaders,omitempty"`
	MaxAge           string   `yaml:"maxAge,omitempty"`
}
type Matcher struct {
	Methods []string `yaml:"methods"`
	Prefix  string   `yaml:"prefix"`
}
type Extauth struct {
	Disable bool `yaml:"disable,omitempty"`
}
type Configs struct {
	Extauth Extauth `yaml:"extauth,omitempty"`
}
type Extensions struct {
	Configs Configs `yaml:"configs,omitempty"`
}
type VHExtauth struct {
	CustomAuth struct{} `yaml:"customAuth,omitempty"`
}
type VHConfigs struct {
	Extauth   *VHExtauth `yaml:"extauth,omitempty"`
	RateLimit *struct{}  `yaml:"rate-limit,omitempty"`
}
type VHExtensions struct {
	Configs VHConfigs `yaml:"configs,omitempty"`
}
type RoutePlugins struct {
	Extensions         *Extensions         `yaml:"extensions,omitempty"`
	Timeout            string              `yaml:"timeout"`
	HeaderManipulation *HeaderManipulation `yaml:"headerManipulation"`
	Retries            *Retries            `yaml:"retries"`
	PrefixRewrite      *PrefixRewrite      `yaml:"prefixRewrite,omitempty"`
}
type ResourceRef struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}
type Kube struct {
	Ref  ResourceRef `yaml:"ref"`
	Port string      `yaml:"port"`
}
type Single struct {
	Upstream *ResourceRef `yaml:"upstream"`
	Kube     *Kube        `yaml:"kube"`
}
type RouteAction struct {
	Single Single `yaml:"single"`
}
type Header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
type RequestHeadersToAdd struct {
	Header Header `yaml:"header"`
	Append *bool  `yaml:"append"`
}
type HeaderManipulation struct {
	RequestHeadersToAdd *[]RequestHeadersToAdd `yaml:"requestHeadersToAdd"`
}
type Retries struct {
	RetryOn       string `yaml:"retryOn"`
	NumRetries    int    `yaml:"numRetries"`
	PerTryTimeout string `yaml:"perTryTimeout"`
}

type PrefixRewrite struct {
	PrefixRewrite string `yaml:"prefixRewrite,omitempty"`
}
type Routes struct {
	Matcher      Matcher       `yaml:"matcher"`
	RoutePlugins *RoutePlugins `yaml:"routePlugins,omitempty"`
	RouteAction  RouteAction   `yaml:"routeAction"`
}
type VirtualHostPlugins struct {
	Extensions VHExtensions `yaml:"extensions,omitempty"`
}
type VirtualHost struct {
	CorsPolicy         *CorsPolicy         `yaml:"corsPolicy,omitempty"`
	Domains            []string            `yaml:"domains"`
	Name               string              `yaml:"name"`
	Routes             []Routes            `yaml:"routes"`
	VirtualHostPlugins *VirtualHostPlugins `yaml:"virtualHostPlugins,omitempty"`
}
type Spec struct {
	DisplayName string      `yaml:"displayName"`
	VirtualHost VirtualHost `yaml:"virtualHost"`
}
