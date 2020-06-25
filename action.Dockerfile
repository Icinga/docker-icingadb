FROM icinga/icingadb-builder

COPY action.bash Dockerfile /

CMD ["/action.bash"]
