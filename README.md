# seed2sdp
[![Build Status](https://travis-ci.com/Gaukas/seed2sdp.svg?branch=master)](https://travis-ci.com/Gaukas/seed2sdp)

---

Generate full-length SDP offers/answers from a shared-secret seed with minimal signalling process.

### Introduction

WebRTC and presumably many other amazing p2p communication implementations relies on SDP(Session Description Protocol). A traditional SDP weighs over 200 bytes and sometimes could be as long as 500 bytes while most of these data are for only integrity and confidentiality purposes and are unnecessary.

seed2sdp(this project) eliminates the random info in SDP by replacing the uncontrollable randomness with deterministic "randomness" based on HKDF readers and excluding all derivable information from the deflated SDP.

### Credits

[Pion](https://pion.ly/)
- Pure Go WebRTC API Implementation [pion/webrtc](https://github.com/pion/webrtc)
- Pure Go ICE Implementation [pion/ice](https:  //github.com/pion/ice)