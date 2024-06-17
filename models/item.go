package models

import (
	"encoding/json"
	"log"
	"os"
)

// represents rock, scissors, pepper and maybe something custom with they weights
// it's not a DB model
type Item struct {
	Name  string `json:"name"`
	Value uint   `json:"value"` // min value is always 0
}

type ResultType uint

const (
	Loss ResultType = iota
	Equal
	Winning
)

var (
	// we set them up on app init
	sessionMaxValue uint
	Items           []*Item
	ItemsByName     map[string]*Item
	ItemsByValue    map[uint]*Item
)

// func implements core logic of the RPS game
// we can use iota, but you want extandable values
func (first *Item) Beats(second *Item) ResultType {
	if first.Value == second.Value {
		return Equal
	}
	if first.Value == 0 && second.Value == sessionMaxValue {
		return Winning
	}
	if first.Value > second.Value {
		return Winning
	}
	return Loss
}

func init() {
	load()
	setMaxValue()
	setItemsNameMap()
	setItemsValueMap()
}

// loads items from json file
func load() {
	cfgFile, err := os.ReadFile("config/items.json")
	if err != nil {
		log.Fatalf("init err %v ", err)
	}
	err = json.Unmarshal(cfgFile, &Items)
	if err != nil {
		log.Fatalf("Unmarshal err: %v", err)
	}
}

// we can get value of last elementh of json, but in this case we have to be sure that json elements are ordered
func setMaxValue() {
	for _, item := range Items {
		if item.Value > sessionMaxValue {
			sessionMaxValue = item.Value
		}
	}
}

// for better performance
func setItemsNameMap() {
	ItemsByName = make(map[string]*Item)
	for _, item := range Items {
		ItemsByName[item.Name] = item
	}
}

func setItemsValueMap() {
	ItemsByValue = make(map[uint]*Item)
	for _, item := range Items {
		ItemsByValue[item.Value] = item
	}
}
