package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ParseItemCmdFlag(args Arguments) (item Item, err error) {
	value, existed := args[ITEM_FLAG]
	if !existed {
		return item, fmt.Errorf(errorTextMessageTemplate1, ITEM_FLAG)
	} else if value == "" {
		return item, fmt.Errorf(errorTextMessageTemplate1, ITEM_FLAG)
	}
	if err = json.Unmarshal([]byte(value), &item); err != nil {
		return item, err
	}
	return item, nil
}

func ParseIdCmdFlag(args Arguments) (string, error) {
	itemId, existed := args[ID_FLAG]
	if !existed {
		return "", fmt.Errorf(errorTextMessageTemplate1, ID_FLAG)
	} else if itemId == "" {
		return "", fmt.Errorf(errorTextMessageTemplate1, ID_FLAG)
	}
	return itemId, nil
}

func ReadAllItems(filePath string) ([]Item, error) {
	items := make([]Item, 0)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(bytes) != 0 {
		if err := json.Unmarshal(bytes, &items); err != nil {
			return nil, err
		}
	}
	return items, nil
}

func FindItemById(itemId string, items []Item) (int, bool) {
	for i := range items {
		if items[i].Id == itemId {
			return i, true
		}
	}
	return 0, false
}

func WriteAllItems(items []Item, filePath string) error {
	file, err := os.Create(filePath)
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(items)
	if err != nil {
		return err
	}
	file.Write(bytes)
	return nil
}
