DNS Made Easy module for Caddy
==============================

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records with [DNS Made Easy](https://dnsmadeeasy.com)

## Caddy module name

```
dns.providers.dnsmadeeasy
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
	"module": "acme",
	"challenges": {
		"dns": {
			"provider": {
				"name": "dnsmadeeasy",
				"api_key": "<DNSMADEEASY_API_KEY>"
				"secret_key": "<DNSMADEEASY_SECRET_KEY>"
				"api_endpoint": "https://<API_ENDPOINT>"
			}
		}
	}
}
```

or with the Caddyfile:

```
# globally
{
	acme_dns dnsmadeeasy {
		api_key {env.DNSMADEEASY_API_KEY}
		secret_key {env.DNSMADEEASY_SECRET_KEY}
		api_endpoint {env.DNSMADEEASY_API_ENDPOINT}
	}
}
```

`api_endpoint` is optional and will default to the [Production endpoint](https://github.com/John-K/dnsmadeeasy/blob/v1.0.0/client.go#L26)
