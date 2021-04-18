# data-channel-create
data-channel-create is based on the Pion WebRTC application `data-channels-create` that shows how you can send/recv DataChannel messages. You need to run 2 instance at the same time for this application to function properly.

* This example is out-of-date. Please refer to data-channel-instance instead.

## Instructions
### Download data-channel-create
```
export GO111MODULE=on
go get github.com/Gaukas/seed2sdp/examples/data-channel-create
```

### Run data-channel-create

Run command below on `host1` where `[seed]` should be a shared secret between `host1` and `host2`.

```
data-channel-create [seed] offer
```

Run command below on `host2` where `[seed]` should be the same shared secret between `host1` and `host2`.

```
data-channel-create [seed] answer
```

### Signalling

Copy the line of text starting with `1,` emitted by `host1` and paste into `host2`

Copy the line of text starting with `2,` emitted by `host2` and paste into `host1`

If you see the text below in STDOUT of both instance:

```
Data channel 'data'-'824633868734' open. Random messages will now be sent to any connected DataChannels every 5 seconds
```

Congrats, data channel is established.=