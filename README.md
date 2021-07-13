# acme
![CI Workflow](https://github.com/chinzhiweiblank/coredns-acme/actions/workflows/go.yml/badge.svg)

![ACME](img/ACME.png)

*acme* is a plugin that automates certificate management: issuance and renewal.

## Description
Manually creating and renewing certificates can result in certificate mismanagement because:
1. All certificates will expire
2. Users need to be reminded to renew certificates
3. Manpower and time is needed to renew certificates manually
4. The process of creating and renewing certificates [manually](ACME.md#ManagingCertificatesManually) is tedious

Managing certificates manually poses a risk to systems in production because:
1. Users can forget to renew certificate until expiration
2. Risk of exposure leads to gaps in ownership and hence Man-in-the-Middle attacks and breaches.

The `acme` plugin automatically creates and renews certificates for you, using the [ACME]((https://datatracker.ietf.org/doc/html/rfc8555/)) protocol. This enables more secure communications and certificate management while saving time and manpower which could be put to better use.

### Why do you need certificates?
| ![Without SSL](img/HTTP.png) |
|:--:|
| Figure 1: Communication without TLS certificate|

Certificates allow you to encrypted communication between the client and the server so that only the intended recipient can access it. Information you send on the Internet is passed from computer to computer to get to the destination server. In Figure 1, your sensitive information like passwords is not encrypted and can be exposed to any server between you and the recipient.

|![With SSL](img/HTTPS.png)|
|:--:|
|Figure 2: Secure communication with TLS certificate|

In Figure 2, When an SSL/TLS certificate is used, the information becomes unreadable to everyone except for the server you are sending the information to. This protects it from hackers and identity thieves.


## How ACME works
See [ACME.md](ACME.md) for the complete explanation.

## Pros and Cons
### Pros
ACME enables automatic renewal, replacement and revocation of domain validated TLS/SSL certificates.

* You no longer have to spend energy and time to keep a watch on their expiry dates and worry about certificates expiring.
* You no longer have to dig out the instructions to renew and configure the certificates.
* You no longer have to worry about data breaches or Man-in-the-Middle attacks that happen when your certificates expire
* Certificates from LetsEncrypt are free!

Just set up ACME once and let it run. At companies, this could save  a lot of manpower and time when there are hundreds of certificates in use.
### Cons
* LetsEncrypt does not offer OV (Organisation Validation) or EV (Extended Validation) certificates as stated in their [FAQ](https://letsencrypt.org/docs/faq/#will-let-s-encrypt-issue-organization-validation-ov-or-extended-validation-ev-certificates).

## Configuration
## Basic
~~~txt
acme {
  domain <DOMAIN>
}
~~~

* `DOMAIN` is the domain name the plugin should be authoritative for.
* Under this configuration, only the **DNS** challenge will be used for ACME.


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

  challenge http port 90
  challenge tlsalpn port 8080
}
~~~
This will perform ACME for `example.com` and perform the following challenges:
1. `HTTP` challenge on port **90**
2. `TLSALPN` challenge on port **8080**
3. `DNS` challenge

## How this plugin works with CoreDNS
`ACME` uses challenges to prove that you own the domain. One challenge is `DNS`, which requires adding DNS records on the authoritative nameserver for your domain. This plugin uses [CoreDNS](https://github.com/coredns/coredns) to create and providing the necessary records for solving this challenge. It can also resolve the other challenges separately.

## Installation
This is a CoreDNS plugin so you need to set up CoreDNS.
1. Clone [CoreDNS](https://github.com/coredns/coredns) and add github.com/chinzhiweiblank/coredns-acme into `go.mod`
2. Clone `https://github.com/chinzhiweiblank/coredns-acme`
3. Add `acme:github.com/chinzhiweiblank/coredns-acme` into `plugin.cfg`
4. Run `go mod edit -replace github.com/chinzhiweiblank/coredns-acme=${PATH_OF_PLUGIN}`. This enables you to build CoreDNS with the `coredns-acme` repository you cloned.

## Disclaimer
* Make sure you have the following conditions: 
  * You own the domain
  * Your CoreDNS server is the authoritative nameserver for the domain

## See Also
1. [Challenge Types](https://letsencrypt.org/docs/challenge-types/)
2. [RFC for ACME](https://datatracker.ietf.org/doc/html/rfc8555/)
3. [ACME Protocol](https://www.thesslstore.com/blog/acme-protocol-what-it-is-and-how-it-works/)
