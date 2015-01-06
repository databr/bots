build: build_go

build_go: clean parliamentarian_bot metrosp_bot ibge_bot sabesp_bot

clean:
	@rm -Rf pkg/*

pkg/parliamentarian_bot: clean
	cd go_bot/parliamentarian_bot && go build -o ../../pkg/parliamentarian_bot
pkg/metrosp_bot: clean
	cd go_bot/metrosp_bot && go build -o ../../pkg/metrosp_bot
pkg/ibge_bot: clean
	cd go_bot/ibge_bot && go build -o ../../pkg/ibge_bot
pkg/sabesp_bot: clean
	cd go_bot/sabesp_bot && go build -o ../../pkg/sabesp_bot


parliamentarian_bot: pkg/parliamentarian_bot

metrosp_bot: pkg/metrosp_bot

ibge_bot: pkg/ibge_bot

sabesp_bot: pkg/sabesp_bot

deploy_go: build_go
	goupx pkg/parliamentarian_bot
	goupx pkg/metrosp_bot
	goupx pkg/ibge_bot
	goupx pkg/sabesp_bot

	rsync -Pavh pkg/parliamentarian_bot $(DATABR_BOT_MACHINE):/usr/local/bin/parliamentarian_bot
	rsync -Pavh pkg/metrosp_bot $(DATABR_BOT_MACHINE):/usr/local/bin/metrosp_bot
	rsync -Pavh pkg/ibge_bot $(DATABR_BOT_MACHINE):/usr/local/bin/ibge_bot
	rsync -Pavh pkg/sabesp_bot $(DATABR_BOT_MACHINE):/usr/local/bin/sabesp_bot

deploy:
	@make GOARCH=amd64 GOOS=linux deploy_go
