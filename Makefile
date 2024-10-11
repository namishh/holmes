build:
	@npx tailwind -i ./assets/app.css -o ./public/app.css
	@templ generate
	@go build -o bin/app ./app/

css:
	@npx tailwind -i ./assets/app.css -o ./public/app.css --watch

run: build-app
	@./bin/app

templ:
	@templ generate --watch

dev:
	@export ENVIRONMENT="dev" ; ~/go/bin/air


clean:
	@rm -rf bin
