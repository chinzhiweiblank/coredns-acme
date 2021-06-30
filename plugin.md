# acme

## Motivation and Use Cases
### Why do you need certificates?
Certificates allow you to perform SSL/TLS where communication between the client and the server is encrypted so that only the intended recipient can access it. Information you send on the Internet is passed from computer to computer to get to the destination server. Any computer in between you and the server can see your sensitive information if it's not encrypted.

When an SSL/TLS certificate is used, the information becomes unreadable to everyone except for the server you are sending the information to. This protects it from hackers and identity thieves.

### Managing Certificates Manually
To generate a TLS/SSL certificate, you need to do the following:
1. Generate a Certificate Signing Request (CSR)
2. Cut and paste the CSR into a Certificate Authority's (CA) web page
3. Prove ownership of the domain(s) in the CSR through the CA's challenges
4. Download the issued certificate and install it on the user's server

Managing certificates manually poses a risk to systems in production because:
1. Users can forget to renew certificate until expiration
2. Risk of exposure leads to gaps in ownership and hence Man-in-the-Middle attacks and breaches.

## How does ACME benefit you?
ACME enables automatic renewal, replacement and revocation of domain validated TLS/SSL certificates.

By doing so,
1. You no longer have to spend energy and time to keep a watch on their expiry dates and worry about certificates expiring.
2. You no longer have to dig out the instructions to renew and configure the certificates.
3. You no longer have to worry about data breaches or Man-in-the-Middle attacks that happen when your certificates expire

Just set up the ACME plugin once and let it run for you. At companies, this could save a lot of manpower and time when there are hundreds of certificates in use.

## How it works
#### Domain Validation
* Server generates a key-value pair: a public and private key. Server keeps the private key
* Certificate Authority (CA) issues one or more challenges to prove that the server controls the domain.
* CA also provides an arbitrary number (a nonce) to sign with the private key. This proves that it controls the key-value pair.
* Server fulfills the challenge and signs the provided nonce.
* CA verifies the nonce and checks if the challenge is fulfilled.
* Server is authorised to do certificate management for the domain with the key-value pair. The key-value pair is now known
as the **authorised** key-value pair.

#### Certificate Issuance
* Server generates a Certificate Signing Request and a public key. It asks the CA to issue a certificate with this public key.
* Server signs the public key in the CSR and the CSR with the **authorised** private key.
* CA verifies both signatures and issues the certificate.
* Server receives the certificate and installs it on the relevant domain.

Likewise, for revocation, a revocation request is generated and signed with the **authorised** private key. It is then sent to the CA to revoke the certificate.
