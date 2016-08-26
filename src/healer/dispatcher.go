package healer

import (
	"sync"
	"github.com/rackspace/gophercloud"
)

type Dispatcher struct {
	// State: active or passive
	State    string
	pool 	 *Pool
}

func NewDispatcher(resultChannel chan<- *EvacContainer) *Dispatcher {
	return &Dispatcher{
		State:    "passive",
		pool:	  NewPool(4, resultChannel),
	}
}

func (d *Dispatcher) dispatch (scheduled_Q *EvacContainer, accepted_Q *EvacContainer) {
	var container *EvacContainer
	for d.State == "active" {
		switch {
		case len(accepted_Q) > 0:
			container = accepted_Q[len(accepted_Q)-1]
			accepted_Q = accepted_Q[:len(accepted_Q)-1]
			container.SetTask("check evacuation")
			d.pool.Run(container)
		case len(scheduled_Q) > 0:
			container = scheduled_Q[len(scheduled_Q)-1]
			scheduled_Q = scheduled_Q[:len(scheduled_Q)-1]
			container.SetTask("evacuate")
			d.pool.Run(container)
		default:
			d.passivate()
		}
	}
}

func (d *Dispatcher) passivate() {
	d.State = "passive"
}

func (d *Dispatcher) activate(evac_Q *EvacContainer, accept_Q *EvacContainer) {
	if d.State != "active" {
		d.State = "active"
		d.dispatch(evac_Q, accept_Q)
	}
}

func (d *Dispatcher) shutdown() {
	d.State = "shutdown"
	d.pool.Shutdown()
}

type Worker interface {
	Task(client *gophercloud.ServiceClient)
}

type Pool struct {
	wg        sync.WaitGroup
	tasksCh   chan Worker
	resultsCh chan Worker
}

func NewPool(numProcesses int, results chan<- *Worker) *Pool {
	p := &Pool{
		tasksCh:   make(chan Worker),
		resultsCh: results,
	}

	p.wg.Add(numProcesses)

	client := gophercloud.ServiceClient{}

	for i:=0; i < numProcesses; i++ {
		go func(){
			for w := range p.tasksCh {
				w.Task(client)
				results <- &w

			}
			p.wg.Done()
		}()
	}
	return p
}

func (p *Pool) Run(task Worker) {
	p.tasksCh <- task
}

func (p *Pool) Shutdown() {
	close(p.tasksCh)
	p.wg.Wait()
}
