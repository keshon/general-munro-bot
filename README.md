### ⚠ Schrödinger's warning
The code in this repository is in Superposition and may work or may not work at all.

# General Munro — Assets backup and Telegram notification for CG Wire Kitsu task tracker

## About
**General Munro** is a web application that integrates with [Kitsu](https://www.cg-wire.com/en/kitsu.html) to do the following:
- Automatic notifications via Telegram messenger for new comments in Kitsu tasks
- Automatic backups to S3 storage of all attached files to comments in Kitsu task

## Quick start
Download latest release (windows only), update conf.toml and run.

## Long start
### Source compilation
Assuming Go is installed in the system compile the code (Windows example):
`go build -ldflags "-s -w" -o app.exe src/main.go`.

The binary file will be `app.exe`.

### Configuration
### Step 1. Create a Telegram bot via @BotFather
Contact Telegram @BotFather and:
 - get a token for the bot.
 - enable groups access.
 - disable group privacy.
 
### Step 2. Update conf.toml file
Rename `empty.conf.toml` to `conf.toml` and update each section accordingly. Each section contains verbose comments but if you stuck feel free to open issue ticket.

#### Step 3. Last touches
##### Make Telegram mentions via "@"
Use Kitsu `Phone` field in User profile to store Telegram usernames like `@someUsername` - this allows bot to mention team members if the bot sends messages to the group.

##### Multlang support
Use conf.toml file to select between Russian and English. You can create your own translation inside `locales` folder if nessesary.

### Deployement
You can deploy the app using Docker. To achieve this with suplied deloy scripts it is necessary **Docker**, **docker-compose** and [Traefik](https://github.com/zorg-industries-limited/ruby-rhod-fantastic-dockers) to be installed on target machine.

In order to build and deploy container in one hit go to `/deploy` folder, update app domain name for `HOST` variable and `KITSU_HOST` for your Kitsu instance in `.env` file - then execute `bash run.sh`. The script will git clone sources, compile them, make a docker image and run a container over it.