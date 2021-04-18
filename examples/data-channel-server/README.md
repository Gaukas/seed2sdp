# data-channel-server
data-channel-create is based on the Pion WebRTC application `data-channels-server` that shows how you can send/recv DataChannel messages. You also need to run `data-channels-client` at the same time for this application to function properly.

## Instructions
### Download data-channel-server
```
export GO111MODULE=on
go get github.com/Gaukas/seed2sdp/examples/data-channel-client
```

### Run data-channel-server

Run command below on `client` where `seed` should be a shared secret between `client` and `server`.

```
data-channel-client serverIP serverPort seed
```

Run command below on `server` where `seed` should be the same shared secret between `client` and `server`.

```
data-channel-server bindtoIP bindtoPort seed
```

* IP must be 1:1 DNAT or bound to an interface

### Signalling

You will see output on `client` similar to below

```
==================== Def Offer ====================
1,9518543359329632256,0,3351379968
===================================================
```

Copy the center line (`1,9518543359329632256,0,3351379968` in this example) into the `server` after `data-channel-server` started. This must be done in 10 seconds after `data-channel-client` starts.

Then you should see on `server`

```
====================== Answer ======================
{"type":"answer","sdp":"v=0\r\no=- 1904260364563689595 1618785297 IN IP4 0.0.0.0\r\ns=-\r\nt=0 0\r\na=fingerprint:sha-256 D4:62:11:9B:02:C2:88:47:A6:D2:93:4F:AA:85:7A:B3:D0:8A:93:B5:90:46:AD:DC:DE:15:A5:08:F5:4E:E1:2C\r\na=group:BUNDLE 0\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=setup:active\r\na=mid:0\r\na=sendrecv\r\na=sctp-port:5000\r\na=ice-ufrag:xKVqtBCWxwrFLgYR\r\na=ice-pwd:eDKLjLyTufRqQWxssboKCZYQnJDcixpe\r\na=candidate:3031764619 1 udp 2130706431 192.168.24.132 27031 typ host\r\na=candidate:3031764619 2 udp 2130706431 192.168.24.132 27031 typ host\r\na=end-of-candidates\r\n"}
====================================================
==================== Def Answer ====================
2,9518543359329632256,0,1771503616
====================================================
```

However, you don't need to do anything. These are debugging info. 

In 10 seconds, you shall see

```
[Info] Peer Connection changed state to: connected
[Info] Successfully opened Data Channel 'Client DataChannel'-'824638048450'. 
```

on both side and there will be a count-up ping-pong happening. After they count to 1024 on both side, the counting will be stopped and datachannel will be destroyed in 10 seconds.

Congrats! You just used webrtc datachannel with single direction signal. 
