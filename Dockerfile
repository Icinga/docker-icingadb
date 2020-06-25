FROM icinga/icingadb-deps
COPY icingadb /

USER icingadb
CMD ["/entrypoint"]
