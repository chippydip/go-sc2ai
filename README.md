# go-sc2ai

Go implementation of the Starcraft II AI API

## Quick start

### Prerequisites

- Install golang >= 1.13  
    https://golang.org/dl/
- Install Starcraft II game  
    https://starcraft2.com/en-us/ (free version is okay)
- Download Ladder  
    - Join Discord messenger: https://discordapp.com/invite/Emm5Ztz
    - Use `!maps` command to get archive link
    - Download, unpack and copy files into Starcraft II `Maps` folder (e.g. `C:/Program Files (x86)/StarCraft II/Maps`)

### Build and run examples

- Create working dir:
    ```bash
    mkdir /path/godev
    cd /path/godev/
    ```
- Clone project:
    ```bash
    git clone git@github.com:chippydip/go-sc2ai.git
    cd go-sc2ai/
    ```
- Build and run example:
    ```bash
    go run .\examples\zerg_rush
    ```

## Start writing your own bot

- Make project dir:
    ```bash
    mkdir /path/mycoolbot
    cd /path/mycoolbot/
    go mod init bot
    go get -u github.com/chippydip/go-sc2ai
    ```
- Copy `main.go` from `examples/stub_bot` into your project dir
- Build and run
    ```bash
    go run .
    ```
- Add some cool new stuff
- Build > Run > Test > Repeat
- ???
- PROFIT

## Useful links

- Get help in Discord: https://discord.gg/Emm5Ztz
- Learn golang: https://tour.golang.org/
- Starcraft II AI Wiki: http://wiki.sc2ai.net/Main_Page
- More maps: https://github.com/ttinies/sc2gameMapRepo
