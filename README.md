[![Xybor founder](https://img.shields.io/badge/xybor-huykingsofm-red)](https://github.com/huykingsofm)
[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/xysched.svg)](https://pkg.go.dev/github.com/xybor-x/xysched)
[![GitHub Repo stars](https://img.shields.io/github/stars/xybor-x/xysched?color=yellow)](https://github.com/xybor-x/xysched)
[![GitHub top language](https://img.shields.io/github/languages/top/xybor-x/xysched?color=lightblue)](https://go.dev/)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/xybor-x/xysched)](https://go.dev/blog/go1.18)
[![GitHub release (release name instead of tag name)](https://img.shields.io/github/v/release/xybor-x/xysched?include_prereleases)](https://github.com/xybor-x/xysched/releases/latest)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/4b74139b892144a4a2a1a2605fc15738)](https://www.codacy.com/gh/xybor-x/xysched/dashboard?utm_source=github.com&utm_medium=referral&utm_content=xybor-x/xysched&utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/4b74139b892144a4a2a1a2605fc15738)](https://www.codacy.com/gh/xybor-x/xysched/dashboard?utm_source=github.com&utm_medium=referral&utm_content=xybor-x/xysched&utm_campaign=Badge_Coverage)
[![Go Report](https://goreportcard.com/badge/github.com/xybor-x/xysched)](https://goreportcard.com/report/github.com/xybor-x/xysched)

# Introduction

Package xysched supports to schedule future tasks with a simple syntax.

# Features

There are two most important objects in this library:

-   `Future` defines tasks.
-   `Scheduler` manages to schedule `Future` instances.

`Scheduler` uses a channel to know when and which `Future` should be run. A
`Future` could be sent to this channel via `After` method and its variants:

```golang
func (s *Scheduler) After(d *time.Duration) chan<-Future
func (s *Scheduler) Now() chan<-Future
func (s *Scheduler) At(t *time.Time) chan<-Future
```

When a `Future` is sent via `After`, it will be called by `Scheduler` after a
duration `d`. This method is non-blocking.

There are some types of `Future`. For example, `Task` and `Cron`.

`Task` is a `Future` running only one time, wheares `Cron` could run
periodically.

For development, `Task` should be the base struct of all `Future` structs.
`Task` supports to add callback `Futures`, which is called after `Task`
completed. It also helps to handle the returned or panicked value, which is
equivalent to javascript `Promise`.

`Cron` also bases on `Task`, so `Cron` has all methods of `Task`.

# Example

1.  Print a message after one second.

```golang
xysched.After(time.Second) <- xysched.NewTask(fmt.Println, "this is a message")
```

2.  Increase x, then print a message.

```golang
var x int = 0

var future = xysched.NewTask(func() { x++ })
future.Callback(fmt.Println, "increase x")

xysched.Now() <- future
```

3.  Print a message every second.

```golang
xysched.Now() <- xysched.NewCron(fmt.Println, "welcome").Secondly()
```

4.  Increase x, then print a message. Loop over seven times. After all, print x.

```golang
var x int = 0
var future = xyshed.NewCron(func(){ x++ }).Secondly().Times(7)
future.Callback(fmt.Println, "increase x")
future.Finish(fmt.Printf, "the final value of x: %d\n", x)

xysched.Now() <- future
```

5.  It is also possible to use `Then` and `Catch` methods to handle the returned
    value of `Future` or recover if it panicked.

```golang
func foo(b bool) string {
    if b {
        panic("foo panicked")
    } else {
        return "foo bar"
    }
}

var future = xysched.NewTask(foo, true)
future.Then(func(s string) { fmt.Println(s) })
future.Catch(func(e error) { fmt.Println(e) })

xysched.Now() <- future
```

6.  Create a new scheduler if it is necessary. Scheduler with non-empty name can
    be used in many places without a global variable.

```golang
// a.go
var scheduler = xysched.NewScheduler("foo")
defer sched.Stop()
scheduler.After(3 * time.Second) <- xysched.NewTask(fmt.Println, "x")

// b.go
var scheduler = xysched.NewScheduler("foo")

// A scheduler should be stopped if it won't be used anymore.
scheduler.Stop()
```

7.  Early stop a future.

```golang
var sched = xysched.NewScheduler("")
defer sched.Stop()

var captured string
var task = xysched.NewTask(func() { captured = "run" })
sched.After(time.Millisecond) <- task
task.Stop()

time.Sleep(2 * time.Millisecond)
xycond.AssertEmpty(captured)
```
