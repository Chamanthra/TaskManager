package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nTask Manager CLI")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. List Tasks")
		fmt.Println("4. Create Task")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option: ")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			login(reader)
		case "2":
			register(reader)
		case "3":
			listTasks()
		case "4":
			createTask(reader)
		case "5":
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			color.Red("Invalid option. Please try again.")
		}
	}
}

func login(reader *bufio.Reader) {
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Call API to login
	color.Green("Login successful (mock)")
}

func register(reader *bufio.Reader) {
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Role (user/admin): ")
	role, _ := reader.ReadString('\n')
	role = strings.TrimSpace(role)

	// Call API to register
	color.Green("Registration successful (mock)")
}

func listTasks() {
	// Call API to get tasks
	fmt.Println("\nYour Tasks:")
	fmt.Println("1. Complete project - High priority - Due: 2023-12-31")
	fmt.Println("2. Buy groceries - Medium priority - Due: 2023-11-15")
}

func createTask(reader *bufio.Reader) {
	fmt.Print("Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("Priority (low/medium/high): ")
	priority, _ := reader.ReadString('\n')
	priority = strings.TrimSpace(priority)

	fmt.Print("Due Date (YYYY-MM-DD): ")
	dueDate, _ := reader.ReadString('\n')
	dueDate = strings.TrimSpace(dueDate)

	// Call API to create task
	color.Green("Task created successfully (mock)")
}
