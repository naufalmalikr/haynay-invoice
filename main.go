package main

import (
	"haynay/haynay-invoice/invoice/handler/bot"
	"haynay/haynay-invoice/invoice/repository"
	"haynay/haynay-invoice/invoice/service"
	"log"
	"os"
	"os/signal"
)

const (
	invoiceBotToken = "5607053437:AAGESIicXRslzHvkWFLqSr1RNAVJhypr4mo"
)

func main() {
	log.Println("Starting bot...")
	invoiceRepository := repository.New()
	invoiceService := service.New(invoiceRepository)
	invoiceBot := bot.New(invoiceBotToken, invoiceService)
	go invoiceBot.Start()
	log.Println("Bot started")
	waitKilled()
	log.Println("Stopping bot...")
	invoiceBot.Stop()
	log.Println("Bot stopped")
}

func waitKilled() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
}
