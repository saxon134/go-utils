package saBroker

var _chan chan *LocalJob
var _handle func(j *LocalJob)

func initLocal(concurrent int) {
	if concurrent <= 0 {
		concurrent = 10
	} else if concurrent > 100 {
		concurrent = 100
	}

	_chan = make(chan *LocalJob, concurrent)
	go func() {
		for {
			j := <-_chan
			_handle(j)
		}
	}()
}

func LocalDo(v *LocalJob) {
	_chan <- v
}
