# General Munro — Telegram notification Bot for Kitsu
![Munro Bot in action](https://github.com/zorg-industries-limited/general-munro-bot/blob/main/misc/munro_in_action.gif)
## About
**General Munro** is a Telegram bot that publish notifications from [Kitsu](https://www.cg-wire.com/en/kitsu.html) via [Custom Actions](https://kitsu.cg-wire.com/custom-actions/).

### ⚠ Schrödinger's warning
This code is in Superposition and may work or may not work at all.

## Source compilation
The bot written in Go so it is necessary to compile the code. Assuming Go is installed in the system use supplied scripts `build.bat` or `bash build.sh` depending your OS architecture.

The binary file will be `bot.exe` or `bot` respectively.

## Docker deploy
You can deploy the bot using Docker. It is necessary **Docker** and **docker-compose** to be installed on target machine with running [Traefik](https://github.com/zorg-industries-limited/ruby-rhod-fantastic-dockers) reverse proxy.

In order to build and deploy container go to `/deploy` folder, update domain name for `HOST` variable in `.env` file and execute `bash run.sh` script that will git clone sources, compile them, make a docker image and run a container over it.


## Configuration
There are certain actions had to be made prior bot compilation.

### @BotFather
Contact @BotFather in Telegram, register new bot and:
 - get a token.
 - enable groups access.
 - disable group privacy.
 
### Kitsu
Inside Kitsu you need to create a Custom Action with the following steps:
- add bot hostname. If you use Docker deploy it is the one that is stored inside `HOST` variable in `.env` file. If you test the bot on local machine - run the bot and see for the CLI output. The hostname should printed by **Fiber** framework.
- check all entities type.
- select ajax mode.

### Config file conf.toml
Bot relies on a config file that should be properly set up. Rename `empty.conf.toml` to `conf.toml` and edit it:
- store bot token to `token` variable inside `[bot]` section.
- update `[kitsu]` section accordingly. Note that login credentials must be set for Studio Manager.
- update `[credentials]` sections.

### Update Credentials section
#### Get chat ID for each Telegram contact bot should write to. 
It can be private chat or a group.

Before going any further you need to become an administrator. To do so let's get our own `chat id` by writing to bot privately with a keyword `lookup` - the bot will print Telegram response object where we need to find `id` key inside `chat`: 
```
"message": {
    ...
    "chat": {
        "id": 285016301,
        ...
    }
    ...
}
```
Save that ID to `admin_chat_id` variable inside `[credentials]` section and restart the bot. You are now admin.

From now on you can repeat the process of getting chat ids for groups by adding bot to the group and use the same `lookup` command with mentioning bot name this time for example:
`@GeneralMunro lookup` 

The bot will send you another Telegram response object to your private chat (because you are admin already). Again you need to find chat id and save it to `chat_id_by_roles` array to the **right** side of `=` (equal) sign. The left side is for names of Tasks statuses in Kitsu.

#### Get short names of [task statuses](https://kitsu.cg-wire.com/customization/#modify-an-existing-task-status) in Kitsu.
Just open Kitsu settings as a Studio Manager and look for `short names` of each Task status.
Write them to the **left** sie of `=`(equal) sign.

At the end final settings for `char_id_by_roles` might look like this:
```
chat_id_by_roles = [
    "RETAKE=12345",
    "WFA=67890",
    "DONE"=24680
]
```

.. where each numerical row represents a group or a private chat in Telegram and a Task status it belongs.

### Last touch
Use Kitsu `Phone` field in User profile to store Telegram usernames like `@someUsername` - this would allow bot to mention team members if the bot sends messages to the group.