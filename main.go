package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const factor = 1

type Product struct {
	ID   string
	Name string
}

func main() {
	channels := runtime.NumCPU() * factor
	products := make([]*Product, 0)
	ids := make([]string, 0)
	for i := 0; i < 10; i++ {
		product := Product{
			ID:   strconv.Itoa(i + 1),
			Name: fmt.Sprintf("Product: %s", strconv.Itoa(i+1)),
		}
		products = append(products, &product)
		ids = append(ids, product.ID)
	}

	productsFoundP := FindProductsByIds(products, ids, channels)
	productsFoundS := FindProductsByIdsSeq(products, ids)
	fmt.Printf("Parallel: %#v\n", productsFoundP)
	fmt.Printf("Sequential: %#v\n", productsFoundS)
}

func FindProductsByIdsSeq(products []*Product, ids []string) []*Product {
	now := time.Now()
	productsFound := make([]*Product, 0)
	for _, id := range ids {
		productFound, err := FindProductById(id, products)
		if err == nil {
			productsFound = append(productsFound, productFound)
		}

	}
	fmt.Printf("sequential version took: %v\n", time.Since(now))
	return productsFound
}

func FindProductsByIds(
	products []*Product,
	ids []string,
	channels int,
) []*Product {
	now := time.Now()
	productsFound := make([]*Product, 0)
	done := make(chan interface{})
	defer close(done)
	idsStream := getStringsStream(ids)
	finders := fanOutProducts(done, products, idsStream, channels)
	for value := range fanInProducts(done, finders...) {
		productsFound = append(productsFound, value)
	}
	fmt.Printf("parallel version took: %v\n", time.Since(now))
	return productsFound
}

func FindProductById(id string, products []*Product) (*Product, error) {
	for _, product := range products {
		sleepRandomTime()
		if product.ID == id {
			return product, nil
		}
	}
	return nil, errors.New("Product not found.")
}

func getStringsStream(ids []string) <-chan string {
	idsStream := make(chan string, len(ids))
	go func(idsStream chan string) {
		defer close(idsStream)
		for i := 0; i < len(ids); i++ {
			idsStream <- ids[i]
		}
	}(idsStream)
	return idsStream
}

func findProductsByIds(
	products []*Product,
	done <-chan interface{},
	idsStream <-chan string,
) <-chan *Product {
	productsStream := make(chan *Product)
	go func() {
		defer close(productsStream)
		for id := range idsStream {
			product, err := FindProductById(id, products)
			if err != nil {
				log.Fatal(err)
			}
			select {
			case <-done:
				return
			case productsStream <- product:
			}
		}
	}()
	return productsStream
}

func fanOutProducts(
	done <-chan interface{},
	products []*Product,
	idsStream <-chan string,
	channels int,
) []<-chan *Product {
	finders := make([]<-chan *Product, channels)
	for i := 0; i < channels; i++ {
		finders[i] = findProductsByIds(products, done, idsStream)
	}
	return finders
}

func fanInProducts(
	done <-chan interface{},
	channels ...<-chan *Product,
) <-chan *Product {
	var wg sync.WaitGroup
	multiplexedStream := make(chan *Product)
	multiplex := func(c <-chan *Product) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()
	return multiplexedStream
}

func computeBinaryInt() int {
	t := time.Now().UnixNano()
	rand.Seed(t)
	target := rand.Intn(2) % 2
	return target
}

func sleepRandomTime() {
	duration := time.Duration(computeBinaryInt() + 1)
	time.Sleep(duration * time.Nanosecond)
}
