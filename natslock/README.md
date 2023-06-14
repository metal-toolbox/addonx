# natslock

This package leverages NATS JetStream key-value stores to implement distributed locking/leader election. 

The initial use-case is for governor addons that run timed reconcile loops and need to avoid multiple executions when several instances of the same addon are started. Using this package, the first instance to obtain the lock becomes the leader and will run all of the timed reconcile loops going forward. When gracefully stopped, it releases the lock and one of the others becomes the new leader at the next run. If a leader becomes unresponsive, the lock automatically times out after a pre-determined interval, allowing for a new leader to be elected.

*Warning:*
The current implementation has the potential to introduce a race condition (because of the KeyValue `Get` and `Put`). If that's critical for your use case you probably shouldn't use this, or at least add some extra checks.

## Example usage

```go
    // assume there's an existing nc *nats.Conn

    // create a jetstream context
    jets, err := nc.JetStream()
    if err != nil {
        return nil, err
    }

    // initialize the KeyValue store with the given name and ttl
    kvStore, err := natslock.NewKeyValue(jets, "my-lock-bucket", time.Hour)
    if err != nil {
        return nil, err
    }

    // initialize the Locker
    locker := natslock.New(natslock.WithKeyValueStore(kvStore))

    // then inside your loop ...

    for {
        // check if this instance is the lead (or acquire the lock if no one else is)
        isLead, err := locker.AcquireLead()
        if err != nil {
            // error acquiring/checking leader lock (you may want to fail here)
            continue
        }

        if !isLead {
            continue
        }

        // execute remainder of code
    }
```
