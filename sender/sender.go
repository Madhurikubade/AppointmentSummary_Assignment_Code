package sender

import (
	"AppointmentSummary_Assignment_Code/database"
	"AppointmentSummary_Assignment_Code/model"
	"AppointmentSummary_Assignment_Code/workerpool"
	"fmt"
)

// CreateAndScheduleSummaryAppointmentMessages processes and schedules messages concurrently
func CreateAndScheduleSummaryAppointmentMessages(date string) error {
	data, err := database.ReadDataForDate(date)
	if err != nil {
		return fmt.Errorf("failed to fetch appointment data: %w", err)
	}

	doctorMap := make(map[int][]model.AppointmentData)
	centerMap := make(map[int][]model.AppointmentData)

	for _, item := range data {
		appt := item.(model.AppointmentData)

		doctorMap[appt.DoctorID] = append(doctorMap[appt.DoctorID], appt)
		centerMap[appt.CenterID] = append(centerMap[appt.CenterID], appt)
	}

	for docID, appts := range doctorMap {
		msg := generateDoctorMessage(docID, appts, date)
		workerpool.WaitGroup.Add(1)
		workerpool.AddTask(func() {
			_ = database.InsertDoctorSummary(msg)
		})
	}

	for centerID, appts := range centerMap {
		msg := generateCenterMessage(centerID, appts, date)
		workerpool.WaitGroup.Add(1)
		workerpool.AddTask(func() {
			_ = database.InsertCenterSummary(msg)
		})
	}

	return nil
}

func generateDoctorMessage(doctorID int, appointments []model.AppointmentData, date string) model.DoctorMessage {
	var name string
	var mobile string
	categoryCount := make(map[string]int)

	for _, appt := range appointments {
		name = appt.DoctorName
		mobile = appt.DoctorMobile
		categoryCount[appt.Category]++
	}

	msg := fmt.Sprintf("Doctor %s had the following appointments on %s:\n", name, date)
	for cat, count := range categoryCount {
		msg += fmt.Sprintf("- %s: %d\n", cat, count)
	}

	return model.DoctorMessage{
		DoctorID:     doctorID,
		DoctorMobile: mobile,
		Date:         date,
		Message:      msg,
	}
}

func generateCenterMessage(centerID int, appointments []model.AppointmentData, date string) model.CenterMessage {
	var name string
	categoryCount := make(map[string]int)

	for _, appt := range appointments {
		name = appt.CenterName
		categoryCount[appt.Category]++
	}

	msg := fmt.Sprintf("Center %s had the following appointments on %s:\n", name, date)
	for cat, count := range categoryCount {
		msg += fmt.Sprintf("- %s: %d\n", cat, count)
	}

	return model.CenterMessage{
		CenterID: centerID,
		Date:     date,
		Message:  msg,
	}
}
