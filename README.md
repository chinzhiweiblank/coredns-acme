# acme

## Name
*acme* automates certificate management: issuance and renewal.

## Description


## Syntax
The basic syntax is
~~~txt
acme {
  domain example.com
}
~~~

## Advanced Syntax
~~~txt
acme {
  domain contoso.com

  # optional parameters
  challenge port <PORT>
}
~~~

The `challenge port <PORT>` 
* The default ports of the challenges, `tlsalpn` and `http01`, are 443 and 80 respectively. They will be used if no port is provided.
* DNS01 challenge will always be used.

You can specify one or more challenges the CA can use to verify that
you own the domain.

## Examples
## See Also
