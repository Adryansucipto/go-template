package config

type SchedulerConfig struct {
	Enabled  bool                   `yaml:"enabled"`
	Timezone string                 `yaml:"timezone"`
	Jobs     map[string]JobSchedule `yaml:"jobs"`
}

type JobSchedule struct {
	Enabled   bool   `yaml:"enabled"`
	Cron      string `yaml:"cron"`
	BatchSize int    `yaml:"batch_size,omitempty"`
}

