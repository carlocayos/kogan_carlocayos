package main

import (
	c "carlokogan/client"
	"context"
	"fmt"
	"github.com/magiconair/properties"
	"net/url"
	"os"
)

// application using the client and computing the average cubic weight
func main() {
	if len(os.Args) == 1 || len(os.Args) > 3 {
		fmt.Fprintf(os.Stderr, "\nInvalid program argument. See format `myapp <config-file-path> <[OPTIONAL] category filter>`\n")
		os.Exit(1)
	}

	// parse the config file
	confPath := os.Args[1]

	// get category filter if supplied, otherwise use the default value
	var categoryFilter string
	if len(os.Args) == 3 {
		categoryFilter = os.Args[2]
	} else {
		categoryFilter = DefaultCategoryFilter
	}

	// load the properties file
	p := properties.MustLoadFile(confPath, properties.UTF8)
	apiEndpointScheme := p.MustGetString(ApiEndpointScheme)
	apiEndpointHost := p.MustGetString(ApiEndpointHost)
	apiEndpointPath := p.MustGetString(ApiEndpointPath)

	// parse the request for errors
	u, err := url.ParseRequestURI(apiEndpointScheme + "://" + apiEndpointHost + apiEndpointPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid uri %s\n", err)
		os.Exit(1)
	}

	// create a new Kogan rest api client
	client, err := c.NewClient(u)
	if err != nil {
		panic(err)
	}

	// create a channel accepting ProductPage
	responseStream := make(chan c.ProductPage)

	// create a goroutine to concurrently send the request while processing the response
	go func() {
		for {
			productPage, err := client.GetProductPage(context.Background(), apiEndpointPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "get product page error %s\n", err)
				os.Exit(1)
			}

			// send the response to the stream channel for calculating the Cubic Weight
			responseStream <- *productPage

			if productPage.Next == "" {
				close(responseStream)
				break
			}
			apiEndpointPath = productPage.Next
		}
	}()

	// blocks and waits for response stream
	var cubicWeight []float64
	for resp := range responseStream {

		// filter the products belonging to the defined Category
		for _, val := range resp.Products {

			if val.Category != categoryFilter {
				continue
			}

			// calculate the Cubic Weight for each product and append the result
			// cubic weight = length * height * width * (250 Cubic Weight Conversion Factor) * (0.01 convert cm to m)
			cWeight := (val.Size.Length * val.Size.Height * val.Size.Width) * CubicWeightConversionFactor * .01
			cubicWeight = append(cubicWeight, cWeight)
		}
	}

	// compute for the Average Cubic Weight
	var sum float64
	for _, val := range cubicWeight {
		sum += val
	}

	// Compute for the total average cubic weight
	average := sum / float64(len(cubicWeight))
	fmt.Printf("\nAverage Cubic Weight of ALL Products under the Category \"%s\" = %.2f kg\n", categoryFilter, average)
}
