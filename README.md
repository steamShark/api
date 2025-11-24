<h1 align='center'>
  🦈steamShark API
</h1>

<p align='center'>
  <a href="https://github.com/sponsors/alexandresanlim"><img alt="version" src="https://img.shields.io/badge/Version-1.0.0-blue" /></a>
  &nbsp;
  <a href="https://github.com/sponsors/alexandresanlim"><img alt="Sponsor" src="https://img.shields.io/badge/Opensource-green" /></a>
</p>
<br />

A Go (Gin + GORM) API for managing **websites**, and **occurrences** used by the SteamShark extension to track phishing and trusted Steam-related websites.


<!-- HEADER SECTION -->
<nav>
    <a href="#description">Description</a> |
    <a href="#features">Features</a> |
    <a href="#roadMap">Road Map</a> |
    <a href="#contributing">Contributing</a>
</nav>

## Features

<div id="features"></div>

- **Websites**: create, list, update, delete
  - Fields: domain, display name, TLD, description, is_scam, is_official, is_trusted, steam_login_present, risk score, risk level, status, notes
  - Relations: tags (many-to-many), occurrences (one-to-many)
- **Occurrences**: user-reported phishing attempts (linked to websites)
- JSON REST API with CORS enabled
- SQLite by default
- Ready for containerization / deployment

## 🛣️ Road Map

<div id="roadMap"></div>

### Current Version 1.0

- [x] Finish base api
    - [x] CRUD for websites
    - [] Add occurences to websites
- [] Make base documentation for API


You can see the <a href="./CHANGELOG.md">changelog</a> on github.

## 🚀 Development

```bash
go run .
```

App runs at [http://localhost:8800](http://localhost:8800)

## 🐳 Docker

### Dockerfile

To build and run with Dockerfile:

```bash
docker build -t steamshark-api .
docker run -p steamshark-api
```

Visit [http://localhost:8800](http://localhost:8800)

### Docker compose

You can simply run
```bash
docker compose up --build -d
```
Visit [http://localhost:8800](http://localhost:8800)

## 🤝Contributing

<div id="contributing"></div>

Everyone is more than welcome to contribute to the project, but for an organized participation, it's important to read the [contributing document](./CONTRIBUTING.md) before doing it!

<style>
nav{
    display: flex;
    flex-direction: row;
    gap: 20px;
}
</style>