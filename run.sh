#/bin/sh
docker run -it --rm \
  -p 8080:8080 \
  -e 'XMLTV_PORT=8080' \
  -e 'XMLTV_DOMAIN=http://host.docker.internal:8080/' \
  xmltv:development