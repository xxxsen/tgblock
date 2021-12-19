package locker

type Locker interface {
	Lock(key string) bool
	Unlock(key string)
}
