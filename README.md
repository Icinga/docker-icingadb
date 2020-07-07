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
	-e ICINGADB_LOGGING_LEVEL=debug \
	-e ICINGADB_REDIS_HOST=redis-icingadb \
	-e ICINGADB_REDIS_PORT=6379 \
	-e ICINGADB_REDIS_PASSWORD=123456 \
	-e ICINGADB_REDIS_POOL_SIZE=42 \
	-e ICINGADB_MYSQL_HOST=mariadb-icingadb \
	-e ICINGADB_MYSQL_PORT=3306 \
	-e ICINGADB_MYSQL_DATABASE=icingadb \
	-e ICINGADB_MYSQL_USER=icingadb \
	-e ICINGADB_MYSQL_PASSWORD=123456 \
	-e ICINGADB_MYSQL_MAX_OPEN_CONNS=42 \
	-e ICINGADB_METRICS_HOST=:: \
	-e ICINGADB_METRICS_PORT=8088 \
	icinga/icingadb
```

The container doesn't need any volumes and
takes the environment variables shown above.

Each environment variable corresponds to a configuration option of Icinga DB.
E.g. `ICINGADB_REDIS_HOST=2001:db8::192.0.2.18` means:

```ini
[redis]
host = "2001:db8::192.0.2.18"
```

Consult the [Icinga DB configuration documentation] on what options there are.

[Icinga DB]: https://github.com/Icinga/icingadb
[Docker]: https://www.docker.com
[Icinga DB configuration documentation]: https://icinga.com/docs/icingadb/latest/doc/03-Configuration/
