# go-sync
Higher level synchronization primitives in Go

Ideas are mostly inspired by [The Little Book of
Semaphores](http://greenteapress.com/semaphores/LittleBookOfSemaphores.pdf)

All types except the semaphore itself are implemented on top of the semaphore
type defined in the sync/semaphore package. Any semaphore implementation can be
used as a drop-in replacement for this. Our current implementation makes use of
go standard sync library's Cond and Mutex type, but we could do away with them
if needed. Possible alternative is to use System V semaphores.
