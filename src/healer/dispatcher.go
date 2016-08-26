package healer

import "sync"

type Dispatcher struct {
	// State: active or passive
	State    string
	taskCh 	 chan *EvacContainer
	resultCh chan *EvacContainer
}

func NewDispatcher(resultChannel chan<- *EvacContainer) *Dispatcher {
	return &Dispatcher{
		State: "passive",
		taskCh: make(chan *EvacContainer),
		resultCh: resultChannel,
	}
}

func (d *Dispatcher) dispatch (evac_Q *EvacContainer, accept_Q *EvacContainer) {
	for d.State == "active" {
		switch {
		case len(accept_Q) > 0:
			// TODO: with lock
			container := accept_Q[len(accept_Q)-1]
			accept_Q = accept_Q[:len(accept_Q)-1]
			container.SetTask("check evacuation")
			d.taskCh <- container
		case len(evac_Q) > 0:
			// TODO: with lock
			container := evac_Q[len(accept_Q)-1]
			evac_Q = evac_Q[:len(accept_Q)-1]
			container.SetTask("evacuate")
			d.taskCh <- container
		default:
			d.State = "passive"
		}
	}
}

func (d *Dispatcher) passivate() {
	d.State = "passive"
	// TODO: close pool
}

func (d *Dispatcher) activate() {
	d.State = "active"
}

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
// TODO: Evacuate
// TODO: Queery for a state of VM
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
