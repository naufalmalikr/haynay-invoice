package main

import (
	_ "github.com/joho/godotenv/autoload"

	"haynay/haynay-invoice/invoice/handler/bot"
	"haynay/haynay-invoice/invoice/repository"
	"haynay/haynay-invoice/invoice/service"
	"log"
	"os"
	"os/signal"
)

func main() {
	log.Println("Starting bot...")
	invoiceRepository := repository.New()
	invoiceService := service.New(invoiceRepository)
	invoiceBot := bot.New(os.Getenv("INVOICE_BOT_TOKEN"), invoiceService)
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
