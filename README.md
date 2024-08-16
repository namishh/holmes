![header](https://i.imgur.com/3EIbYcO.png)

<div align="center">

<img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white"> <img src="https://img.shields.io/badge/sqlite-%2307405e.svg?style=for-the-badge&logo=sqlite&logoColor=white"> <img src="https://img.shields.io/badge/tailwindcss-%2338B2AC.svg?style=for-the-badge&logo=tailwind-css&logoColor=white"> <img src="https://img.shields.io/badge/DigitalOcean-%230167ff.svg?style=for-the-badge&logo=digitalOcean&logoColor=white">

</div>

### Holmes
A simplistic powerful template / starter kit for making your own cryptic / scavenger hunts written completely in Golang and [templ](https://github.com/a-h/templ)

### Features
- No javascript *
- Very easy to extend
- Built in admin panel to make


### Run Locally

- Clone the project

```bash
  git clone https://github.com/namishh/holmes
```

- Go to the project directory

```bash
  cd holmes
```

- Install go dependencies

```bash
  # air for live reload
  go install github.com/air-verse/air@latest
  # templ
  go install github.com/a-h/templ/cmd/templ@latest
```

- Install npm dependencies for tailwind

```bash
  npm i
```

- Start the server

```bash
  ## start templ
  make templ

  ## start tailwind watcher
  make css

  ## run the server
  make dev
```

- To build the project
```bash
  make build
```


### Screenshots

- Admin Panel
![admin panel](https://i.imgur.com/48Vrz7j.png)

### Todo
- [x] Setup the project
- [x] Auth
  - [x] Login and Register
  - [x] Logout
- [x] Admin
  - [x] Delete User
  - [x] Add Question
  - [x] Delete Question
  - [x] Edit Question
    - [x] Edit Details
    - [x] Edit Media
  - [x] Hints
    - [x] Add hint
    - [x] Delete hint
- [ ] Game flow
  - [ ] Get Question
  - [ ] Submit Answer
  - [ ] Get Hint
  - [ ] Get Score
- [ ] Leaderboard
- [ ] Host on digital ocean
- [ ] Extra
  - [ ] Complete the home page
  - [ ] Rate limiting

### Aknowledgements

-  [emarifer/go-echo-templ-htmx](https://github.com/emarifer/go-echo-templ-htmx)
