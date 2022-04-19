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
