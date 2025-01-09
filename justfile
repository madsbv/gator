alias gc := goose-create
goose-create *args:
    just run-goose create {{args}}

alias g := run-goose
run-goose *args:
    goose -env ".env" {{args}}
