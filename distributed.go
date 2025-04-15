//go:generate mockgen -destination=mocks/distributed.go -package=gocronmocks . Elector,Locker,Lock
package gocron

import (
	"context"
	"time"
)

// Elector determines the leader from instances asking to be the leader. Only
// the leader runs jobs. If the leader goes down, a new leader will be elected.
type Elector interface {
	// IsLeader should return  nil if the job should be scheduled by the instance
	// making the request and an error if the job should not be scheduled.
	IsLeader(context.Context) error
}

// Locker represents the required interface to lock jobs when running multiple schedulers.
// The lock is held for the duration of the job's run, and it is expected that the
// locker implementation handles time splay between schedulers.
// The lock key passed is the job's name - which, if not set, defaults to the
// go function's name, e.g. "pkg.myJob" for func myJob() {} in pkg
type Locker interface {
	// Lock if an error is returned by lock, the job will not be scheduled.
	// If ttl is nil, the lock will be held for the default ttl of the underlying
	// distributed lock implementation.
	// For example, for the Redis distributed lock implementation from github.com/go-co-op/gocron-redis-lock, the ttl is controlled by redislock.WithExpiry, and defaults to 8s.
	// If ttl is not nil, the lock will be held for minimum between the duration of the job's run, the jobs timeout, and the ttl.
	Lock(ctx context.Context, key string, ttl *time.Duration) (Lock, error)
}

// Lock represents an obtained lock. The lock is released after the execution of the job
// by the scheduler.
type Lock interface {
	Unlock(ctx context.Context) error
}
