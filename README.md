# Barbora Notifier

![Docker Image CI](https://github.com/apsega/barbora-notifier/workflows/Docker%20Image%20CI/badge.svg?branch=master) ![Go](https://github.com/apsega/barbora-notifier/workflows/Go/badge.svg?branch=master)

App that checks if there's available slot in Barbora eshop.

## Features

* Checks [Barbora eshop](https://www.barbora.lt/) for slot availability;
* Prints and notifies to Slack when slot is available.

## Getting Started

Setup Slack incoming webhook - https://api.slack.com/messaging/webhooks and get unique URL, something like `https://hooks.slack.com/services/xxx/xxx/xxx`

Clone repo and build application Docker container:

```sh
docker build --tag barbora-notifier:1.1 .
```

Run Docker container:

```sh
docker run --rm barbora-notifier:1.1 --email=$EMAIL --password=$PASSWORD --webhook="$SLACK_WEBHOOK_URL"
```

i.e.:

```sh
$ docker run --rm barbora-notifier:1.1 --email=myemail@gmail.com --password=p455w0rD@l33T --webhook="https://hooks.slack.com/services/xxx/xxx/xxx"

balandžio 4 d. (šeštadienis)
Delivery time: 08 - 09
Available: false

Delivery time: 09 - 10
Available: false

Delivery time: 10 - 11
Available: false

Delivery time: 11 - 12
Available: false

Delivery time: 12 - 13
Available: false

Delivery time: 13 - 14
Available: false

Delivery time: 14 - 15
Available: false

Delivery time: 15 - 16
Available: false

Delivery time: 16 - 17
Available: false

Delivery time: 17 - 18
Available: false

Delivery time: 18 - 19
Available: false

Delivery time: 19 - 20
Available: false

Delivery time: 20 - 21
Available: false
```

