package main

import (
	"github.com/streadway/amqp"
)

var vhostTemplate string = `{{- $UpName := .Name -}}

upstream {{ $UpName }} {
  server {{ .Origin }} max_fails=0;
}

server {
    listen {{ or .Port "80" }};
    server_name {{ .Name }} ;

    access_log /var/log/nginx/{{ .Name }}_access.log;
    error_log  /var/log/nginx/{{ .Name }}_error.log;

    location / {
        proxy_pass http://{{ $UpName }};
    }
}
`

var nginxTemplate string = `error_log  /tmp/error.log;
pid        /tmp/nginx.pid;

events {
    worker_connections  65536;
}

http {
    access_log /tmp/access.log;
    client_header_timeout 10s;
    client_body_timeout   10s;
    default_type  application/octet-stream;
    sendfile    on;

{{ .Configs }}

}
`

func main() {
	_, _ = amqp.Dial("amqp://localhost:5672/")
}

