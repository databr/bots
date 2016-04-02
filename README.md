# Bots

Junte-se a nós no Slack [clicando aqui](http://databr.herokuapp.com/), venha conversar com os mantenedores e interessados em API e dados publicos 

## Como adicionar um bot?

* Crie um repositorio para seu bot
* Adicione nele um Dockerfile e monte seu ambiente
* Crie seu bot, ao cria-lo lembre-se:
  * Seu Dockerfile deve ter o `CMD` que rode o bot
  * Formate seu bot para BUSCAR -> EXTRAIR -> SALVAR, ele não deve ser um daemon
  * O banco de dados é acessivel pela variavel de ambiente `MONGO_URL` (ex: user1:easypass@server/db)
  * O nome da database está na variavel de ambiente `MONGO_DATABASE_NAME` (ex: database)
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

Em qualquer momento abra uma issues nesse repositorio ou me contate para que eu possa te ajudar, pelo twitter [_dukex](https://twitter.com/_dukex) ou email [duke at databr.io](mailto:duke@databr.io)
