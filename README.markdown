# Euro Exchange Rates Resource

This is an example resource for the [concourse-resource-go](https://github.com/homeport/concourse-resource-go) interface.

# Development

## Check

Native:

```command
$ jo -d . source.url=https://api.frankfurter.app/latest | go run . check
```

Docker:

```command
jo -d . source.url=https://api.frankfurter.app/latest | docker run --rm -i euro-exchange-rates-resource /opt/resource/check
```

## Get

Native:

```command
$ jo -d . source.url=https://api.frankfurter.app/latest 'params.currencies[]=EUR' version.date=2024-01-15 | go run . get /tmp
```

Get what check discovered:

```command
$ jo -d . source.url=https://api.frankfurter.app/latest 'params.currencies[]=EUR' version=$(
  jo -d . source.url=https://api.frankfurter.app/latest | go run . check
) \
  | jq '.version=.version[0]' \
  | go run . get $(mktemp -d)
```

Docker:

```command
$ jo -d . source.url=https://api.frankfurter.app/latest 'params.currencies[]=EUR' version=$(
  jo -d . source.url=https://api.frankfurter.app/latest | docker run --rm -i euro-exchange-rates-resource /opt/resource/check
) \
  | jq '.version=.version[0]' \
  | docker run --rm -i euro-exchange-rates-resource /opt/resource/in /tmp
```
