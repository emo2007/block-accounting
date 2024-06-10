package system

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type SystemQueue struct {
	m          sync.Mutex
	_init_size int
	buf        []*any
	_buf_p     atomic.Int64
	_read_p    atomic.Int64
}

func NewSystemQueue(size int) *SystemQueue {
	if size < 50 {
		size = 50
	}

	return &SystemQueue{
		_init_size: size,
		buf:        make([]*any, size),
	}
}

func (s *SystemQueue) Put(_ context.Context, job any) error {
	s.m.Lock()
	defer s.m.Unlock()

	p := s._buf_p.Load()

	// Resize buf if needed
	if int64(len(s.buf)) == p {
		s.resize()
	}

	s.buf[p] = &job

	s._buf_p.Add(1)

	return nil
}

func (s *SystemQueue) Pop(_ context.Context) (any, error) {
	s.m.Lock()
	defer s.m.Unlock()

	p := s._buf_p.Load()
	if int64(len(s.buf)) < p {
		s.resize()

		return nil, fmt.Errorf("system-queue: error _p index out of buffer range")
	}

	jobp := s.buf[p]

	if jobp == nil {
		return nil, fmt.Errorf("system-queue: error nil job")
	}

	job := *jobp
	s._buf_p.Add(-1)
	s.buf[p-1] = nil

	return job, nil
}

func (s *SystemQueue) resize() {
	s.buf = append(s.buf, make([]*any, s._init_size/2)...)
}
