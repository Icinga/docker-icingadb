<!-- Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+ -->

# Icinga DB - Docker image

This image integrates [Icinga DB] into your [Docker] environment.

## Usage

```bash
docker network create icinga

docker run --rm -d \
	--network icinga \
	--name redis-icingadb \
	redis

docker run --rm -d \
	--network icinga \
	--name mariadb-icingadb \
	-e MYSQL_RANDOM_ROOT_PASSWORD=1 \
	-e MYSQL_DATABASE=icingadb \
	-e MYSQL_USER=icingadb \
	-e MYSQL_PASSWORD=123456 \
	mariadb

docker run -d \
	--network icinga \
	--restart always \
	-e ICINGADB_REDIS_ADDRESS=redis-icingadb:6379 \
	-e ICINGADB_REDIS_PASSWORD=123456 \
	-e ICINGADB_DATABASE_HOST=mariadb-icingadb \
	-e ICINGADB_DATABASE_PORT=3306 \
	-e ICINGADB_DATABASE_DATABASE=icingadb \
	-e ICINGADB_DATABASE_USER=icingadb \
	-e ICINGADB_DATABASE_PASSWORD=123456 \
	icinga/icingadb
```

The container doesn't need any volumes and
takes the environment variables shown above.

Each environment variable corresponds to a configuration option of Icinga DB.
E.g. `ICINGADB_REDIS_ADDRESS=redis-icingadb:6379` means:

```yaml
redis:
  address: redis-icingadb:6379
```

Consult the [Icinga DB configuration documentation] on what options there are.

### Connect via TLS

```bash
docker run -d \
	--network icinga \
	--restart always \
	-e ICINGADB_REDIS_ADDRESS=redis-icingadb:6379 \
	-e ICINGADB_REDIS_PASSWORD=123456 \
	-e ICINGADB_REDIS_TLS=true \
	-e ICINGADB_REDIS_CERT='-----BEGIN CERTIFICATE-----
MIIBAzCBrgIBKjANBgkqhkiG9w0BAQQFADANMQswCQYDVQQDDAI0MjAeFw0yMTA1
MTcxMDI3MDlaFw0yMTA1MTgxMDI3MDlaMA0xCzAJBgNVBAMMAjQyMFwwDQYJKoZI
hvcNAQEBBQADSwAwSAJBANkBa53UGhd9RYiAZPGOz0/Y9P4/o6uHw/Eh4ExgCrpx
17NNV1JSAQlVnHtVANGmdz9J0c0MWC2ya3o39BbK7/cCAwEAATANBgkqhkiG9w0B
AQQFAANBACma7rGAI3khftF9du1KwivWzeGPHJwZBMfL/F99d2ckTyQozLTTL/p3
U1aTnHBR8cl5yTMAD8onBa/j7HhvL/Q=
-----END CERTIFICATE-----' \
	-e ICINGADB_REDIS_KEY='-----BEGIN PRIVATE KEY-----
MIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT8wggE7AgEAAkEA2QFrndQaF31FiIBk
8Y7PT9j0/j+jq4fD8SHgTGAKunHXs01XUlIBCVWce1UA0aZ3P0nRzQxYLbJrejf0
Fsrv9wIDAQABAkEAm0xV7MRey9Kd0Vs5Ylm2aUk1w0Jd6iKmCkoZD+9nnhcKSNuR
Jf3I9OAXYWCOIEszfrFyAQDTdp9UrOyeE9U7SQIhAP7cREMVA0NryBqYwJketN54
3unUJGBkVeumyXA/EMIFAiEA2fnScmRn4cXqqxe9Dkgn2RiogTkCZ8h5BdY67xta
nssCIF6gT+QMUDrfMNvXLWNsyED15eYxsxPrDQ/CzHYVpFY1AiEAz080gatQyX+s
kpB/NCgYDffPuyb3TLFzuMNpRaOkakUCIHtBnos4xywZBqDdRIenbxRdQHX/llUx
r1WLl8RkIQ3V
-----END PRIVATE KEY-----' \
	-e ICINGADB_REDIS_CA='-----BEGIN CERTIFICATE-----
MIIBAzCBrgIBKjANBgkqhkiG9w0BAQQFADANMQswCQYDVQQDDAI0MjAeFw0yMTA1
MTcxMDI3MDlaFw0yMTA1MTgxMDI3MDlaMA0xCzAJBgNVBAMMAjQyMFwwDQYJKoZI
hvcNAQEBBQADSwAwSAJBANkBa53UGhd9RYiAZPGOz0/Y9P4/o6uHw/Eh4ExgCrpx
17NNV1JSAQlVnHtVANGmdz9J0c0MWC2ya3o39BbK7/cCAwEAATANBgkqhkiG9w0B
AQQFAANBACma7rGAI3khftF9du1KwivWzeGPHJwZBMfL/F99d2ckTyQozLTTL/p3
U1aTnHBR8cl5yTMAD8onBa/j7HhvL/Q=
-----END CERTIFICATE-----' \
	-e ICINGADB_DATABASE_HOST=mariadb-icingadb \
	-e ICINGADB_DATABASE_PORT=3306 \
	-e ICINGADB_DATABASE_DATABASE=icingadb \
	-e ICINGADB_DATABASE_USER=icingadb \
	-e ICINGADB_DATABASE_PASSWORD=123456 \
	-e ICINGADB_DATABASE_TLS=true \
	-e ICINGADB_DATABASE_CERT='-----BEGIN CERTIFICATE-----
MIIBAzCBrgIBKjANBgkqhkiG9w0BAQQFADANMQswCQYDVQQDDAI0MjAeFw0yMTA1
MTcxMDI3MDlaFw0yMTA1MTgxMDI3MDlaMA0xCzAJBgNVBAMMAjQyMFwwDQYJKoZI
hvcNAQEBBQADSwAwSAJBANkBa53UGhd9RYiAZPGOz0/Y9P4/o6uHw/Eh4ExgCrpx
17NNV1JSAQlVnHtVANGmdz9J0c0MWC2ya3o39BbK7/cCAwEAATANBgkqhkiG9w0B
AQQFAANBACma7rGAI3khftF9du1KwivWzeGPHJwZBMfL/F99d2ckTyQozLTTL/p3
U1aTnHBR8cl5yTMAD8onBa/j7HhvL/Q=
-----END CERTIFICATE-----' \
	-e ICINGADB_DATABASE_KEY='-----BEGIN PRIVATE KEY-----
MIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT8wggE7AgEAAkEA2QFrndQaF31FiIBk
8Y7PT9j0/j+jq4fD8SHgTGAKunHXs01XUlIBCVWce1UA0aZ3P0nRzQxYLbJrejf0
Fsrv9wIDAQABAkEAm0xV7MRey9Kd0Vs5Ylm2aUk1w0Jd6iKmCkoZD+9nnhcKSNuR
Jf3I9OAXYWCOIEszfrFyAQDTdp9UrOyeE9U7SQIhAP7cREMVA0NryBqYwJketN54
3unUJGBkVeumyXA/EMIFAiEA2fnScmRn4cXqqxe9Dkgn2RiogTkCZ8h5BdY67xta
nssCIF6gT+QMUDrfMNvXLWNsyED15eYxsxPrDQ/CzHYVpFY1AiEAz080gatQyX+s
kpB/NCgYDffPuyb3TLFzuMNpRaOkakUCIHtBnos4xywZBqDdRIenbxRdQHX/llUx
r1WLl8RkIQ3V
-----END PRIVATE KEY-----' \
	-e ICINGADB_DATABASE_CA='-----BEGIN CERTIFICATE-----
MIIBAzCBrgIBKjANBgkqhkiG9w0BAQQFADANMQswCQYDVQQDDAI0MjAeFw0yMTA1
MTcxMDI3MDlaFw0yMTA1MTgxMDI3MDlaMA0xCzAJBgNVBAMMAjQyMFwwDQYJKoZI
hvcNAQEBBQADSwAwSAJBANkBa53UGhd9RYiAZPGOz0/Y9P4/o6uHw/Eh4ExgCrpx
17NNV1JSAQlVnHtVANGmdz9J0c0MWC2ya3o39BbK7/cCAwEAATANBgkqhkiG9w0B
AQQFAANBACma7rGAI3khftF9du1KwivWzeGPHJwZBMfL/F99d2ckTyQozLTTL/p3
U1aTnHBR8cl5yTMAD8onBa/j7HhvL/Q=
-----END CERTIFICATE-----' \
	icinga/icingadb
```

## Build it yourself

```bash
git clone https://github.com/Icinga/icingadb.git
#pushd icingadb
#git checkout v1.0.0
#popd

./build.bash ./icingadb
```

[Icinga DB]: https://github.com/Icinga/icingadb
[Docker]: https://www.docker.com
[Icinga DB configuration documentation]: https://icinga.com/docs/icingadb/latest/doc/03-Configuration/
