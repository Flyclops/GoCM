GoCM
====

Super-simple asynchronous GCM notification send service written in Go


Motivation
-----------
The world already has [pyapns](https://github.com/samuraisam/pyapns) for asynchronous sending of Apple Push Notifications. We needed similar async functionality for Google Cloud Messages that we use for push notifications on Android.

Since we're (slowly) moving most of our backend codebase to [Go](http://golang.org), any new code is written in Go.

[pyapns](https://github.com/samuraisam/pyapns) is so nice because [Apple Push Notification Services](https://developer.apple.com/library/ios/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/Chapters/ApplePushService.html) require a long-open socket through which all messages are sent. This minimizes the number of handshakes that need to take place, and gives the app an open pipe through which to dump messages.

GCM, however, is a RESTful service (as of the writing of this README, the "send" URL is https://android.googleapis.com/gcm/send), that requires one HTTP connection per message sent. It does support multicasting to up to 1000 devices at once, but the message must be identical to all devices.

In our testing, from API servers running on AWS EC2 c1.medium cloud servers, sending a single APNS push notification takes < 10ms on an open socket, while sending a single GCM notification can take between 150-250ms. Since we send hundreds of notifications per second, this became unacceptable.

The service
------------

In order to not reinvent the wheel, we started with [Alex Lockwood's](https://github.com/alexjlockwood) open source GCM package for Go, called simply [gcm](https://github.com/alexjlockwood/gcm). Surrounding that, we instantiate a standard web server that takes an ```-apikey``` argument (and also configurable ```-ipaddress``` and ```-port``` arguments if the defaults are not desired) and listens for incoming POST requests. 

The POST request should include two key-value pairs, ```token``` (the GCM device token) and ```payload``` (the JSON packet to send to GCM).

The server returns immediately, while pushing the main bulk of the work on to a new goroutine.

Example
--------

Start the server

```shell
./GoCM --apikey <GCM_API_KEY>
```

Send a message via ```curl```:
```shell
curl -d "token=<GCM_DEVICE_TOKEN>&payload={\"title\": \"This is the title\", \"subtitle\": \"This is the subtitle\", \"tickerText\": \"This is the ticker text\", \"datestamp\": \"2014-03-07T18:01:04.702100\"}" localhost:5601/gcm/send
```


To Do
-----------

- Output logging to a log file specified on the command-line
- Perhaps make runnable on a UNIX socket

