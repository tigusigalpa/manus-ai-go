package main

import (
	"fmt"
	"log"
	"os"

	manusai "github.com/tigusigalpa/manus-ai-go"
)

func main() {
	apiKey := os.Getenv("MANUS_AI_API_KEY")
	if apiKey == "" {
		log.Fatal("MANUS_AI_API_KEY environment variable is required")
	}

	client, err := manusai.NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("=== Creating a Task ===")
	task, err := client.CreateTask("Write a poem about Go programming", &manusai.TaskOptions{
		AgentProfile: manusai.AgentProfileManus16,
		TaskMode:     "chat",
	})
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	fmt.Printf("Task created successfully!\n")
	fmt.Printf("Task ID: %s\n", task.TaskID)
	fmt.Printf("Task Title: %s\n", task.TaskTitle)
	fmt.Printf("Task URL: %s\n", task.TaskURL)

	fmt.Println("\n=== Getting Task Details ===")
	taskDetail, err := client.GetTask(task.TaskID)
	if err != nil {
		log.Fatalf("Failed to get task: %v", err)
	}

	fmt.Printf("Status: %s\n", taskDetail.Status)
	fmt.Printf("Credit Usage: %.2f\n", taskDetail.CreditUsage)

	fmt.Println("\n=== Listing Tasks ===")
	tasks, err := client.GetTasks(&manusai.TaskFilters{
		Limit: 10,
		Order: "desc",
	})
	if err != nil {
		log.Fatalf("Failed to list tasks: %v", err)
	}

	fmt.Printf("Found %d tasks\n", len(tasks.Data))
	for i, t := range tasks.Data {
		fmt.Printf("%d. %s - %s (Status: %s)\n", i+1, t.ID, t.Title, t.Status)
	}

	fmt.Println("\n=== Updating Task ===")
	newTitle := "Updated Task Title"
	updatedTask, err := client.UpdateTask(task.TaskID, &manusai.TaskUpdate{
		Title: &newTitle,
	})
	if err != nil {
		log.Fatalf("Failed to update task: %v", err)
	}

	fmt.Printf("Task updated: %s\n", updatedTask.Title)

	fmt.Println("\n=== Deleting Task ===")
	deleteResult, err := client.DeleteTask(task.TaskID)
	if err != nil {
		log.Fatalf("Failed to delete task: %v", err)
	}

	fmt.Printf("Task deleted: %v\n", deleteResult.Deleted)
}
