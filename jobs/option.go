package jobs

import "time"

type Option func(j *Job)

func WithInterval(interval time.Duration) Option {
	return func(j *Job) {
		if interval == 0 {
			interval = time.Second
		}
		j.interval = interval
	}
}

func WithRepeat(repeat int) Option {
	return func(j *Job) {
		j.repeat = repeat
	}
}

func WithName(name string) Option {
	return func(j *Job) {
		j.name = name
	}
}

type ItemOption func(j *JobItem)

func WithItemId(id int64) ItemOption {
	return func(j *JobItem) {
		j.JobId = id
	}
}
func WithItemDebug(debug bool) ItemOption {
	return func(j *JobItem) {
		j.Debug = debug
	}
}
func WithItemReplace(replace bool) ItemOption {
	return func(j *JobItem) {
		j.Replace = replace
	}
}
