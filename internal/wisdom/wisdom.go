package wisdom

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

// Interface for getting wisdom quotes
type Wisdomer interface {
	Init(string) error
	GetQuote() []byte
}

type Wisdom struct {
	quotes []string
}

func NewWisdom() *Wisdom {
	w := &Wisdom{
		quotes: make([]string, 0),
	}
	return w
}

func (w *Wisdom) Init(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w.quotes = append(w.quotes, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (w *Wisdom) GetQuote() []byte {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(w.quotes))
	return []byte(w.quotes[index])
}
