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

	fmt.Println("=== Creating File Record ===")
	fileResult, err := client.CreateFile("document.pdf")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}

	fmt.Printf("File ID: %s\n", fileResult.ID)
	fmt.Printf("Filename: %s\n", fileResult.Filename)
	fmt.Printf("Status: %s\n", fileResult.Status)

	fmt.Println("\n=== Uploading File Content ===")
	fileContent := []byte("This is sample PDF content")
	err = client.UploadFileContent(fileResult.UploadURL, fileContent, "application/pdf")
	if err != nil {
		log.Fatalf("Failed to upload file content: %v", err)
	}

	fmt.Println("File uploaded successfully!")

	fmt.Println("\n=== Creating Task with File Attachment ===")
	attachment := manusai.NewAttachmentFromFileID(fileResult.ID)
	
	task, err := client.CreateTask("Analyze this document", &manusai.TaskOptions{
		AgentProfile: manusai.AgentProfileManus16,
		Attachments:  []interface{}{attachment},
	})
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	fmt.Printf("Task created with attachment!\n")
	fmt.Printf("Task ID: %s\n", task.TaskID)
	fmt.Printf("Task URL: %s\n", task.TaskURL)

	fmt.Println("\n=== Creating Task with URL Attachment ===")
	urlAttachment := manusai.NewAttachmentFromURL("https://example.com/image.jpg")
	
	task2, err := client.CreateTask("Describe this image", &manusai.TaskOptions{
		AgentProfile: manusai.AgentProfileManus16,
		Attachments:  []interface{}{urlAttachment},
	})
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	fmt.Printf("Task created with URL attachment!\n")
	fmt.Printf("Task ID: %s\n", task2.TaskID)

	fmt.Println("\n=== Creating Task from Local File ===")
	localAttachment, err := manusai.NewAttachmentFromFilePath("example.txt")
	if err != nil {
		fmt.Printf("Note: Could not load local file (this is expected in demo): %v\n", err)
	} else {
		task3, err := client.CreateTask("Process this file", &manusai.TaskOptions{
			AgentProfile: manusai.AgentProfileManus16,
			Attachments:  []interface{}{localAttachment},
		})
		if err != nil {
			log.Fatalf("Failed to create task: %v", err)
		}
		fmt.Printf("Task created from local file!\n")
		fmt.Printf("Task ID: %s\n", task3.TaskID)
	}

	fmt.Println("\n=== Listing Files ===")
	files, err := client.ListFiles()
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}

	fmt.Printf("Found %d files\n", len(files.Data))
	for i, f := range files.Data {
		fmt.Printf("%d. %s - %s (Status: %s)\n", i+1, f.ID, f.Filename, f.Status)
	}

	fmt.Println("\n=== Deleting File ===")
	deleteResult, err := client.DeleteFile(fileResult.ID)
	if err != nil {
		log.Fatalf("Failed to delete file: %v", err)
	}

	fmt.Printf("File deleted: %v\n", deleteResult.Deleted)
}
