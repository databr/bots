# Bots

Join us on IRC at #databr to chat with other databr maintainers! ([web access](http://webchat.freenode.net/?channels=databr))

## Como adicionar um bot?

* Crie um repositorio para seu bot
* Adicione nele um Dockerfile e monte seu ambiente
* Crie seu bot, ao cria-lo lembre-se:
  * Seu Dockerfile deve ter o `CMD` que rode o bot
  * Formate seu bot para BUSCAR -> EXTRAIR -> SALVAR, ele não deve ser um daemon
  * O banco de dados é acessivel pela variavel de ambiente `DATABASE_URL`
* Fork esse projeto
* Crie um branch com o nome do seu bot (ex: `git checkout -b metrosp-bot`)
* Adicione seu projeto como subtree, como no exemplo abaixo:
```
$ git remote add metrosp-bot git@github.com:databr/metrosp-bot.git
$ git fetch metrosp-bot
$ git subtree add --prefix=metrosp-bot metrosp-bot/master
$ git push
```
* Crie um pull request

Em qualquer momento abra uma issues nesse repositorio ou me contate para que eu possa te ajudar, pelo twitter [_dukex](https://twitter.com/_dukex) ou email [duke at databr.io](duke@databr.io)