# Gloo-vs-upgrader
Converts Gloo 0.x VirtualServices to Gloo 1.0 VirtualServices


## Usage

Use `-k` to convert to Kube Svc destinations (https://docs.solo.io/gloo/latest/gloo_routing/virtual_services/routes/route_destinations/kubernetes_services/) from 0.x Upstreams.

```bash
go build .
./gloo-vs-upgrader examples/*
writing new VirtualService as examples/api-vs2.yaml
writing new VirtualService as examples/auth-vs2.yaml
writing new VirtualService as examples/noauth-vs2.yaml
writing new VirtualService as examples/rewrites-vs2.yaml
```
