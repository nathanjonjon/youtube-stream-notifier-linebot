# youtube-stream-notifier-linebot


## Table of contents
* [General info](#general-info)
* [Project Structure](#project-structure)
* [Usage](#usage)


## General Info
This is a LineBot deployed on Heroku that notifies users and groups when a given YouTube channel launches video or stream. The project was meant to help my fellows in church to subscribe to our preach live streams via `Line`, especially for the elderly. Note that it is a simplfied one compared to the running service on production, and is shared to help others who want to build something alike.


## Project Structure
- `api` : Callback handlers and apis that interacts with `pubsubhubub` and database
- `controller` : Controllers of database and linebot
- `main.go` : Initiate database, linebot and automate subscription

## Usage
- Create an account in [Line developer console](https://developers.line.biz/console/)
- Configure webhook, message channel and register credentials
- Deploy codes to your app in [Heroku](https://dashboard.heroku.com/apps) and set the `Config Vars`
- Add on a `Postgresql` database to your app from [here](https://elements.heroku.com/addons/heroku-postgresql)
- Upload a video or start a stream, and stay tuned to the linebot!

