# How does ACME work?
In the beginning, the client needs to register an account with a CA and add the domain under it. Then it needs to prove that it owns the domain through domain validation.

## Domain Validation
![Domain Validation](img/DomainValidation.png)
### Steps:
1. Client first generates a public and private key. Client keeps the private key.
2. The CA issues one or more challenges (DNS/HTTPS/TLS-ALPN) to prove that the client controls the domain.
3. CA also sends a nonce to sign with the private key. This proves that the client controls the public and private keys.
4. Client fulfills the challenge and signs the provided nonce.
5. LetsEncrypt verifies the nonce and checks whether the challenge is fulfilled.
6. Server is authorised to do certificate management for the domain with the key-value pair. The key-value pair is now known as the **authorised** key-value pair.

### ACME Challenges
These challenges are for proving to the CA that the client owns the domain. In this case, we refer to the client as the one requesting for the certificate and the server as the Certificate Authority.
1. [HTTP](https://datatracker.ietf.org/doc/html/rfc8555#section-8.3)
* Client constructs a key authorization from the token in the challenge and the client's account key. 
* Client then provisions it as a resource on the HTTP server for the domain and notifies the server. The key authorization will be placed at **http://{domain}/.well-known/acme-challenge/{token}**.
* Server will try to retrieve the key authorization from the URL and verify its value matches.
2. [DNS](https://datatracker.ietf.org/doc/html/rfc8555#section-8.4)
* Client constructs a key authorization from the token in the challenge and the client's account key. It computes the SHA256 digest of it.
* Client provisions a TXT record with the digest under **_acme-challenge.{domain}**, the validation domain. Client notifies the server.
* Server will try to retrieve the TXT record under the validation domain name and verify its value matches.
3. [TLS-ALPN](https://datatracker.ietf.org/doc/html/rfc8737)
* Known as TLS with Application-Layer Protocol Negotiation (TLS ALPN). It allows clients to negotiate what protocol to use for communication (at the application level).
* Client configures a TLS server to respond to specific
connection attempts using the ALPN extension with identifying
information.
* Server validates control of the domain name by connecting to a TLS server at one of the addresses resolved for the domain name and verifying that a certificate with specific content is
presented.

### Certificate Issuance
![Certificate Issuance](img/DomainVerification.png)
* Server generates a Certificate Signing Request and a public key. It asks the CA to issue a certificate with this public key.
* Server signs the public key in the CSR and the CSR with the **authorised** private key.
* CA verifies both signatures and issues the certificate.
* Server receives the certificate and installs it on the relevant domain.

Likewise, for revocation, a revocation request is generated and signed with the **authorised** private key. It is then sent to the CA to revoke the certificate.

