# OLX Crawler

OLX Crawler is an app written in Go to observe olx.pl ads. Do you want to buy something new for you, your family or friends? Then this app fits your needs. You don't need to check every ad on your own. Our app will do this for you!

## Motivation

It's my friend's idea. I don't care about websites like allegro or olx.pl, but I wanted to help him, so I made this tinny application.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system. **WARNING: The development environment is set up only for windows.**

### Prerequisites

You need to install this to run this app:

- Golang
- Node.JS
- Yarn
- windres

### Installing

1. Clone this repository.
2. Navigate to the right directory.
3. Type "dev" in your command prompt.
4. Open the second command prompt, navigate to the directory where is application source code, then to client directory.
5. Type "yarn install" and wait until the download finished.
6. Type "yarn run start" in your command prompt.
7. Navigate to localhost:3000. If there is any content, then everything works.

## How to build the app to a standalone .exe file

I prepared batch files to make this easy as I can.

1. Navigate to client directory.
2. Type "yarn run build" in your command prompt and wait until the build finished.
3. Navigate to the root directory.
4. Run build-windows.bat.

## Screenshots

## Tech/framework used

<b>Built with</b>

- [Echo](https://echo.labstack.com/) - The go web framework used to serve API (and frontend if we talk about production build).
- [colly](https://github.com/gocolly/colly) - Used to scrape the website.
- [robfig/cron](https://github.com/robfig/cron) - Used to perform scheduled tasks.
- [gorm](https://gorm.io/) - The ORM library for Golang used to communicate with the sqlite3 database.
- [React](https://reactjs.org) - The web framework used to build UI.
- [Material-UI](https://material-ui.com/) - The React UI framework.

## Author

This app is fully made by [Dawid Wysokiński (Kichiyaki)](https://dawid-wysokinski.pl).
