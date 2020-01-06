package exchange

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"sync"
)

type ItemDetail struct {
	Item Item `json:"item"`
}

type Item struct {
	Icon string `json:"icon"`
	IconLarge string `json:"icon_large"`
	ID int `json:"id"`
	Type string `json:"type"`
	TypeIcon string `json:"typeIcon"`
	Name string `json:"name"`
	Description string `json:"description"`
	Members string `json:"members"`
	Current PriceCurrentDetail `json:"current"`
	Today PriceTodayDetail `json:"today"`
	Day30 PriceChange `json:"day30"`
	Day90 PriceChange `json:"day90"`
	Day180 PriceChange `json:"day180"`
}

type PriceChange struct {
	Trend string `json:"trend"`
	Change string `json:"change"`
}

type PriceCurrentDetail struct {
	Trend string `json:"trend"`
	Price float64 `json:"price"`
}

type PriceTodayDetail struct {
	Trend string `json:"trend"`
	Price string `json:"price"`
}

type ItemGraph struct {
	Daily map[string]int `json:"daily"`
	Average map[string]int `json:"average"`
}

type response struct {
	body []byte `json:"body"`
	err error `json:"error"`
}

// FetchItem returns the price, trend and misc info for an item
func FetchItem(itemID int) (*ItemDetail, error){
	uri := fmt.Sprintf("%s?item=%d", GE_DETAIL, itemID)

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &ItemDetail{}

	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}


// FetchGraph returns the daily and avg price chart for the previous 180 days for an item
func FetchGraph(itemID int) (*ItemGraph, error) {
	uri := fmt.Sprintf("%s/%d.json", GE_GRAPH, itemID)

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &ItemGraph{}

	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func newGraphURI(itemID int) string {
	return fmt.Sprintf("%s/%d.json", GE_GRAPH, itemID)
}

func newDetailURI(itemID int) string {
	return fmt.Sprintf("%s?item=%d", GE_DETAIL, itemID)
}

func mapResponse(body []byte, respStruct interface{}) error {
	return json.Unmarshal(body, respStruct)
}

func BatchFetchItem(itemIDs []int) ([]*ItemDetail, error) {
	uris := make([]string, len(itemIDs))
	for _, v := range itemIDs{
		uris = append(uris, newDetailURI(v))
	}

	res := batchGet(uris)

	itms := make([]*ItemDetail, 0)
	for _, v := range res {
		if v.err == nil {
			unstrct := &ItemDetail{}
			err := mapResponse(v.body, unstrct)
			if err == nil {
				itms = append(itms, unstrct)
			}
		}
	}

	return itms, nil
}

func BatchFetchGraph(itemIDs []int) ([]*ItemGraph, error) {
	uris := make([]string, len(itemIDs))
	for _, v := range itemIDs{
		uris = append(uris, newGraphURI(v))
	}

	res := batchGet(uris)
	
	itms := make([]*ItemGraph, 0)
	for _, v := range res {
		if v.err == nil {
			unstrct := &ItemGraph{}
			err := mapResponse(v.body, unstrct)
			if err == nil {
				itms = append(itms, unstrct)
			}
		}
	}

	return itms, nil
}

func batchGet(uris []string) ([]*response) {
	res := []*response{}
	chs := []chan *response{}

	for _, v := range uris {
		ch := make(chan *response)
		chs = append(chs, ch)
		go fetchData(ch, v)
	}

	for resp := range merge(chs) {
		res = append(res, resp)
	}

	return res
}

func fetchData(c chan<- *response, uri string) {
	defer close(c)

	resp, err := http.Get(uri)
	if err != nil {
		c <- &response{
			body: nil,
			err: err,
		}
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c <- &response{
			body: nil,
			err: err,
		}
	}

	c <- &response{
		body: body,
		err: nil,
	}
}


func merge(chs []chan *response) <-chan *response {
	var wg sync.WaitGroup

	out := make(chan *response)

	output := func(c <-chan *response) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(chs))

	for _, c := range chs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}


