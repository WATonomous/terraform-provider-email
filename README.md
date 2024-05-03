# Terraform Provider Email

Terraform provider for sending emails.

## Development

```bash
go mod tidy # install dependencies
make build # build the provider
make install # install the provider
```

Please see `./examples` for example usage.

## Testing

First start a local SMTP server:

```bash
docker run --rm -it -p 3000:80 -p 2525:25 rnwood/smtp4dev:v3
```

Then run the acceptance tests:

```bash
make testacc
```

## Generate docs

```bash
go generate ./...
```

## References

https://developer.hashicorp.com/terraform/tutorials/providers/provider-setup