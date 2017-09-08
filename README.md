# Description [![Build Status]

ratelimiter is a simple golang implementation of a thread safe, basic rate limiter.

# Usage

### Creating a Rate Limiter

To create a rate limiter, simply:

```
limiter := NewRateLimiter();
limiter.SetRate(5000.0); // 5000 Permit per second
```


### Using The Rate Limiter

There are two ways to aquire permits:

  * Blocking
  * Blocking with Timeout

For the basic blocking:

```
limiter->aquire(1); // aquires 1 permit. Will block the thread until it is allowed to proceed
```

For the timeout blocking:

```
limiter->try_aquire(2, 3000); // aquires 2 permit, and max timeout is 3000ms
```

On both cases, if try_aquire can aquire it's permits within the specified time (from now), it will block as long as necessary by the rate limiter, and then return true when it has aquired the permits. If it cannot aquire those permits within the specified time, then it will return IMMEDIATELY, with a value of false.
