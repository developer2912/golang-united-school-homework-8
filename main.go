package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

type Item struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Arguments map[string]string

const (
	ID_FLAG        = "id"
	OPERATION_FLAG = "operation"
	ITEM_FLAG      = "item"
	FILENAME_FLAG  = "fileName"

	ADD_ITEM       = "add"
	REMOVE_ITEM    = "remove"
	GET_ALL_ITEMS  = "list"
	GET_ITEM_BY_ID = "findById"

	errorTextMessageTemplate1 = "-%s flag has to be specified"
	errorTextMessageTemplate2 = "Operation %s not allowed!"
	errorTextMessageTemplate3 = "Item with id %s already exists"
	errorTextMessageTemplate4 = "Item with id %s not found"
)

var (
	validOperations           = []string{"add", "remove", "findById", "list"}
	validCmdFlags             = []string{OPERATION_FLAG, FILENAME_FLAG, ITEM_FLAG, ID_FLAG}
	defaultCmdFlagValue       = ""
	defaultCmdFlagDescription = ""
)

func parseArgs() Arguments {
	cmdFlags := make([]string, 4)
	for i, cmdflag := range validCmdFlags {
		flag.StringVar(&cmdFlags[i], cmdflag, defaultCmdFlagValue, defaultCmdFlagDescription)
	}
	flag.Parse()

	arguments := make(Arguments)
	for i, cmdflag := range validCmdFlags {
		if cmdFlags[i] != defaultCmdFlagValue {
			arguments[cmdflag] = cmdFlags[i]
		}
	}
	return arguments
}

func Perform(args Arguments, writer io.Writer) error {
	operation, fileName := "", ""
	for _, cmdflag := range validCmdFlags {
		switch cmdflag {
		case OPERATION_FLAG:
			value, existed := args[OPERATION_FLAG]
			if !existed {
				return fmt.Errorf(errorTextMessageTemplate1, OPERATION_FLAG)
			} else if value == "" {
				return fmt.Errorf(errorTextMessageTemplate1, OPERATION_FLAG)
			}
			isOperationValid := false
			for i := range validOperations {
				if validOperations[i] == value {
					isOperationValid = true
					break
				}
			}
			if !isOperationValid {
				return fmt.Errorf(errorTextMessageTemplate2, value)
			}
			operation = value
		case FILENAME_FLAG:
			value, existed := args[cmdflag]
			if !existed {
				return fmt.Errorf(errorTextMessageTemplate1, cmdflag)
			} else if value == "" {
				return fmt.Errorf(errorTextMessageTemplate1, cmdflag)
			}
			fileName = value
		default:
		}
	}
	switch operation {
	case GET_ALL_ITEMS:
		items, err := ReadAllItems(fileName)
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(items)
		if err != nil {
			return err
		}
		if len(bytes) != 0 {
			_, err := writer.Write(bytes)
			if err != nil {
				return err
			}
		}
	case ADD_ITEM:
		item, err := ParseItemCmdFlag(args)
		if err != nil {
			return err
		}
		items, err := ReadAllItems(fileName)
		if err != nil {
			return err
		}
		if _, itemFound := FindItemById(item.Id, items); itemFound {
			errorMessage := fmt.Sprintf(errorTextMessageTemplate3, item.Id)
			writer.Write([]byte(errorMessage))
		} else {
			items = append(items, item)
			WriteAllItems(items, fileName)
		}
	case REMOVE_ITEM:
		id, err := ParseIdCmdFlag(args)
		if err != nil {
			return err
		}
		items, err := ReadAllItems(fileName)
		if err != nil {
			return err
		}
		index, itemFound := FindItemById(id, items)
		if !itemFound {
			message := fmt.Sprintf(errorTextMessageTemplate4, id)
			writer.Write([]byte(message))
		} else {
			items = append(items[:index], items[index+1:]...)
			WriteAllItems(items, fileName)
		}
	case GET_ITEM_BY_ID:
		id, err := ParseIdCmdFlag(args)
		if err != nil {
			return err
		}
		items, err := ReadAllItems(fileName)
		if err != nil {
			return err
		}
		index, itemFound := FindItemById(id, items)
		if itemFound {
			bytes, err := json.Marshal(items[index])
			if err != nil {
				return err
			}
			writer.Write(bytes)
		}
	default:
	}
	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
