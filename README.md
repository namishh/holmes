![header](https://i.imgur.com/wht7eCr.png)

<div align="center">

<img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white"> <img src="https://img.shields.io/badge/sqlite-%2307405e.svg?style=for-the-badge&logo=sqlite&logoColor=white"> <img src="https://img.shields.io/badge/tailwindcss-%2338B2AC.svg?style=for-the-badge&logo=tailwind-css&logoColor=white"> <img src="https://raw.githubusercontent.com/zeropsio/recipe-shared-assets/080df0587759a1692ade5e693509ce49fd5c5870/zerops-shield.svg">

</div>

### Holmes
A simplistic powerful template / starter kit for making your own cryptic / scavenger hunts written completely in Golang and [templ](https://github.com/a-h/templ)

### Features
- No javascript *
- Very easy to extend
- Built in admin panel to manage users, questions and hints
- Adding images, videos and audios to questions
- Built in authentication system
- Built in rate limiting
- Built in leaderboard
- Easy to add custom routes

\* (Except for the tailwind stuff and toggling the navbar)


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

- Fill in the `.env` file

```env 
DB_NAME="database.db"
SECRET=""
ADMIN_PASS="holmes"
ENVIRONMENT="DEV"
BUCKET_NAME="xxx"
BUCKET_ENDPOINT="xxx.xx.xxx"
BUCKET_ACCESSKEY="xxx-xxxxxxxxxxxxxxx"
BUCKET_SECRETKEY="xxxxxxxxxxxxx"
```

- Start the server

```bash
  ## start templ
  make templ

  ## start tailwind watcher
  make css

  ## run the server
  export ENVIRONMENT="DEV" ; make dev
```

- To build the project
```bash
  make build
```


### Screenshots

- Home Page 

<img src="https://i.imgur.com/sXEpJrh.png">

- Auth Pages

<img src="https://i.imgur.com/5L8OnaH.png">
<img src="https://i.imgur.com/Qsdiifl.png">
<img src="https://i.imgur.com/hW2oJq8.png">

- Admin Panel

<img src="https://i.imgur.com/m87qXdv.png">

- All the questions

<img src="https://i.imgur.com/82IJ6qC.png">

- Questions

<img src="https://i.imgur.com/0Q7v31r.png">
<img src="https://i.imgur.com/GJRdpmv.png">

- Hints

<img src="https://i.imgur.com/c89NvH9.png">
<img src="https://i.imgur.com/YvgrxDi.png">

- Single Question

<img src="https://i.imgur.com/ijS1lgD.png">
<img src="https://i.imgur.com/w9axYGy.png">

- Opened Hint 

<img src="https://i.imgur.com/exJicUs.png">

- Leaderboard

<img src="https://i.imgur.com/EhlDYx7.png">

- Error Pages

<img src="https://i.imgur.com/hFF6u05.png">
<img src="https://i.imgur.com/snKUvK9.png">

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
- [x] Game flow
  - [x] Get all questions
  - [x] Get Question
  - [x] Submit Answer
  - [x] Update Score
  - [x] Get Hint
  - [x] End Game
- [x] Leaderboard
- [x] Error Pages
- [x] Rate limiting
- [x] Hosted on Zerops

### Acknowledgements

-  [emarifer/go-echo-templ-htmx](https://github.com/emarifer/go-echo-templ-htmx)
-  [zerops](https://zerops.io) for the cheap af hosting and object storage.
