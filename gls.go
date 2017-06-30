package gls

import "github.com/huandu/goroutine"
import "sync"

func Set(key string, value interface{}) {
	goid := Goid()
	Gls.GetGlsData(goid).Set(goid, key, value)
}

func Get(key string) interface{} {
	goid := Goid()
	return Gls.GetGlsData(goid).Get(goid, key)
}

func Cleanup() {
	goid := Goid()
	Gls.GetGlsData(goid).Cleanup(goid)
}

func Goid() int64 {
	return goroutine.GoroutineId()
}

// Copy gls in current goroutine to child goroutine.
func Go(f func()) {
	parent := Goid()
	go func(){
		values := Values{}
		Gls.GetGlsData(parent).GetValues(parent, values)
		goid := Goid()
		Gls.GetGlsData(goid).SetValues(goid, values)
		defer Cleanup()
		f()
	}()
}

type GlsGroup struct {
	Group [1024]*GlsData
}

func newGlsGroup() *GlsGroup {
	obj := &GlsGroup{}
	for i := 0; i < len(obj.Group); i++ {
		obj.Group[i] = newGlsData()
    }
	return obj
}

func (this *GlsGroup) GetGlsData(goid int64) *GlsData {
	return this.Group[goid & 1023]
}

var Gls = newGlsGroup()

type Values map[interface{}]interface{}

type GlsData struct {
	lock sync.RWMutex
	data map[int64]Values
}

func newGlsData() *GlsData {
	return &GlsData {
		data : make(map[int64]Values),
    }
}

func (this *GlsData) Set(goid int64, key string, value interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()

	values, exists := this.data[goid]
	if !exists {
		values = Values{}
		this.data[goid] = values
    }

	values[key] = value
}

func (this *GlsData) Get(goid int64, key string) interface{} {
	this.lock.RLock()
	defer this.lock.RUnlock()

	values, exists := this.data[goid]
	if !exists {
		return nil
	}

	return values[key]
}

func (this *GlsData) Cleanup(goid int64) {
	this.lock.Lock()
	defer this.lock.Unlock()

	delete(this.data, goid)
}

func (this *GlsData) GetValues(goid int64, out Values) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	values, exists := this.data[goid]
	if !exists {
		return
	}

	for k, v := range values {
		out[k] = v
    }
}

func (this *GlsData) SetValues(goid int64, values Values) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.data[goid] = values
}

