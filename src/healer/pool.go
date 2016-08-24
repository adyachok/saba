package healer

import "sync"

type Worker interface {
	Task()
}

type Pool struct {
	wg    sync.WaitGroup
	tasks chan Worker
}

// TODO: Add logic of two chanals
// TODO: receive VMs for evacuation
// TODO: send info about evacuation
func NewPool(numProcesses int) *Pool {
	p := &Pool{
		tasks: make(chan Worker),
	}

	p.wg.Add(numProcesses)

	for i:=0; i < numProcesses; i++ {
		go func(){
			for w := range p.tasks {
				w.Task()
			}
			p.wg.Done()
		}()
	}
	return p
}

func (p *Pool) Run(task Worker) {
	p.tasks <- task
}

func (p *Pool) Shutdown() {
	close(p.tasks)
	p.wg.Wait()
}
