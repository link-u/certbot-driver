# CertBot Driver (uses Route53 DNS)
[![Build on Linux](https://github.com/link-u/certbot-driver/workflows/Build%20on%20Linux/badge.svg)](https://github.com/link-u/certbot-driver/actions?query=workflow%3A%22Build+on+Linux%22)
[![Build on macOS](https://github.com/link-u/certbot-driver/workflows/Build%20on%20macOS/badge.svg)](https://github.com/link-u/certbot-driver/actions?query=workflow%3A%22Build+on+macOS%22)
[![Build on Windows](https://github.com/link-u/certbot-driver/workflows/Build%20on%20Windows/badge.svg)](https://github.com/link-u/certbot-driver/actions?query=workflow%3A%22Build+on+Windows%22)  
[![Build debian packages](https://github.com/link-u/certbot-driver/workflows/Build%20debian%20packages/badge.svg)](https://github.com/link-u/certbot-driver/actions?query=workflow%3A%22Build+debian+packages%22)

It controls certbot to create and renew certs, using AWS Route 53 DNS Plugin.

# How to build

## Pre-requirements

 - Docker
   - This program uses docker internally.
 - Golang
 - AWS IAM for Route53
   - Please read [this page](https://certbot-dns-route53.readthedocs.io/en/stable/) and please prepare it.

## How to build

```bash
make certbot-driver

./certbot-driver --help
```

# How to use

## usage
```bash
% ./certbot-driver
usage: certbot-driver [<flags>] <command> [<args> ...]

Control certbot automatically

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  create --cert.directory=(path/to/cert) --email-address=(aoba@example.com) --aws.iam=(iam.conf) [<flags>] <domains>...
    create new certs

  renew --cert.directory=(path/to/cert) --email-address=(aoba@example.com) --aws.iam=(iam.conf) [<flags>]
    renew existing certs
```

## create

It creates a certificate for `example.com` and `*.example.com`.

```bash
% certbot-driver \
  --cert.directory=data/example.com \
  --email-address=your-name@example.com \
  --aws.iam=route53.iam.conf \
  'example.com' '*.example.com'
```

## renew

It keeps or renew certificates.

```bash
% certbot-driver \
  --cert.directory=data/example.com \
  --email-address=your-name@example.com \
  --aws.iam=route53.iam.conf
```

# How to use certificates?

Certificates are stores in the directory as in `/etc/letsencrypt`.

You can use
 - `path/to/certs/live/example.com/privkey.pem`
 - `path/to/certs/live/example.com/fullchain.pem`

In nginx, apache or other HTTP servers.

Please see example for more details:

```bash
cd path/to/cert
find .
./csr
./csr/0000_csr-certbot.pem
./keys
./keys/0000_key-certbot.pem
./renewal
./renewal/example.com.conf
./archive
./archive/example.com
./archive/example.com/chain1.pem
./archive/example.com/fullchain1.pem
./archive/example.com/privkey1.pem
./archive/example.com/cert1.pem
./live
./live/README
./live/example.com
./live/example.com/privkey.pem
./live/example.com/chain.pem
./live/example.com/cert.pem
./live/example.com/fullchain.pem
./live/example.com/README
./accounts
./accounts/acme-v02.api.letsencrypt.org
./accounts/acme-v02.api.letsencrypt.org/directory
./accounts/acme-v02.api.letsencrypt.org/directory/7b0ea06ef2adc55dd70bdf6902e9b10e
./accounts/acme-v02.api.letsencrypt.org/directory/7b0ea06ef2adc55dd70bdf6902e9b10e/private_key.json
./accounts/acme-v02.api.letsencrypt.org/directory/7b0ea06ef2adc55dd70bdf6902e9b10e/regr.json
./accounts/acme-v02.api.letsencrypt.org/directory/7b0ea06ef2adc55dd70bdf6902e9b10e/meta.json
./accounts/acme-staging-v02.api.letsencrypt.org
./accounts/acme-staging-v02.api.letsencrypt.org/directory
./accounts/acme-staging-v02.api.letsencrypt.org/directory/3a3615f3d27cc339e1d4e5ed52275f45
./accounts/acme-staging-v02.api.letsencrypt.org/directory/3a3615f3d27cc339e1d4e5ed52275f45/private_key.json
./accounts/acme-staging-v02.api.letsencrypt.org/directory/3a3615f3d27cc339e1d4e5ed52275f45/regr.json
./accounts/acme-staging-v02.api.letsencrypt.org/directory/3a3615f3d27cc339e1d4e5ed52275f45/meta.json
./renewal-hooks
./renewal-hooks/post
./renewal-hooks/pre
./renewal-hooks/deploy
```
