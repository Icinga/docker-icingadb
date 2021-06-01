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
	-v icingadb:/data \
	-e ICINGADB_REDIS_ADDRESS=redis-icingadb:6379 \
	-e ICINGADB_REDIS_PASSWORD=123456 \
	-e ICINGADB_DATABASE_HOST=mariadb-icingadb \
	-e ICINGADB_DATABASE_PORT=3306 \
	-e ICINGADB_DATABASE_DATABASE=icingadb \
	-e ICINGADB_DATABASE_USER=icingadb \
	-e ICINGADB_DATABASE_PASSWORD=123456 \
	icinga/icingadb
```

The container expects a volume on `/data` and
takes the environment variables shown above.

Each environment variable corresponds to a configuration option of Icinga DB.
E.g. `ICINGADB_REDIS_ADDRESS=2001:db8::192.0.2.18` means:

```yaml
redis:
  address: 2001:db8::192.0.2.18
```

Consult the [Icinga DB configuration documentation] on what options there are.

Icinga DB automatically imports and upgrades the SQL database schema.
In a HA setup this may lead to a broken database.
Pass the environment variable `ICINGADB_SLEEP=30` to one of the instances
to sleep e.g. 30 seconds, so the other instance has this time to finish the upgrade.

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
