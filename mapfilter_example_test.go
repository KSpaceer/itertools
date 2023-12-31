package itertools_test

import (
	"fmt"
	"github.com/KSpaceer/itertools"
	"strings"
	"time"
)

type Account string

type Transaction struct {
	From      Account
	To        Account
	Amount    int
	Timestamp time.Time
}

type Message struct {
	FromTo string
	Amount int
	Time   string
}

func messageToTransaction(msg Message) Transaction {
	var tx Transaction
	parts := strings.SplitN(msg.FromTo, "---", 2)
	tx.From, tx.To = Account(parts[0]), Account(parts[1])
	tx.Amount = msg.Amount
	tx.Timestamp, _ = time.Parse(time.RFC3339Nano, msg.Time)
	return tx
}

type FraudChecker struct {
	definitelyFraudulent Account
}

func (fc *FraudChecker) IsFraudulentTransaction(tx Transaction) bool {
	return tx.To == fc.definitelyFraudulent || tx.From == fc.definitelyFraudulent
}

type Alerter struct{}

func (Alerter) Alert(tx Transaction) {
	var sb strings.Builder
	sb.WriteString("!!!FOUND FRAUD TRANSACTION!!!\n")
	sb.WriteString(fmt.Sprintf("From: %q\n", tx.From))
	sb.WriteString(fmt.Sprintf("To: %q\n", tx.To))
	sb.WriteString(fmt.Sprintf("Amount: %d\n", tx.Amount))
	fmt.Println(sb.String())
}

type MessageConsumer []Message

func (mc *MessageConsumer) StartConsume() <-chan Message {
	messages := make([]Message, len(*mc))
	copy(messages, *mc)

	ch := make(chan Message)

	go func() {
		for _, msg := range messages {
			ch <- msg
		}
		close(ch)
	}()
	return ch
}

func Example_mapFilter() {
	const fraud = "1"

	consumer := MessageConsumer{
		{
			FromTo: "2---3",
			Amount: 15,
			Time:   time.Now().Format(time.RFC3339Nano),
		},
		{
			FromTo: "3---1",
			Amount: 15,
			Time:   time.Now().Format(time.RFC3339Nano),
		},
		{
			FromTo: "5---6",
			Amount: 13,
			Time:   time.Now().Format(time.RFC3339Nano),
		},
		{
			FromTo: "7---8",
			Amount: 5,
			Time:   time.Now().Format(time.RFC3339Nano),
		},
		{
			FromTo: "4---1",
			Amount: 10,
			Time:   time.Now().Format(time.RFC3339Nano),
		},
		{
			FromTo: "1---0",
			Amount: 25,
			Time:   time.Now().Format(time.RFC3339Nano),
		},
	}

	ch := consumer.StartConsume()

	fraudChecker := FraudChecker{definitelyFraudulent: fraud}

	alerter := Alerter{}

	iter := itertools.Map(
		itertools.NewChanIterator(ch),
		messageToTransaction,
	).Filter(fraudChecker.IsFraudulentTransaction)

	iter.Range(func(tx Transaction) bool {
		alerter.Alert(tx)
		return true
	})
	// Output:
	// !!!FOUND FRAUD TRANSACTION!!!
	// From: "3"
	// To: "1"
	// Amount: 15
	//
	// !!!FOUND FRAUD TRANSACTION!!!
	// From: "4"
	// To: "1"
	// Amount: 10
	//
	// !!!FOUND FRAUD TRANSACTION!!!
	// From: "1"
	// To: "0"
	// Amount: 25
	//
}
