package zero

import "github.com/0x1306a94/go-message-push/cmap"

type AssociatedManager struct {
	container cmap.ConcurrentMap
}

func NewAssociatedManager() *AssociatedManager {
	return &AssociatedManager{
		container: cmap.New(),
	}
}

func (a *AssociatedManager) Set(key string, conn *Conn) {
	a.container.Set(key, conn)
}

func (a *AssociatedManager) Del(key string) {
	a.container.Remove(key)
}
