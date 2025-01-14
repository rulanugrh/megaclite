# Megaclite
![banner](.github/index.png)

<img src="https://user-images.githubusercontent.com/73097560/115834477-dbab4500-a447-11eb-908a-139a6edaec5c.gif">


# Background Project
This project was actually due to my campus assignment on database security. Finally from there I thought of the idea to build a webmail system by applying the PGP algorithm so that the login process and sending mail the data is securely encrypted which is only known by the user. 

# Requirements
- [Go](https://golang.org/) 1.18 or higher
- [Docker](https://www.docker.com/) 20.10.17 or higher
- [Docker Compose](https://docs.docker.com/compose/) 1.29.
- [Templ](https://templ.guide/) latest version
- [Taiwind CSS](https://tailwindcss.com/) latest version
- [MySQL](https://www.mysql.com/) latest version

# Installation
first, clone this repository
```bash
$ git clone https://github.com/rulanugrh/megaclite.git
```

after that you can copy depedency `.env.example` to `.env`
```bash
$ cp .env.example .env

# install package
$ go mod tidy
```
and then, you can migrate with this command
```bash
$ go run main.go migrate

# and for seed you can run
$ go run main.go seed
```

if you want running view template you can running this command
```bash
# you must install templ befor running this command
$ templ generate
$ go run main.go serve
```

or you want running API you run this command
```bash
$ go run main.go api
```

> Note: Documentation API in http://localhost:port/docs
