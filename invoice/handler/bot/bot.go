package bot

import (
	"fmt"
	"haynay/haynay-invoice/invoice/service"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type bot struct {
	bot     *tgbotapi.BotAPI
	config  tgbotapi.UpdateConfig
	service service.Service
	users   map[int64]*user
	steps   map[string]*step
}

func New(telegramToken string, service service.Service) (Bot, error) {
	telegramBot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return nil, err
	}

	config := tgbotapi.NewUpdate(0)
	config.Timeout = 30
	bot := bot{
		bot:     telegramBot,
		config:  config,
		service: service,
		users:   map[int64]*user{},
		steps:   map[string]*step{}}

	bot.initializeSteps()
	return &bot, nil
}

func (b *bot) initializeSteps() {
	b.steps["start"] = newStep("start").
		addChoiceTo("buyer", initializeInvoice(b), noRollback, "invoice", "1").
		allowFreeTextTo("start", noProcess, noRollback).
		simpleOutput("Assalaamu'alaikum, silahkan pilih:\n1. Ketik \"Invoice\" untuk mulai bikin invoice")

	b.steps["buyer"] = newStep("buyer").
		allowFreeTextTo("phone", setBuyer, noRollback).
		setOutput(buyerOutput).
		cancelable()

	b.steps["phone"] = newStep("phone").
		allowFreeTextTo("address", setPhone, noRollback).
		simpleOutput("Masukkan no HP pembeli").
		cancelable().
		editableFor(b.steps["buyer"])

	b.steps["address"] = newStep("address").
		allowFreeTextTo("order_item", setAddressAndInitializeOrder, removeLastOrder).
		simpleOutput("Masukkan alamat pembeli beserta kode pos").
		cancelable().
		editableFor(b.steps["phone"])

	b.steps["order_item"] = newStep("order_item").
		allowFreeTextTo("order_quantity", setOrderItem, noRollback).
		setOutput(orderItemOutput).
		cancelable().
		setEditedStep(editFromOrderItem(b))

	b.steps["order_quantity"] = newStep("order_quantity").
		allowFreeTextTo("order_price", setOrderQuantity, noRollback).
		setOutput(orderQuantityOutput).
		cancelable().
		editableFor(b.steps["order_item"])

	b.steps["order_price"] = newStep("order_price").
		allowFreeTextTo("order_more", setOrderPrice, noRollback).
		setOutput(orderPriceOutput).
		cancelable().
		editableFor(b.steps["order_quantity"])

	b.steps["order_more"] = newStep("order_more").
		addChoiceTo("courier_name", noProcess, noRollback, "ok", "oke", "1").
		addChoiceTo("order_item", initializeOrder, removeLastOrder, "tambah", "add", "2").
		rejectFreeText().
		setOutput(orderMoreOutput).
		cancelable().
		editableFor(b.steps["order_price"])

	b.steps["courier_name"] = newStep("courier_name").
		allowFreeTextTo("courier_info", setCourier, noRollback).
		simpleOutput("Masukkan nama kurir").
		cancelable().
		editableFor(b.steps["order_more"])

	b.steps["courier_info"] = newStep("courier_info").
		addChoiceTo("courier_price", noProcess, noRollback, "ok", "oke", "1").
		addChoiceTo("courier_price", setCOD, rollbackCourierInfo, "cod", "2").
		addChoiceTo("courier_price", setFreeShippingFee, rollbackCourierInfo, "gratis", "free", "3").
		rejectFreeText().
		simpleOutput("Pilih salah satu:\n1. Ketik \"OK\" untuk lanjut\n2. Ketik \"COD\" jika pembeli melakukan COD\n3. Ketik \"Gratis\" jika gratis ongkir").
		cancelable().
		editableFor(b.steps["courier_name"])

	b.steps["courier_price"] = newStep("courier_price").
		allowFreeTextTo("stock", setShippingFee, noRollback).
		setOutput(courierPriceOutput).
		cancelable().
		editableFor(b.steps["courier_info"])

	b.steps["stock"] = newStep("stock").
		addChoiceTo("confirmation", noProcess, rollbackPO, "ready", "1").
		addChoiceTo("preorder_duration", setPO, rollbackPO, "po", "2").
		rejectFreeText().
		simpleOutput("Pilih salah satu:\n1. Ketik \"Ready\" jika produk ready stock\n2. Ketik \"PO\" jika produk preorder").
		cancelable().
		editableFor(b.steps["courier_price"])

	b.steps["preorder_duration"] = newStep("preorder_duration").
		allowFreeTextTo("confirmation", setPODuration, noRollback).
		simpleOutput("Masukkan maksimal jumlah pekan untuk preorder").
		cancelable().
		editableFor(b.steps["stock"])

	b.steps["confirmation"] = newStep("confirmation").
		addChoiceTo("save", noProcess, noRollback, "ok", "oke").
		addChoiceTo("print", noProcess, noRollback, "print", "view").
		rejectFreeText().
		setOutput(confirmationOutput(b)).
		cancelable().
		editableFor(b.steps["stock"])

	b.steps["save"] = newStep("save").
		immediateTo("print", saveInvoice(b), deleteInvoice(b)).
		simpleOutput("Invoice berhasil disimpan")

	b.steps["print"] = newStep("print").
		immediateTo("start", noProcess, noRollback).
		setOutput(printInvoice(b))

	b.steps["cancel"] = newStep("cancel").
		immediateTo("start", noProcess, noRollback).
		simpleOutput("Invoice dibatalkan")

	b.steps["edit"] = newStep("edit").
		asDynamic(editAction).
		simpleOutput("Silahkan ubah data")
}

func (b *bot) Start(panicChan chan interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panicChan <- r
		}
	}()

	updates := b.bot.GetUpdatesChan(b.config)
	log.Println("Ready to receive messages")
	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		input := update.Message.Text
		log.Printf("Message from %s (%d): %s\n", update.Message.From.UserName, userID, input)
		for _, output := range b.process(userID, input) {
			b.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, output))
		}
	}
}

func (b *bot) process(userID int64, input string) []string {
	user := b.getOrInitializeUser(userID)
	step := b.steps[user.CurrentStep]

	// NOTE: don't know why got warning in line below, it says
	//       "this value of err is never used (SA4006) go-staticcheck",
	//       though it used for handleError
	action, err := b.getAction(step, input)
	step, action, outputs, err := b.processAction(step, action, user, input)
	return append(outputs, b.handleError(step, action, user, err)...)
}

func (b *bot) getOrInitializeUser(userID int64) *user {
	if user, found := b.users[userID]; found {
		return user
	}

	b.users[userID] = &user{CurrentStep: "start"}
	return b.users[userID]
}

func (b *bot) getAction(step *step, input string) (*action, error) {
	if step.ImmediateAction != nil {
		return step.ImmediateAction, nil
	} else if action := step.ChoicesAction[strings.ToLower(input)]; action != nil {
		return action, nil
	} else if step.FreeTextAction != nil {
		return step.FreeTextAction, nil
	}

	return nil, errUnexpected
}

func (b *bot) processAction(step *step, action *action, user *user, input string) (*step, *action, []string, error) {
	outputs := []string{}
	var err error
	for action != nil {
		nextStepName := action.NextStep
		nextStep := b.steps[nextStepName]
		if nextStep != nil && nextStep.DynamicAction != nil {
			outputs = append(outputs, nextStep.Output(user)...)
			action = nextStep.DynamicAction(step, user)
			continue
		}

		if err = action.Process(input, user); err != nil {
			break
		}

		outputs = append(outputs, nextStep.Output(user)...)
		user.CurrentStep = nextStepName
		step = nextStep
		action = nextStep.ImmediateAction
	}

	return step, action, outputs, err
}

func (b *bot) handleError(step *step, action *action, user *user, err error) []string {
	outputs := []string{}
	if err == nil {
		return outputs
	}

	log.Printf("Process error: %s\n", err.Error())
	outputs = append(outputs, fmt.Sprintf("Error: %s", err.Error()))
	if action == nil {
		return outputs
	}

	if err = action.Rollback(user); err == nil {
		outputs = append(outputs, "Silahkan ulangi")
		return append(outputs, step.Output(user)...)
	}

	log.Printf("Rollback error: %s\n", err.Error())
	return append(outputs, fmt.Sprintf("Error: %s", errUnexpected.Error()))
}

func (b *bot) Stop() {
	b.bot.StopReceivingUpdates()
}
