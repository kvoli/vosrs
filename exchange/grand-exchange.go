package exchange

import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
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
