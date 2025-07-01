package main

import (
	"AppointmentSummary_Assignment_Code/config"
	"AppointmentSummary_Assignment_Code/database"
	"AppointmentSummary_Assignment_Code/sender"
	"AppointmentSummary_Assignment_Code/workerpool"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide the date argument in format YYYY-MM-DD")
	}
	date := os.Args[1]

	database.ConnectMysqlDB()
	defer database.CloseMysql()
	// create worker pool
	log.Println("Creating and starting worker pool......")
	workerpool.NewWorkerPool(config.AppConfig.WorkerPoolConf.MaxWorker, config.AppConfig.WorkerPoolConf.BufferSize)
	workerpool.Run()
	log.Println("Worker pool created")
	err := sender.CreateAndScheduleSummaryAppointmentMessages(date)
	if err != nil {
		log.Fatalf("Failed to create summary messages: %v", err)
	}

	// Wait for all tasks to complete
	workerpool.WaitGroup.Wait()
	workerpool.CloseWorkerPool()
	log.Println("Summary messages created successfully for date:", date)
}
