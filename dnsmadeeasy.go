package dnsmadeeasy

import (
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"

	dme "github.com/john-k/dnsmadeeasy"
	"github.com/libdns/dnsmadeeasy"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *dnsmadeeasy.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.dnsmadeeasy",
		New: func() caddy.Module { return &Provider{new(dnsmadeeasy.Provider)} },
	}
}

// Before using the provider config, resolve placeholders in the API token.
// Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	p.Provider.APIKey = caddy.NewReplacer().ReplaceAll(p.Provider.APIKey, "")
	p.Provider.SecretKey = caddy.NewReplacer().ReplaceAll(p.Provider.SecretKey, "")
	p.Provider.APIEndpoint = dme.Prod
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//		dnsmadeeasy {
//		    api_key <api_key>
//		    secret_key <secret_key>
//		    api_endpoint https://api.dnsmadeeasy.com/V2.0/
//		}
//
//		api_endpoint is optional, and points to the Production endpoint by default
//	     valid endpoints can be found at:
//		    https://github.com/John-K/dnsmadeeasy/blob/v1.0.0/client.go#L24-L27
//
// Expansion of placeholders is left to the JSON config caddy.Provisioner (above).
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_key":
				if p.Provider.APIKey != "" {
					return d.Err("API key already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.APIKey = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "secret_key":
				if p.Provider.SecretKey != "" {
					return d.Err("secret_key already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.SecretKey = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "api_endpoint":
				// don't error out if this is already set, since
				// we default to Production endpoint as a courtesy
				if !d.NextArg() {
					return d.ArgErr()
				}
				endpoint := dme.BaseURL(d.Val())
				if endpoint == dme.Prod {
					p.Provider.APIEndpoint = dme.Prod
				} else if endpoint == dme.Sandbox {
					p.Provider.APIEndpoint = dme.Sandbox
				} else {
					return d.Err(fmt.Sprintf("Unknown API Endpoint: %s", endpoint))
				}
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.APIKey == "" {
		return d.Err("missing API key")
	}
	if p.Provider.SecretKey == "" {
		return d.Err("missing Secret key")
	}
	if p.Provider.APIEndpoint == "" {
		return d.Err("missing API Endpoint")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
