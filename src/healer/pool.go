package healer
//
//import "sync"
//
//type Worker interface {
//	Task()
//}
//
//type Pool struct {
//	wg    sync.WaitGroup
//	tasks chan Worker
//}
//
//func NewPool(numProcesses int) *Pool {
//	p := &Pool{
//		wg: sync.WaitGroup{numProcesses},
//		tasks: make(chan Worker),
//	}
//
//	for i:=0; i < numProcesses; i++ {
//		go func(){
//			for w := range p.tasks {
//				w.Task()
//			}
//			p.wg.Done()
//		}()
//	}
//	return &p
//}
//
//func (p *Pool) Run(task Worker) {
//	p.tasks <- task
//}
//
//func (p *Pool) Shutdown() {
//	close(p.tasks)
//	p.wg.Wait()
//}
