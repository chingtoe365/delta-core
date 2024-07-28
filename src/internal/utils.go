package internal

import (
	"delta-core/domain"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ReadTradeItemsFromJsonFile() domain.TradeItems {
	jsonFile, err := os.Open("assets/trade_items.json")
	if err != nil {
		fmt.Print(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var tradeItems domain.TradeItems
	json.Unmarshal(byteValue, &tradeItems)
	return tradeItems
}

func ReadTradeSignalsFromJsonFile() domain.TradeSignals {
	jsonFile, err := os.Open("assets/trade_signals.json")
	if err != nil {
		fmt.Print(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var tradeSignals domain.TradeSignals
	json.Unmarshal(byteValue, &tradeSignals)
	return tradeSignals
}

func ReadTradeSignalCategoriesFromJsonFile() domain.TradeSignalCategories {
	jsonFile, err := os.Open("assets/trade_signal_categories.json")
	if err != nil {
		fmt.Print(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var tradeSignalCategories domain.TradeSignalCategories
	json.Unmarshal(byteValue, &tradeSignalCategories)
	return tradeSignalCategories
}

func GetAllTradeItems() []domain.TradeItem {
	tradeItems := ReadTradeItemsFromJsonFile()
	return tradeItems.Items
}

func GetTradeItemByName(itemName string) (domain.TradeItem, error) {
	tradeItems := GetAllTradeItems()
	for _, item := range tradeItems {
		if item.Name == itemName {
			return item, nil
		}
	}
	return domain.TradeItem{}, os.ErrNotExist
}

func GetTradeItemByCategory(category string) []domain.TradeItem {
	tradeItems := GetAllTradeItems()
	var results []domain.TradeItem
	for _, item := range tradeItems {
		if item.Category == category {
			results = append(results, item)
		}
	}
	return results
}

func GetAllTradeSignals() []domain.TradeSignal {
	tradeSignals := ReadTradeSignalsFromJsonFile()
	return tradeSignals.Items
}

func GetSignalByName(name string) (domain.TradeSignal, error) {
	tradeSignals := GetAllTradeSignals()
	for _, signal := range tradeSignals {
		if signal.Name == name {
			return signal, nil
		}
	}
	return domain.TradeSignal{}, os.ErrNotExist
}

func EscapeTradeItemName(itemName string) string {
	title := strings.ReplaceAll(itemName, " ", "-")
	title = strings.ReplaceAll(title, "/", "-")
	title = strings.ToLower(title)
	return title
}

func BuildTaskTitle(item domain.TradeItem, signal domain.TradeSignal) string {
	return fmt.Sprintf("%s/%s", EscapeTradeItemName(item.Name), signal.Name)
}

func ParseFromTitle(title string) (domain.TradeItem, domain.TradeSignal, error) {
	// TODO: stablize solution
	splits := strings.Split(title, "/")
	itemName := splits[0]
	signalName := splits[1]
	itemSplits := strings.Split(itemName, "-")
	tradeItems := GetAllTradeItems()
	tradeSignals := GetAllTradeSignals()
	var foundItem domain.TradeItem
	var foundSignal domain.TradeSignal
	itemFound := false
	signalFound := false
	for _, item := range tradeItems {
		if strings.HasPrefix(strings.ToLower(item.Name), itemSplits[0]) {
			foundItem = item
			itemFound = true
			break
		}
	}
	for _, signal := range tradeSignals {
		if signal.Name == signalName {
			foundSignal = signal
			signalFound = true
			break
		}
	}
	if itemFound && signalFound {
		return foundItem, foundSignal, nil
	} else {
		return foundItem, foundSignal, os.ErrNotExist
	}

}

func FilterTasks(tradeItem string, tradeItemCategory string, tradeSignal string, tasks []domain.Task) []domain.Task {
	var tasksOut []domain.Task
	if (strings.EqualFold(tradeItem, "ALL") || strings.EqualFold(tradeItemCategory, "ALL")) && strings.EqualFold(tradeSignal, "ALL") {
		tasksOut = tasks
	} else {
		for _, task := range tasks {
			item, signal, err := ParseFromTitle(task.Title)
			if err != nil {
				return nil
			}
			itemEscaped := EscapeTradeItemName(strings.ToLower(tradeItem))
			titleLowerCased := strings.ToLower(task.Title)
			if strings.EqualFold(tradeSignal, "ALL") {
				if tradeItem != "" {
					if strings.HasPrefix(titleLowerCased, itemEscaped) {
						tasksOut = append(tasksOut, task)
					}
				} else {
					if tradeItemCategory != "" && strings.EqualFold(tradeItemCategory, item.Category) {
						tasksOut = append(tasksOut, task)
					}
				}
			} else {
				if tradeSignal != "" && strings.EqualFold(tradeSignal, signal.Name) {
					if strings.EqualFold(tradeItem, "ALL") || strings.EqualFold(tradeItemCategory, "ALL") {
						tasksOut = append(tasksOut, task)
					} else {
						if tradeItem != "" {
							if strings.HasPrefix(titleLowerCased, itemEscaped) {
								tasksOut = append(tasksOut, task)
							}
						} else {
							if tradeItemCategory != "" && strings.EqualFold(tradeItemCategory, item.Category) {
								tasksOut = append(tasksOut, task)
							}
						}
					}
				}
			}
		}
	}
	return tasksOut
}
