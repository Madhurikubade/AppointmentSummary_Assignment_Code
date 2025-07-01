package model

type Configurations struct {
	Prefix string
	Server struct {
		Host string
		Port string
	}

	MysqlConf      MysqlConfig
	WorkerPoolConf WorkerPool
}

type AppointmentData struct {
	CenterID     int
	CenterName   string
	DoctorID     int
	DoctorName   string
	DoctorMobile string
	PatientID    int
	PatientName  string
	Category     string
	TreatmentCat string
	StartTime    string
	EndTime      string
}

type MysqlConfig struct {
	Username             string
	Password             string
	Net                  string
	Address              string
	DatabaseName         string
	AllowNativePasswords bool
	ParseTime            bool
	MaxConnections       int
	MaxIdleTimeMinutes   int
}

type DoctorMessage struct {
	DoctorID     int
	DoctorMobile string
	Message      string
	Date         string
}

type CenterMessage struct {
	CenterID int
	Message  string
	Date     string
}

type WorkerPool struct {
	MaxWorker   int `json:"MaxWorker"`
	QueuedTaskC chan func()
	BufferSize  int `json:"BufferSize"`
}
