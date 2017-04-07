# Задание: mapi & agent.

Необходимо написать код двух программ:
* mapi - запрашивает конфигурацию клиентов CDN по ссылке **http://gcore-api-branch-1.apidev.spcdn.ru/agent_config** и 
  через RabbitMQ отправляет его нескольким серверам, на которых agent-ы ожидают эти данные; затем выходит

* agent - ожидают данные от mapi (бесконечно) и после получения генерируют конфиг nginx на их основе, который распечатывают на стандартный вывод (stdout)

Условия:
 1. и mapi, и agent подключаются как RabbitMQ по адресу **amqp://localhost:5672/**
 1. agent генерируют конфиг nginx-а с помощью двух шаблонов, **vhostTemplate** и **nginxTemplate** и библиотеки https://golang.org/pkg/text/template/:

```
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
```

3. образец того, что печатает agent:

```
error_log  /tmp/error.log;
pid        /tmp/nginx.pid;

events {
    worker_connections  65536;
}
...
upstream node2.nodeapp.top {
  server node.nodeapp.top max_fails=0;
}

server {
    listen 80;
    server_name node2.nodeapp.top ;

    access_log /var/log/nginx/node2.nodeapp.top_access.log;
    error_log  /var/log/nginx/node2.nodeapp.top_error.log;

    location / {
        proxy_pass http://node2.nodeapp.top;
    }
}
...
```
4. Код следует собирать из запускать следующим образом:
```
cd test-task
go get github.com/streadway/amqp
go build agent mapi
bin/agent &
bin/mapi
```

Дополнительный вопрос: что не так с конфигом, который генерирует agent по заданным шаблонам?
