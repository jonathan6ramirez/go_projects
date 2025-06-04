package main

import (
	"encoding/csv"
	"fmt"
	"sync"

	// "fmt"
	"log"
	"os"

	// importing colly
	"github.com/gocolly/colly"
)

type Product struct {
	Url, Image, Name, Price string
}

/* INFO:   This program is meant to practice scraping the internet. as of time of writing this it
   can scape a particular page but not do any kind of pagination.
   The general flow is to:
	   1. Initialize the collector object
	   2. Define the function that colects the data
	   3. Define the function that handles moving to the next page
	   4. Define the functions that outputs the data once its scraped
	   5. Lastly visit the website to collect the data.
   https://www.zenrows.com/blog/web-scraping-golang#headless-browser-go
*/

func main() {
	// Instantiate a new collector object
	c := colly.NewCollector(
		colly.AllowedDomains("www.scrapingcourse.com"),
	)

	// initialize the slice of structs that will contain the scraped data
	var products []Product

	// define a sync to filter visited URLs
	var visitedUrls sync.Map

	// onhtml callback for scraping product information
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		// initialize a new product instance
		product := Product{}

		// scrape the target data
		product.Url = e.ChildAttr("a", "href")
		product.Image = e.ChildAttr("img", "src")
		product.Name = e.ChildText(".product-name")
		product.Price = e.ChildText(".price")

		// add the product instance with scraped data to the list of products.
		products = append(products, product)
	})

	// onhtml callback for handling pagination
	c.OnHTML("a.next", func(e *colly.HTMLElement) {
		// extract the next page URL from the next button
		nextPage := e.Attr("href")

		//check if the nextPage URL has been visited
		if _, found := visitedUrls.Load(nextPage); !found {
			fmt.Println("scraping: ", nextPage)

			// mark the URL as visited
			visitedUrls.Store(nextPage, struct{}{})

			// visit the next page
			e.Request.Visit(nextPage)
		}
	})

	// store the data to a csv after extraction
	c.OnScraped(func(r *colly.Response) {
		// open the csv file
		file, err := os.Create("products.csv")
		if err != nil {
			log.Fatalf("Failed to create output csv file", err)
		}
		defer file.Close()

		// initialize a file writer
		writer := csv.NewWriter(file)

		// write the csv headers
		headers := []string{
			"Url",
			"Image",
			"Name",
			"Price",
		}
		writer.Write(headers)

		// write each product as a csv row
		for _, product := range products {
			// convert a product to an array or strings
			record := []string{
				product.Url,
				product.Image,
				product.Name,
				product.Price,
			}

			writer.Write(record)
		}
		defer writer.Flush()
	})

	// open the target URL
	c.Visit("https://www.scrapingcourse.com/ecommerce")
}
