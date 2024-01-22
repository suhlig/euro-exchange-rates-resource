# Euro Exchange Rates Resource

This is an example resource for the [concourse-resource-go](https://github.com/suhlig/concourse-resource-go) interface. It fetches currency exchange rates from the European Central Bank via [Frankfurter](https://github.com/hakanensari/frankfurter).

# Development

## Check

Native:

```command
$ jo -d . source.url=https://api.frankfurter.app | go run . check
```

Docker:

```command
$ jo -d . source.url=https://api.frankfurter.app | docker run --rm -i euro-exchange-rates-resource /opt/resource/check
```

## Get

Native:

```command
$ jo -d . source.url=https://api.frankfurter.app 'params.currencies[]=USD' version.date=2024-01-15 | go run . get /tmp
```

Get what check discovered:

```command
$ jo -d . source.url=https://api.frankfurter.app 'params.currencies[]=USD' version=$(
  jo -d . source.url=https://api.frankfurter.app | go run . check
) \
  | jq '.version=.version[0]' \
  | go run . get $(mktemp -d)
```

Docker:

```command
$ jo -d . source.url=https://api.frankfurter.app 'params.currencies[]=SEK' 'params.currencies[]=USD' version=$(
  jo -d . source.url=https://api.frankfurter.app | docker run --rm -i euro-exchange-rates-resource /opt/resource/check
) \
  | jq '.version=.version[0]' \
  | docker run --rm -i euro-exchange-rates-resource /opt/resource/in /tmp
```

# TODO

* `float32` is not ideal for money. Consider [shopspring/decimal](https://github.com/shopspring/decimal) or store everything in microcents.
