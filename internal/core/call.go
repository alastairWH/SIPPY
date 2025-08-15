package core

import "sync"

type Call struct {
	Caller   string
	Callee   string
	Active   bool
}

type CallManager struct {
	calls map[string]*Call
	mu    sync.RWMutex
}

func NewCallManager() *CallManager {
	return &CallManager{calls: make(map[string]*Call)}
}

func (cm *CallManager) StartCall(caller, callee string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.calls[caller+":"+callee] = &Call{Caller: caller, Callee: callee, Active: true}
}

func (cm *CallManager) EndCall(caller, callee string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.calls, caller+":"+callee)
}

func (cm *CallManager) IsActive(caller, callee string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	_, ok := cm.calls[caller+":"+callee]
	return ok
}

func (cm *CallManager) Calls() map[string]*Call {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	copy := make(map[string]*Call)
	for k, v := range cm.calls {
		copy[k] = v
	}
	return copy
}
