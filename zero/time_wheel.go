package zero

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type SlotElement interface {
	GetID() string
	Closed() bool
}

type slot struct {
	id       int
	elements map[string]SlotElement
}

func newSlot(id int) *slot {
	s := &slot{
		id:       id,
		elements: make(map[string]SlotElement),
	}
	return s
}

func (s *slot) add(e SlotElement) {
	s.elements[e.GetID()] = e
}

func (s *slot) remove(e SlotElement) {
	delete(s.elements, e.GetID())
}

type Handler func(e SlotElement)
type TimeWheel struct {
	sync.RWMutex
	duration     time.Duration
	perWheel     int
	currentIndex int
	wheel        []*slot
	cache        map[string]*slot
	handler      Handler
	ticker       *time.Ticker
	task         chan SlotElement
	stop         chan struct{}
}

func NewTimeWheel(d time.Duration, wheel int, f Handler) *TimeWheel {
	if d < 1 || wheel < 1 || f == nil {
		return nil
	}
	t := &TimeWheel{
		duration:     d,
		perWheel:     wheel,
		handler:      f,
		currentIndex: 0,
		task:         make(chan SlotElement, 10000),
		stop:         make(chan struct{}),
	}
	t.cache = make(map[string]*slot)
	t.wheel = make([]*slot, wheel)
	for i := 0; i < wheel; i++ {
		t.wheel[i] = newSlot(i)
	}
	return t
}

func (t *TimeWheel) Start() {
	if t.ticker == nil {
		t.ticker = time.NewTicker(t.duration)
		go t.run()
	}
}

func (t *TimeWheel) Stop() {
	if t.ticker != nil {
		go func() {
			t.stop <- struct{}{}
		}()
	}
}
func (t *TimeWheel) Add(e SlotElement) {
	t.task <- e
}

func (t *TimeWheel) Remove(e SlotElement) {
	if v, ok := t.cache[e.GetID()]; ok {
		v.remove(e)
		logrus.Infof("slot: %d 移除连接: %s", v.id, e.GetID())
	}
}

func (t *TimeWheel) prevTickIndex() int {
	t.RLock()
	defer t.RUnlock()
	i := t.currentIndex
	if i == 0 {
		return t.perWheel - 1
	}
	return i - 1
}

func (t *TimeWheel) run() {
	for {
		select {
		case <-t.stop:
			t.ticker.Stop()
			t.ticker = nil
		case <-t.ticker.C:
			if t.currentIndex == t.perWheel {
				t.currentIndex = 0
			}
			slot := t.wheel[t.currentIndex]
			count := len(slot.elements)
			if count > 0 {
				logrus.Infof("slot: %d 有过期连接: %d", slot.id, count)
				for _, v := range slot.elements {
					slot.remove(v)
					if t.handler != nil {
						t.handler(v)
					}
				}
			}
			t.currentIndex++
		case v := <-t.task:
			t.Remove(v)
			if v.Closed() {
				continue
			}
			slot := t.wheel[t.prevTickIndex()]
			slot.add(v)
			logrus.Infof("slot: %d 添加连接: %s", slot.id, v.GetID())
			t.cache[v.GetID()] = slot
		}
	}
}
