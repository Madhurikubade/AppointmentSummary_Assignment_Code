# Appointment Summary Messages
This assignment should give you a quick overview of the kinds of tasks handled on the backend.

You need to first set up your local database by inserting the data from the given .csv files. You will then summarize them two ways through Golang code, and write the summarized messages to new tables.

**The expected format of the summary messages is explained in sender/sender.go.**

You have been given some boilerplate code for the flow. Please read through it, starting from main.go, and clarify as needed.

### You can import any packages and change any function signatures as needed. Please focus on writing clean and performant code. You can also modify the boilerplate - the end result is what is important!
### Note that you might need to create indexes in the tables to optimize your queries. You are free to do so.

You can run the program as ./main date (or .\main.exe date on Windows).

#### Brief explanation of each table given to you in the .csv files:
- Center: Represents a specific clinic. Contains a unique id (CenterID) and the name of the center
- Patient: Represents a specific person. Contains a unique id (PatientID), the person's name info (Salutation and Name), and their mobile number
- DoctorStaff: Represents a doctor. Contains a unique id (DoctorStaffID), the doctor's name, and their mobile
- Appointment: Represents an appointment a patient has with a doctor at a specific center. Contains a unique identifier, the center, doctor, and patient IDs, the start and end time, the status (S for scheduled, C for cancelled), and the treatment category (Consultation, Not Specified, Oral Surgery, Ortho, Endo). **Note that cancelled appointments should not show up in any summary message, in either the doctor-wise messages or in the center-wise messages**