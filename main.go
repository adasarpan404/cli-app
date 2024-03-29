package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"sort"
)

type Expense struct {
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

type ExpenseTracker struct {
	Expenses []Expense `json:"expenses"`
}

func main() {
	expenseTracker := loadExpensesFromFile("expenses.dat")
	action := flag.String("action", "view", "Specify the action: Add, View or summarize")
	name := flag.String("name", "", "Specify the name of your expense")
	cost := flag.Float64("cost", 0, "Specify the cost of the expense")
	minCost := flag.Float64("minCost", 0, "Specify the min cost for searching")
	maxCost := flag.Float64("maxCost", 0, "Specify the max cost for searching")
	sortOrder := flag.String("sortOrder", "asc", "Specify the sort order: asc or desc")
	flag.Parse()
	switch *action {
	case "add":
		expenseTracker.AddExpense(*name, *cost)
	case "view":
		expenseTracker.SortByCost(*sortOrder)
		expenseTracker.ViewExpenses()
	case "summarize":
		expenseTracker.SortByCost(*sortOrder)
		expenseTracker.SummarizeExpenses()
	case "search":
		expenseTracker.SearchExpenses(*name, *minCost, *maxCost)
	case "sort":
		expenseTracker.SortByCost(*sortOrder)
	default:
		fmt.Println("Invalid action. Use 'add', 'view', or 'summarize'.")
		os.Exit(1)
	}
	saveExpensesToFile(expenseTracker, "expenses.dat")
}

func loadExpensesFromFile(filename string) ExpenseTracker {

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		emptyFile, createErr := os.Create(filename)
		if createErr != nil {
			fmt.Println("Error creating expenses file:", createErr)
			return ExpenseTracker{}
		}
		defer emptyFile.Close()
	} else if err != nil {
		fmt.Println("Error checking expenses file:", err)
		return ExpenseTracker{}
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading expenses file:", err)
		return ExpenseTracker{}
	}

	var expenseTracker ExpenseTracker
	if len(data) > 0 {
		reader := bytes.NewReader(data)
		decoder := gob.NewDecoder(reader)
		err = decoder.Decode(&expenseTracker)
		if err != nil {
			fmt.Println("Error decoding expenses data:", err)
			return ExpenseTracker{}
		}
	}

	return expenseTracker
}

func saveExpensesToFile(expenseTracker ExpenseTracker, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating expenses file:", err)
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(expenseTracker)
	if err != nil {
		fmt.Println("Error encoding expenses data:", err)
	}
}

func (et *ExpenseTracker) AddExpense(name string, cost float64) {
	if name == "" || cost <= 0 {
		fmt.Println("Invalid expense details. Both name and cost are required.")
		return
	}
	newExpense := Expense{Name: name, Cost: cost}
	et.Expenses = append(et.Expenses, newExpense)
	fmt.Printf("Expense added: %s ($%.2f)\n", newExpense.Name, newExpense.Cost)
}

func (et *ExpenseTracker) ViewExpenses() {
	if len(et.Expenses) == 0 {
		fmt.Println("No expenses recorded.")
		return
	}

	fmt.Println("All Expenses:")
	for _, expense := range et.Expenses {
		fmt.Printf("%s: $%.2f\n", expense.Name, expense.Cost)
	}
}

func (et *ExpenseTracker) SummarizeExpenses() {
	if len(et.Expenses) == 0 {
		fmt.Println("No expenses recorded.")
		return
	}

	totalExpense := 0.0
	for _, expense := range et.Expenses {
		totalExpense += expense.Cost
	}

	fmt.Printf("Total Expenses: $%.2f\n", totalExpense)
}

func (et *ExpenseTracker) SearchExpenses(name string, minCost, maxCost float64) {
	if name == "" && minCost <= 0 && maxCost <= 0 {
		fmt.Println("Invalid search key. Please provide")
		return
	}
	found := false
	for _, expense := range et.Expenses {
		nameMatch := name == "" || expense.Name == name

		// Check if the cost is within the specified range (if provided)
		costInRange := (minCost <= 0 || expense.Cost >= minCost) &&
			(maxCost <= 0 || expense.Cost <= maxCost)

		if nameMatch && costInRange {
			fmt.Printf("Expense found: %s ($%.2f)\n", expense.Name, expense.Cost)
			found = true
		}
	}
	if !found {
		fmt.Println("No expenses found matching the specified criteria.")
	}
}

func (et *ExpenseTracker) SortByCost(sortOrder string) {
	sort.Slice(et.Expenses, func(i, j int) bool {
		if sortOrder == "asc" {
			return et.Expenses[i].Cost < et.Expenses[j].Cost
		} else if sortOrder == "desc" {
			return et.Expenses[i].Cost > et.Expenses[j].Cost
		}
		// Default to ascending order if sortOrder is not "asc" or "desc"
		return et.Expenses[i].Cost < et.Expenses[j].Cost
	})
}
