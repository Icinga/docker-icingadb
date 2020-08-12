# Icinga DB Docker image | (c) 2020 Icinga GmbH | GPLv2+

FROM icinga/icingadb-builder

COPY action.bash Dockerfile /

CMD ["/action.bash"]
