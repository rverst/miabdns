# miabdns

miabdns is a simple server to implement a dynamic dns service for 
[Mail-in-a-Box](https://mailinabox.email/), written in [Go](https://golang.org/).

[![Go Report Card](https://goreportcard.com/badge/github.com/rverst/miabdns)](https://goreportcard.com/report/github.com/rverst/miabdns)


## Overview


## Installation

### Binary Release

* Download a binary release from the [release page](https://github.com/rverst/miabdns/releases).

* ~~Run the docker container from [docker hub](https://hub.docker.com/)~~ 

## Usage (to be continued...)

You need a configuration file, which provides credentials for your Mail-in-a-Box
instance(s) and a list of users that are allowed to use the service:

```json
    {"boxes": [
        {
            "name": "test",
            "server": "https://box.example.org",
            "username": "testUser",
            "password": "SuperSecret"
        }
    ],
    "users": [
        {
            "box": "test",
            "username": "user1",
            "password": "pwd1",
            "allow_ip4": true,
            "allow_ip6": false,
            "allowed_domains": [
                "www\\.example\\.org",
                ".*\\.?user1\\.example\\.org"
            ]
        }
    ]}
```
  
**Run `miabdns` for available commands.**


## Dependencies

miabdns uses and relies on the following, awesome libraries (in lexical order):

| Dependency | License |
| :------------- | :------------- |
| [github.com/integrii/flaggy](https://github.com/integrii/flaggy) | [The Unlicense](https://github.com/integrii/flaggy/blob/master/LICENSE) |
| [github.com/rverst/go-miab](https://github.com/rverst/go-miab) | [MIT License](https://github.com/rverst/go-miab/blob/master/LICENSE.txt) |

