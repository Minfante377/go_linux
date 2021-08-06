# Remote execution of applications running on a linux server using Telegram

## About the project

The aim of this project is to allow the users to execute different linux
commands remotelly using the Telegram cli.
This is a multi-thread application that handles simultaneously the requests of
multiple users. Configuration is aimed to be simple.

## Usage

- Clone this repository on the Linux server you want to execute commands on.
- Run make install to install all the dependencies.
- Configure the .env file.
- Configure the .secrets file. In this file you must specify:

  - TELEGRAM_TOKEN=YOUR BOT TOKEN
  - USER=USER SUDO PASSWORD
- Execute make run to directly run the application or make build to build the
  executable.

# Work in progress...
