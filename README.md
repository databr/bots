# Bots

Join us on IRC at #databr to chat with other databr maintainers! ([web access](http://webchat.freenode.net/?channels=databr))

## Como adicionar um bot?

* Crie um repositorio para seu bot
* Adicione nele um Dockerfile e monte seu ambient
* Crie seu bot, ao criar lembre-se:
  * Seu Dockerfile deve ter um `CMD`que rode o bot
  * Formate seu bot para BUSCA -> PARSER -> SALVA, ele não deve esperar ficar como daemon
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
