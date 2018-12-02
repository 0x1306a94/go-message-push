package zero

import (
	"fmt"
	"testing"
	"time"
)

type mockConn struct {
	id int
}

func (m *mockConn) GetID() string {
	return fmt.Sprintf("id:%d", m.id)
}

func TestTimeWheel(t *testing.T) {

	timeWheel := NewTimeWheel(time.Second*1, 30, func(e SlotElement) {
		fmt.Println("心跳超时: ", e.GetID())
	})
	timeWheel.Start()
	for i := 0; i < 5; i++ {
		id := i + 1
		m := &mockConn{id: id}
		timeWheel.Add(m)
	}
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for {
			<-ticker.C
			fmt.Println("更新心跳: id:5")
			m := &mockConn{id: 5}
			timeWheel.Add(m)
		}
	}()

	select {}
}
