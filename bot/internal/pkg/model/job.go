package model

type JobType int

const (
	GetTemperature JobType = iota
)

type JobSettings struct {
	Type     JobType `json:"type"`
	TickRate int     `json:"tick_rate"`
}

type GetTemperatureJob struct {
	JobSettings
}

func NewGetTemperatureJob(tickRate int) GetTemperatureJob {
	return GetTemperatureJob{
		JobSettings: JobSettings{
			Type:     GetTemperature,
			TickRate: tickRate,
		},
	}
}

func (j GetTemperatureJob) GetSettings() JobSettings {
	return j.JobSettings
}
