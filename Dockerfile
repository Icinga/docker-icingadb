FROM icinga/icingadb-deps
COPY icingadb mysql.schema.sql.bz2 /

USER icingadb
CMD ["/entrypoint"]
