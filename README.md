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

## Observações

Na minha máquina o tempo da request feita para a api "economia.awesomeapi.com.br" excede os 200ms, eu acredito que seja por conta da minha internet ou então pela demora da api, todos os teste local eu fiz aumentando esse tempo para 300, e o cliente aumentei para 400. Mas no projeito deixei como está solicitado.
