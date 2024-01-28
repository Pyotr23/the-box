package model

import "context"

type JobType int

const (
	GetTemperature JobType = iota
)

type JobSettings struct {
	Type     JobType `json:"type"`
	TickRate int     `json:"tick_rate"` // Count of seconds between ticks
}

type temperatureGetter interface {
	GetTemperature(ctx context.Context, id int) (string, error)
}

type GetTemperatureJob struct {
	JobSettings
	DeviceIDs []int `json:"device_ids"`
	getter    temperatureGetter
}

func NewGetTemperatureJob(tickRate int, deviceIDs []int, getter temperatureGetter) GetTemperatureJob {
	return GetTemperatureJob{
		JobSettings: JobSettings{
			Type:     GetTemperature,
			TickRate: tickRate,
		},
		DeviceIDs: deviceIDs,
		getter:    getter,
	}
}

func (j GetTemperatureJob) GetSettings() JobSettings {
	return j.JobSettings
}

func (j GetTemperatureJob) Do() {

}
