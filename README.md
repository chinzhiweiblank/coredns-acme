# acme

## Name
*acme* is a CoreDNS plugin that automates certificate management: issuance and renewal.
This enables DNS over TLS, a protocol for DNS resolution.

## Description

DNS over TLS is currently done through the CoreDNS `tls` plugin. You have to manually create and provide a certificate and key. You also have to do it when the certificate expires.

However, the `acme` plugin automatically creates and renews the certificate for you, using the `ACME` protocol. 

`ACME` uses challenges to prove that you own the domain. One challenge is `DNS01`, which requires adding DNS records on the authoritative nameserver for your domain. CoreDNS, as a nameserver, can resolve this by creating and providing the necessary records to resolve this challenge.

This [post](https://www.thesslstore.com/blog/acme-protocol-what-it-is-and-how-it-works/) provides a detailed explanation of the protocol.

## Configuration Syntax
## Basic
~~~txt
acme {
  domain <DOMAIN>
}
~~~

* `DOMAIN` is the domain name the plugin should be authoritative for, e.g. contoso.com
* Under this configuration, only the **DNS01** challenge will be used for ACME.

## Advanced
~~~txt
acme {
  domain DOMAIN

  # optional parameters
  challenge <CHALLENGE> port <PORT>
}
~~~
You can specify one or more challenges the CA can use to verify your ownership of the domain.
* `CHALLENGE` is the name of the challenge you will use for ACME. There are only two options: `tlsalpn` and `http01`.
* `PORT` is the port number to use for each challenge. Make sure the ports are open and accessible.


## Examples
### Basic
~~~txt
acme {
  domain contoso.com
}
~~~
This will perform ACME for `contoso.com` and use the `DNS01` challenge only.

### Advanced
This configuration:
~~~txt
acme {
  domain example.com

  challenge http01 port 90
  challenge tlsalpn port 8080
}
~~~
This will perform ACME for `example.com` and perform the following challenges:
1. `HTTP01` challenge on port **90**
2. `TLSALPN` challenge on port **8080**
3. `DNS01` challenge

## Installation
1. Clone CoreDNS and add github.com/chinzhiweiblank/coredns-acme into `go.mod`
2. Clone `https://github.com/chinzhiweiblank/coredns-acme`
3. Add `acme:github.com/chinzhiweiblank/coredns-acme` into `plugin.cfg`
4. Execute `go mod edit -replace github.com/chinzhiweiblank/coredns-acme=${PATH_OF_PLUGIN}`. This enables you to build CoreDNS with the `coredns-acme` repository you cloned.

## See Also
1. [Challenge Types](https://letsencrypt.org/docs/challenge-types/)
2. [RFC for ACME](https://datatracker.ietf.org/doc/html/rfc8555/)
