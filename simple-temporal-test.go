package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// Simple workflow to test Temporal integration
func SimpleWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting simple workflow", "name", name)

	// Simulate some work
	workflow.Sleep(ctx, 2*time.Second)

	result := fmt.Sprintf("Hello %s! Workflow completed at %s", name, time.Now().Format(time.RFC3339))
	logger.Info("Workflow completed", "result", result)

	return result, nil
}

// Simple activity
func SimpleActivity(ctx context.Context, input string) (string, error) {
	log.Printf("Executing activity with input: %s", input)
	return fmt.Sprintf("Processed: %s", input), nil
}

func main() {
	log.Println("🚀 Testing Temporal Integration")
	log.Println("===============================")

	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort:  "localhost:7233",
		Namespace: "default",
	})
	if err != nil {
		log.Fatalf("❌ Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	log.Println("✅ Connected to Temporal server")

	// Create worker
	w := worker.New(c, "test-task-queue", worker.Options{})

	// Register workflow and activity
	w.RegisterWorkflow(SimpleWorkflow)
	w.RegisterActivity(SimpleActivity)

	// Start worker in background
	go func() {
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Worker failed: %v", err)
		}
	}()

	log.Println("✅ Worker started")

	// Execute workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("test-workflow-%d", time.Now().Unix()),
		TaskQueue: "test-task-queue",
		WorkflowExecutionTimeout: 1 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    1 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    3,
		},
	}

	log.Println("🏃 Executing workflow...")
	start := time.Now()

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, SimpleWorkflow, "Volcano LLM")
	if err != nil {
		log.Fatalf("❌ Failed to execute workflow: %v", err)
	}

	log.Printf("✅ Workflow started with ID: %s, RunID: %s", we.GetID(), we.GetRunID())

	// Wait for result
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalf("❌ Workflow failed: %v", err)
	}

	duration := time.Since(start)
	log.Printf("✅ Workflow completed in %v", duration)
	log.Printf("📋 Result: %s", result)

	// Show Temporal UI information
	log.Println("")
	log.Println("🌐 View workflow execution in Temporal UI:")
	log.Printf("   http://localhost:8088/namespaces/default/workflows/%s/%s", we.GetID(), we.GetRunID())
	log.Println("")

	// Show completion message
	log.Println("📊 Workflow execution details:")
	log.Printf("   - Workflow Type: SimpleWorkflow")
	log.Printf("   - Workflow ID: %s", we.GetID())
	log.Printf("   - Run ID: %s", we.GetRunID())
	log.Printf("   - Duration: %v", duration)

	log.Println("")
	log.Println("🎉 Temporal integration test completed successfully!")
}