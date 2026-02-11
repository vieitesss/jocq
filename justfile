alias f := file

_default:
    @just --list

run *parameters:
    go run ./cmd/jocq {{ parameters }}

file:
    go run ./cmd/jocq -f assets/example.json
