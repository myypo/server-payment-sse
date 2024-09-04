default: 
    just --list --unsorted

@local:
    #!/usr/bin/env bash

    docker compose --file ./infra/local/local.docker-compose.yaml up -d --build
    until docker exec payment-sse-database pg_isready -p 5432 ; do sleep 0.25 ; done
    sleep 2

@test:
    #!/usr/bin/env bash

    set -a
    source .env.test
    set +a

    docker compose --file ./infra/test/test.docker-compose.yaml up -d --build
    until docker exec test-payment-sse-database pg_isready -p 12586 ; do sleep 0.25 ; done
    go test -count=1 ./...
    docker compose --file ./infra/test/test.docker-compose.yaml down -v

@create-migration name:
    goose -s -dir=./migration create {{name}} sql
