# client-server-api

## Utilização

```
cd client-server-api
go run server/main.go
go run client/main.go
```

## Conferindo os dados no banco de dados

```
sqlite3 server/quotes.db
select * from dollars;
```

## Conferindo o arquivo de cotação

```
cat client/cotacao.txt
```