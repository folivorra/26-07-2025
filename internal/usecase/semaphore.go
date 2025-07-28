package usecase

type Semaphore struct {
	C chan struct{}
}

func NewSemaphore(max int) *Semaphore {
	return &Semaphore{
		C: make(chan struct{}, max),
	}
}

func (s *Semaphore) Acquire() {
	s.C <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.C
}
