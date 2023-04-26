package main

import (
	_ "github.com/joho/godotenv/autoload"

	"fmt"
	"haynay/haynay-invoice/invoice/handler/bot"
	"haynay/haynay-invoice/invoice/repository"
	"haynay/haynay-invoice/invoice/service"
	"log"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Bot error: %s\n", err.Error())
	}
}

func run() error {
	log.Println("Starting bot...")
	invoiceIDSeq, err := strconv.ParseInt(os.Getenv("INVOICE_ID_SEQ"), 10, 64)
	if err != nil {
		log.Printf("Error parsing INVOICE_ID_SEQ: %s\n", err.Error())
	}

	invoiceRepository := repository.New(invoiceIDSeq)
	invoiceService := service.New(invoiceRepository)
	invoiceBot, err := bot.New(os.Getenv("INVOICE_BOT_TOKEN"), invoiceService)
	if err != nil {
		return err
	}

	panicChan := make(chan interface{}, 1)
	go invoiceBot.Start(panicChan)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	select {
	case <-sigchan:
	case c := <-panicChan:
		return fmt.Errorf("%s", c)
	}

	log.Println("Stopping bot...")
	invoiceBot.Stop()
	invoiceRepository.Shutdown()
	log.Println("Bot stopped")
	return nil
}
