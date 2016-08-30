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

func NewDispatcher(client *gophercloud.ServiceClient, resultChannel chan<- *EvacContainer) *Dispatcher {
	return &Dispatcher{
		State:    "passive",
		pool:	  NewPool(client, 4, resultChannel),
	}
}

func (d *Dispatcher) dispatch (qm *QueueManager) {
	var container *EvacContainer
	for d.State == "active" {
		switch {
		case len(qm.Accepted_Q) > 0:
			qm.lock.RLock()
			container = qm.Accepted_Q[len(qm.Accepted_Q)-1]
			qm.Accepted_Q = qm.Accepted_Q[:len(qm.Accepted_Q)-1]
			qm.lock.RUnlock()
			d.pool.Run(container)
		case len(qm.Scheduled_Q) > 0:
			qm.lock.RLock()
			container = qm.Scheduled_Q[len(qm.Scheduled_Q)-1]
			qm.Scheduled_Q = qm.Scheduled_Q[:len(qm.Scheduled_Q)-1]
			qm.lock.RUnlock()
			d.pool.Run(container)
		default:
			d.passivate()
		}
	}
}

func (d *Dispatcher) passivate() {
	d.State = "passive"
}

func (d *Dispatcher) activate(qm *QueueManager) {
	if d.State != "active" {
		d.State = "active"
		d.dispatch(qm)
	}
}

func (d *Dispatcher) shutdown() {
	d.State = "shutdown"
	d.pool.Shutdown()
	return
}


type Pool struct {
	wg        sync.WaitGroup
	tasksCh   chan *EvacContainer
	resultsCh chan<- *EvacContainer
}

func NewPool(client *gophercloud.ServiceClient, numProcesses int, results chan<- *EvacContainer) *Pool {
	p := &Pool{
		tasksCh:   make(chan *EvacContainer),
		resultsCh: results,
	}

	p.wg.Add(numProcesses)

	for i:=0; i < numProcesses; i++ {
		go func(){
			for w := range p.tasksCh {
				w.Task(client)
				results <- w

			}
			p.wg.Done()
		}()
	}
	return p
}

func (p *Pool) Run(task *EvacContainer) {
	p.tasksCh <- task
}

func (p *Pool) Shutdown() {
	close(p.tasksCh)
	p.wg.Wait()
}
