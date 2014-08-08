build:
	rm -Rf camarabook-data
	go build

all: build
	./camarabook-data --save-deputies-from-search
	./camarabook-data --save-deputies-from-xml

  # save_from_deputado_about_parser
  # save_images_from_deputies_json_parser
	# cotas
	# video
	# proposition
