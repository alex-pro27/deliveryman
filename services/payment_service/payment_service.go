package payment_service

import (
	"fmt"
	"git.samberi.com/dois/delivery_api/config"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	CANCELLED       = "CANCELLED"
	PAID            = "PAID"
	CREATED         = "CREATED"
	REGISTERED      = "REGISTERED"
	ERROR           = "ERROR"
	ORDER_NOT_FOUND = "ORDER_NOT_FOUND"
)

type Queue map[string]bool

type Customer struct {
	UID   string
	Phone string
}

type Product struct {
	Name          string
	Department    int32
	Price         float64
	Quantity      float64
	Measure       string
	Code          string
	Barcode       string
	TaxType       int32
	TaxRate       float64
	TaxSum        float64
	Discount      float64
	PriceDiscount float64
	Amount        float64
}

type Order struct {
	UID          string
	Products     []Product
	PayStatusURL string
	SuccessURL   string
	FailURL      string
	ShopNumber   uint
	Customer     Customer
}

func getUrl(path string) string {
	return config.Config.PaymentService.Server + path
}

func readBody(resp *http.Response) (body string, err error) {
	b, _ := ioutil.ReadAll(resp.Body)
	body = string(b)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("body: %s; status: %d", body, resp.StatusCode)
	}
	return body, err
}

func ConfirmOrder(order Order) error {
	httpClient := new(http.Client)
	resp, err := httpClient.PostForm(
		getUrl(config.Config.PaymentService.Methods.CreateOrder),
		url.Values{},
	)
	if err != nil {
		return err
	}
	_, err = readBody(resp)
	return err
}
