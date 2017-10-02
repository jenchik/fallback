# fallback
The pattern, similar to RWMutex. With low latency between RLock and WLock, as well as calling a slow lazy asynchronous data loader. The code looks like a WLock call inside RLock. Please, readme https://github.com/golang/go/issues/4026
