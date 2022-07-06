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
	-e ICINGADB_REDIS_HOST=redis-icingadb \
	-e ICINGADB_REDIS_PORT=6379 \
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
E.g. `ICINGADB_REDIS_HOST=redis-icingadb` means:

```yaml
redis:
  host: redis-icingadb
```

Consult the [Icinga DB configuration documentation] on what options there are.

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
