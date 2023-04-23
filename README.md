# haynay-invoice

An (over-engineered 😁) Telegram bot [@HaynayInvoiceBot](http://t.me/HaynayInvoiceBot) for generate an invoice via interactive form.

## Features

- Automatic ID generation
- Automatic calculation
- Common input validations
- Cancel form
- Edit the previous step

## Technologies

- Clean architecture
- Using [Telegram Bot API](https://github.com/go-telegram-bot-api/telegram-bot-api)
- No database

## How to Develop

Clone this repo

```sh
git clone git@github.com:naufalmalikr/haynay-invoice.git
cd haynay-invoice
```

Install dependencies

```sh
go mod download
```

Copy the sample of environment variables, then edit as you need

```sh
cp .envsample .env
```

## How to Run

Run without compile

```sh
go run .
```
