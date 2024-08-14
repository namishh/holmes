build-app:
	@go build -o bin/app ./app/

create:
	@go run ./app/create.go

css:
	@tailwind -i ./assets/app.css -o ./public/app.css --watch

gen-css:
	@tailwind -i ./assets/app.css -o ./public/app.css

run: build-app
	@./bin/app

templ:
	@templ generate --watch

gen-templ:
	@templ generate

dev:
	@./bin/air


clean:
	@rm -rf bin
