// Copyright 2018 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package pool

import (
	"errors"
	"time"
)

// ErrScheduleTimeout returned by Pool to indicate that there no free
// goroutines during some period of time.
var ErrScheduleTimeout = errors.New("schedule error: timed out")

// Pool contains logic of goroutine reuse.
type Pool struct {
	sem  chan struct{}
	work chan func()
}

// NewPool creates new goroutine pool with given size. It also creates a work
// queue of given size. Finally, it spawns given amount of goroutines
// immediately.
func NewPool(size, queue, spawn int) (*Pool, error) {
	if spawn <= 0 && queue > 0 {
		return nil, errors.New("dead queue configuration detected")
	}

	if spawn > size {
		return nil, errors.New("spawn > workers")
	}

	p := &Pool{
		sem:  make(chan struct{}, size),
		work: make(chan func(), queue),
	}

	for i := 0; i < spawn; i++ {
		p.sem <- struct{}{}

		go p.worker(func() {})
	}

	return p, nil
}

// Schedule schedules task to be executed over pool's workers.
func (p *Pool) Schedule(task func()) {
	p.schedule(task, nil)
}

// ScheduleTimeout schedules task to be executed over pool's workers.
// It returns ErrScheduleTimeout when no free workers met during given timeout.
func (p *Pool) ScheduleTimeout(timeout time.Duration, task func()) error {
	return p.schedule(task, time.After(timeout))
}

func (p *Pool) schedule(task func(), timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return ErrScheduleTimeout
	case p.work <- task:
		return nil
	case p.sem <- struct{}{}:
		go p.worker(task)

		return nil
	}
}

func (p *Pool) worker(task func()) {
	defer func() {
		<-p.sem
	}()

	task()

	for task := range p.work {
		task()
	}
}
