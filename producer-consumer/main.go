package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var pizzaMade, pizzaFailed, total int
var NumberOfPizzas = 10

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12)
		msg := ""
		success := false

		// simulate failure "x"
		if rnd < 5 {
			pizzaFailed++
		} else {
			pizzaMade++
		}

		total++

		fmt.Printf("Making pizza #%d. Please wait for %d seconds...\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		// simulate failure "y"
		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready", pizzaNumber)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}

		return &p
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we are making
	var i = 0

	// run forever or until we receive a quit notification
	for {
		// try to make pizza
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			// similar to switch statements
			// they are helpful for channels
			select {
			// tried to make a a pizza(we sent something to the data channel)
			case pizzaMaker.data <- *currentPizza:
			// failed to make pizza
			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func main() {
	rand.Seed(time.Now().UnixNano())

	color.Cyan("The Pizzeria is open for business!")
	color.Cyan("----------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background
	go pizzeria(pizzaJob)

	// create and run consumer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery!", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad!")
			}
		} else {
			color.Cyan("Done making pizzas...")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error closing channel!", err)
			}
		}
	}

	// print out the ending message
	color.Cyan("-----------------")
	color.Cyan("Done for the day.")
	color.Cyan("We made %d pizzas, but failed to make %d, with %d attempts in total", pizzaMade, pizzaFailed, total)

	switch {
	case pizzaFailed > 9:
		color.Red("It was an awful day!")
	case pizzaFailed >= 6:
		color.Red("It was not a very good day...")
	case pizzaFailed >= 4:
		color.Red("It was an ok day...")
	case pizzaFailed >= 2:
		color.Red("It was a pretty good day!")
	default:
		color.Green("It was a great day!!!")
	}

}
