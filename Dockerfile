# Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+

FROM icinga/icingadb-deps
COPY icingadb *.sql.bz2 /

USER icingadb
CMD ["/entrypoint"]
