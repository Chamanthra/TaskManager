package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	baseURL = "http://localhost:8080/api"
)

var (
	token  string
	client = &http.Client{}
	reader = bufio.NewReader(os.Stdin)
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

func main() {
	checkServerConnection()

	for {
		fmt.Printf("\n%s\n", blue("Task Manager CLI"))
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. Exit")
		fmt.Print("Choose an option: ")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			login()
		case "2":
			register()
		case "3":
			fmt.Println(green("Goodbye!"))
			os.Exit(0)
		default:
			fmt.Println(red("Invalid option. Please try again."))
		}
	}
}

func checkServerConnection() {
	resp, err := http.Get(baseURL + "/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println(red("Server is not available. Please ensure the server is running."))
		os.Exit(1)
	}
	fmt.Println(green("Connected to server successfully"))
}

func login() {
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	credentials := map[string]string{
		"user_name": username,
		"password":  password,
	}

	resp, err := sendRequest("POST", "/login", credentials, nil)
	if err != nil {
		fmt.Println(red("Login failed:", err))
		return
	}

	respMap, ok := resp.(map[string]interface{})
	if !ok {
		fmt.Println(red("Invalid login response format"))
		return
	}

	tokenData, ok := respMap["token"].(string)
	if !ok {
		fmt.Println(red("Token missing in response"))
		return
	}

	token = tokenData
	fmt.Println(green("\nLogin successful!"))

	role, _ := respMap["role"].(string)
	if role == "admin" {
		adminDashboardMenu()
	} else {
		postLoginMenu()
	}
}

func register() {
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	userData := map[string]string{
		"user_name": username,
		"password":  password,
	}

	_, err := sendRequest("POST", "/register", userData, nil)
	if err != nil {
		fmt.Println(red("Registration failed:", err))
		return
	}

	fmt.Println(green("\nRegistration successful! Please login."))
}

func postLoginMenu() {
	// Get user role first
	resp, err := sendRequest("GET", "/profile", nil, nil)
	if err != nil {
		fmt.Println(red("Failed to get user profile:", err))
		return
	}

	profile, ok := resp.(map[string]interface{})
	if !ok {
		fmt.Println(red("Invalid profile data"))
		return
	}

	role := profile["role"].(string)

	for {
		fmt.Printf("\n%s\n", blue("Task Manager - Logged In"))
		fmt.Println("1. Task Management")
		fmt.Println("2. Comment Management")
		fmt.Println("3. File Management")
		fmt.Println("4. Notifications")
		fmt.Println("5. User Profile")

		if role == "admin" {
			fmt.Println("6. Admin Dashboard")
			fmt.Println("7. Logout")
		} else {
			fmt.Println("6. Logout")
		}

		fmt.Print("Choose an option: ")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			taskManagementMenu()
		case "2":
			commentManagementMenu()
		case "3":
			fileManagementMenu()
		case "4":
			notificationMenu()
		case "5":
			userProfileMenu()
		case "6":
			if role == "admin" {
				adminDashboardMenu()
			} else {
				token = ""
				fmt.Println(green("Logged out successfully"))
				return
			}
		case "7":
			if role == "admin" {
				token = ""
				fmt.Println(green("Logged out successfully"))
				return
			} else {
				fmt.Println(red("Invalid option"))
			}
		default:
			fmt.Println(red("Invalid option"))
		}
	}
}

func adminDashboardMenu() {
	for {
		fmt.Printf("\n%s\n", blue("Admin Dashboard"))
		fmt.Println("1. List All Users")
		fmt.Println("2. Delete User")
		fmt.Println("3. View All Tasks")
		fmt.Println("4. Back to Main Menu")
		fmt.Print("Choose an option: ")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			listAllUsers()
		case "2":
			deleteUser()
		case "3":
			viewAllTasks()
		case "4":
			return
		default:
			fmt.Println(red("Invalid option"))
		}
	}
}

func listAllUsers() {
	resp, err := sendRequest("GET", "/users", nil, nil)
	if err != nil {
		fmt.Println(red("Failed to get users:", err))
		return
	}

	if users, ok := resp.([]interface{}); ok {
		fmt.Println(green("\nAll Users:"))
		for i, user := range users {
			userMap := user.(map[string]interface{})
			fmt.Printf("%d. %s (Role: %s)\n", i+1, userMap["username"], userMap["role"])
		}
	}
}

func deleteUser() {
	fmt.Print("Enter User ID to delete: ")
	userID, _ := reader.ReadString('\n')
	userID = strings.TrimSpace(userID)

	fmt.Print(red("Are you sure you want to delete this user? (y/n): "))
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm == "y" || confirm == "yes" {
		_, err := sendRequest("DELETE", "/users/"+userID, nil, nil)
		if err != nil {
			fmt.Println(red("Failed to delete user:", err))
		}
	} else {
		fmt.Println(yellow("User deletion cancelled"))
	}
}

func viewAllTasks() {
	resp, err := sendRequest("GET", "/tasks/all", nil, nil)
	if err != nil {
		fmt.Println(red("Failed to get tasks:", err))
		return
	}

	if tasks, ok := resp.([]interface{}); ok {
		fmt.Println(green("\nAll Tasks:"))
		for i, task := range tasks {
			taskMap := task.(map[string]interface{})
			fmt.Printf("%d. %s (Status: %s, Owner: %s)\n",
				i+1, taskMap["title"], taskMap["status"], taskMap["user"].(map[string]interface{})["username"])
			if due, ok := taskMap["due_date"].(string); ok && due != "" {
				fmt.Printf("   Due: %s\n", due)
			}
		}
	}
}

func taskManagementMenu() {
	for {
		fmt.Printf("\n%s\n", blue("Task Management"))
		fmt.Println("1. Create Task")
		fmt.Println("2. List Tasks")
		fmt.Println("3. Update Task")
		fmt.Println("4. Delete Task")
		fmt.Println("5. Back to Main Menu")
		fmt.Print("Choose an option: ")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			createTask()
		case "2":
			listTasks()
		case "3":
			updateTask()
		case "4":
			deleteTask()
		case "5":
			return
		default:
			fmt.Println(red("Invalid option"))
		}
	}
}

func createTask() {
	fmt.Print("Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("Priority (low/medium/high/critical): ")
	priority, _ := reader.ReadString('\n')
	priority = strings.TrimSpace(priority)

	fmt.Print("Due Date (YYYY-MM-DD): ")
	dueDateStr, _ := reader.ReadString('\n')
	dueDateStr = strings.TrimSpace(dueDateStr)

	var dueDate time.Time
	if dueDateStr != "" {
		var err error
		dueDate, err = time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			fmt.Println(red("Invalid date format. Use YYYY-MM-DD"))
			return
		}
	}

	task := map[string]interface{}{
		"title":       title,
		"description": description,
		"priority":    priority,
	}

	if !dueDate.IsZero() {
		task["due_date"] = dueDate.Format(time.RFC3339)
	}

	_, err := sendRequest("POST", "/tasks", task, nil)
	if err != nil {
		fmt.Println(red("Failed to create task:", err))
	}
}

func listTasks() {
	resp, err := sendRequest("GET", "/tasks", nil, nil)
	if err != nil {
		fmt.Println(red("Failed to get tasks:", err))
		return
	}

	if tasks, ok := resp.([]interface{}); ok {
		fmt.Println(green("\nYour Tasks:"))
		for i, task := range tasks {
			taskMap := task.(map[string]interface{})
			fmt.Printf("%d. %s (Status: %s, Priority: %s)\n", i+1, taskMap["title"], taskMap["status"], taskMap["priority"])
			if due, ok := taskMap["due_date"].(string); ok && due != "" {
				fmt.Printf("   Due: %s\n", due)
			}
		}
	}
}

func updateTask() {
	fmt.Print("Enter Task ID to update: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	// Get current task first
	fmt.Println(yellow("\nGetting current task details..."))
	resp, err := sendRequest("GET", "/tasks/"+id, nil, nil)
	if err != nil {
		fmt.Println(red("Failed to get task:", err))
		return
	}

	task, ok := resp.(map[string]interface{})
	if !ok {
		fmt.Println(red("Invalid task data"))
		return
	}

	// Display current task info
	fmt.Println(green("Current Task:"))
	fmt.Printf("Title: %s\n", task["title"])
	fmt.Printf("Description: %s\n", task["description"])
	fmt.Printf("Status: %s\n", task["status"])
	fmt.Printf("Due Date: %s\n", task["due_date"])
	fmt.Printf("Priority: %s\n", task["priority"])

	// Get updates
	fmt.Print("\nNew Title (leave blank to keep current): ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("New Description (leave blank to keep current): ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("New Status (todo/in_progress/done/archived): ")
	status, _ := reader.ReadString('\n')
	status = strings.TrimSpace(status)

	fmt.Print("New Due Date (YYYY-MM-DD, leave blank to keep current): ")
	dueDateStr, _ := reader.ReadString('\n')
	dueDateStr = strings.TrimSpace(dueDateStr)

	updates := make(map[string]interface{})
	if title != "" {
		updates["title"] = title
	}
	if description != "" {
		updates["description"] = description
	}
	if status != "" {
		updates["status"] = status
	}
	if dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			fmt.Println(red("Invalid date format. Use YYYY-MM-DD"))
			return
		}
		updates["due_date"] = dueDate.Format(time.RFC3339)
	}

	if len(updates) == 0 {
		fmt.Println(yellow("No updates provided"))
		return
	}

	_, err = sendRequest("PUT", "/tasks/"+id, updates, nil)
	if err != nil {
		fmt.Println(red("Failed to update task:", err))
	}
}

func deleteTask() {
	fmt.Print("Enter Task ID to delete: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	fmt.Print(red("Are you sure you want to delete this task? (y/n): "))
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm == "y" || confirm == "yes" {
		_, err := sendRequest("DELETE", "/tasks/"+id, nil, nil)
		if err != nil {
			fmt.Println(red("Failed to delete task:", err))
		}
	} else {
		fmt.Println(yellow("Task deletion cancelled"))
	}
}

func commentManagementMenu() {
	for {
		fmt.Printf("\n%s\n", blue("Comment Management"))
		fmt.Println("1. Add Comment")
		fmt.Println("2. View Comments")
		fmt.Println("3. Delete Comment")
		fmt.Println("4. Back to Main Menu")
		fmt.Print("Choose an option: ")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			addComment()
		case "2":
			viewComments()
		case "3":
			deleteComment()
		case "4":
			return
		default:
			fmt.Println(red("Invalid option"))
		}
	}
}

func addComment() {
	fmt.Print("Enter Task ID: ")
	taskID, _ := reader.ReadString('\n')
	taskID = strings.TrimSpace(taskID)

	fmt.Print("Enter Comment: ")
	content, _ := reader.ReadString('\n')
	content = strings.TrimSpace(content)

	comment := map[string]string{
		"content": content,
	}

	_, err := sendRequest("POST", "/tasks/"+taskID+"/comments", comment, nil)
	if err != nil {
		fmt.Println(red("Failed to add comment:", err))
	}
}

func viewComments() {
	fmt.Print("Enter Task ID: ")
	taskID, _ := reader.ReadString('\n')
	taskID = strings.TrimSpace(taskID)

	resp, err := sendRequest("GET", "/tasks/"+taskID+"/comments", nil, nil)
	if err != nil {
		fmt.Println(red("Failed to get comments:", err))
		return
	}

	if comments, ok := resp.([]interface{}); ok {
		fmt.Println(green("\nComments:"))
		for i, comment := range comments {
			commentMap := comment.(map[string]interface{})
			fmt.Printf("%d. %s\n", i+1, commentMap["content"])
			if user, ok := commentMap["user"].(map[string]interface{}); ok {
				fmt.Printf("   - By: %s\n", user["user_name"])
			}
		}
	}
}

func deleteComment() {
	fmt.Print("Enter Comment ID to delete: ")
	commentID, _ := reader.ReadString('\n')
	commentID = strings.TrimSpace(commentID)

	fmt.Print(red("Are you sure you want to delete this comment? (y/n): "))
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm == "y" || confirm == "yes" {
		_, err := sendRequest("DELETE", "/comments/"+commentID, nil, nil)
		if err != nil {
			fmt.Println(red("Failed to delete comment:", err))
		}
	} else {
		fmt.Println(yellow("Comment deletion cancelled"))
	}
}

func fileManagementMenu() {
	for {
		fmt.Printf("\n%s\n", blue("File Management"))
		fmt.Println("1. Upload File")
		fmt.Println("2. Download File")
		fmt.Println("3. Delete File")
		fmt.Println("4. Back to Main Menu")
		fmt.Print("Choose an option: ")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			uploadFile()
		case "2":
			downloadFile()
		case "3":
			deleteFile()
		case "4":
			return
		default:
			fmt.Println(red("Invalid option"))
		}
	}
}

func uploadFile() {
	fmt.Print("Enter Task ID: ")
	taskID, _ := reader.ReadString('\n')
	taskID = strings.TrimSpace(taskID)

	fmt.Print("Enter file path to upload: ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(red("Error opening file:", err))
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		fmt.Println(red("Error creating form file:", err))
		return
	}
	io.Copy(part, file)
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}

	_, err = sendRequest("POST", "/tasks/"+taskID+"/files", body, headers)
	if err != nil {
		fmt.Println(red("Failed to upload file:", err))
	}
}

func downloadFile() {
	fmt.Print("Enter File ID to download: ")
	fileID, _ := reader.ReadString('\n')
	fileID = strings.TrimSpace(fileID)

	fmt.Print("Enter destination path: ")
	destPath, _ := reader.ReadString('\n')
	destPath = strings.TrimSpace(destPath)

	req, err := http.NewRequest("GET", baseURL+"/files/"+fileID, nil)
	if err != nil {
		fmt.Println(red("Error creating request:", err))
		return
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(red("Request failed:", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(red("Error downloading file:", string(body)))
		return
	}

	out, err := os.Create(destPath)
	if err != nil {
		fmt.Println(red("Error creating file:", err))
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(red("Error saving file:", err))
		return
	}

	fmt.Println(green("File downloaded successfully to:", destPath))
}

func deleteFile() {
	fmt.Print("Enter File ID to delete: ")
	fileID, _ := reader.ReadString('\n')
	fileID = strings.TrimSpace(fileID)

	fmt.Print(red("Are you sure you want to delete this file? (y/n): "))
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm == "y" || confirm == "yes" {
		_, err := sendRequest("DELETE", "/files/"+fileID, nil, nil)
		if err != nil {
			fmt.Println(red("Failed to delete file:", err))
		}
	} else {
		fmt.Println(yellow("File deletion cancelled"))
	}
}

func notificationMenu() {
	resp, err := sendRequest("GET", "/notifications", nil, nil)
	if err != nil {
		fmt.Println(red("Failed to get notifications:", err))
		return
	}

	if notifications, ok := resp.([]interface{}); ok {
		fmt.Println(green("\nYour Notifications:"))
		for i, notification := range notifications {
			notifMap := notification.(map[string]interface{})
			status := notifMap["status"].(string)
			statusColor := green
			if status == "unread" {
				statusColor = yellow
			}
			fmt.Printf("%d. [%s] %s\n", i+1, statusColor(status), notifMap["message"])
			if task, ok := notifMap["task"].(map[string]interface{}); ok {
				fmt.Printf("   Related to task: %s\n", task["title"])
			}
		}

		fmt.Println("\n1. Mark notification as read")
		fmt.Println("2. Back to main menu")
		fmt.Print("Choose an option: ")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			markNotificationAsRead()
		case "2":
			return
		default:
			fmt.Println(red("Invalid option"))
		}
	}
}

func markNotificationAsRead() {
	fmt.Print("Enter Notification ID to mark as read: ")
	notifID, _ := reader.ReadString('\n')
	notifID = strings.TrimSpace(notifID)

	_, err := sendRequest("PUT", "/notifications/"+notifID+"/read", nil, nil)
	if err != nil {
		fmt.Println(red("Failed to mark notification as read:", err))
	} else {
		fmt.Println(green("Notification marked as read"))
	}
}

func userProfileMenu() {
	for {
		resp, err := sendRequest("GET", "/profile", nil, nil)
		if err != nil {
			fmt.Println(red("Failed to get profile:", err))
			return
		}

		profile, ok := resp.(map[string]interface{})
		if !ok {
			fmt.Println(red("Invalid profile data"))
			return
		}

		fmt.Printf("\n%s\n", blue("User Profile"))
		fmt.Printf("Username: %s\n", profile["username"])
		fmt.Printf("Role: %s\n", profile["role"])
		fmt.Printf("Member since: %s\n", profile["created_at"])

		fmt.Println("\n1. Update profile")
		fmt.Println("2. Change password")
		fmt.Println("3. Back to main menu")
		fmt.Print("Choose an option: ")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			updateProfile()
		case "2":
			changePassword()
		case "3":
			return
		default:
			fmt.Println(red("Invalid option"))
		}
	}
}

func updateProfile() {
	fmt.Print("New Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	updates := map[string]string{
		"username": username,
	}

	_, err := sendRequest("PUT", "/profile", updates, nil)
	if err != nil {
		fmt.Println(red("Failed to update profile:", err))
	} else {
		fmt.Println(green("Profile updated successfully"))
	}
}

func changePassword() {
	fmt.Print("Current Password: ")
	currentPass, _ := reader.ReadString('\n')
	currentPass = strings.TrimSpace(currentPass)

	fmt.Print("New Password: ")
	newPass, _ := reader.ReadString('\n')
	newPass = strings.TrimSpace(newPass)

	updates := map[string]string{
		"current_password": currentPass,
		"new_password":     newPass,
	}

	_, err := sendRequest("PUT", "/profile/password", updates, nil)
	if err != nil {
		fmt.Println(red("Failed to change password:", err))
	} else {
		fmt.Println(green("Password changed successfully"))
	}
}

func sendRequest(method, path string, body interface{}, headers map[string]string) (interface{}, error) {
	var req *http.Request
	var err error

	if body != nil {
		var jsonData []byte
		switch v := body.(type) {
		case *bytes.Buffer:
			// For file uploads, use the buffer directly
			req, err = http.NewRequest(method, baseURL+path, v)
		default:
			jsonData, err = json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("error marshaling request body: %v", err)
			}
			req, err = http.NewRequest(method, baseURL+path, bytes.NewBuffer(jsonData))
		}
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}
	} else {
		req, err = http.NewRequest(method, baseURL+path, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}
	}

	// Add authorization header if token exists
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Set content type for JSON requests
	if body != nil && headers == nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add custom headers if provided
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("server returned %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse JSON response if content type is JSON
	if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		var result interface{}
		if err := json.Unmarshal(responseBody, &result); err != nil {
			return nil, fmt.Errorf("error parsing JSON response: %v", err)
		}
		return result, nil
	}

	return string(responseBody), nil
}
