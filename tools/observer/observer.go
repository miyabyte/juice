package observer

import "sync"

var initOnce sync.Once

type Oberver interface {
	Name() string
	Run(data interface{})
}

type SubjectObserverM struct {
	subjectObservers map[string]map[string]Oberver
}

func (som *SubjectObserverM) RegisterObserver(subjectName string, observer Oberver) (err error) {
	initOnce.Do(func() {
		som.subjectObservers = make(map[string]map[string]Oberver)
	})

	if _, ok := som.subjectObservers[subjectName]; !ok {
		som.subjectObservers[subjectName] = make(map[string]Oberver)
	}

	if _, ok := som.subjectObservers[subjectName][observer.Name()]; !ok {
		som.subjectObservers[subjectName][observer.Name()] = observer
	}
	return
}

func (som *SubjectObserverM) NotifySubjectObserver(subjectName string, data interface{}) (err error) {
	if obs, ok := som.subjectObservers[subjectName]; ok {
		for _, ob := range obs {
			ob.Run(data)
		}
	}
	return
}
