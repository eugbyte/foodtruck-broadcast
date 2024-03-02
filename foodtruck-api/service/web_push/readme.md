# Dev
1. Start the localstack container
```
make compose-down
make compose
```

2. Deploy the aws resources, e.g. lambdas, via localstack
```
make tflocal
```

3. Navigate to the endpoint, as shown in the printed output

Example:
```
curl --location --request GET 'http://9e1ad919.execute-api.localhost.localstack.cloud:4566/api/v1/subscription' \
--header 'Content-Type: application/json' \
--data '{
    "geohash": "abc",
    "endpoint": "https://mock.com",
    "lastsend": 1682143194727,
    "expiration": 1682143194727,
    "auth": "abc",
    "p256dh": "abc",
    "optIn": true
}'
```

Alternatively, run `make tflocal-refresh` which combines steps 1 to 2.

# Gotachs
1. The endpoint for localstack services, when invoked by the localstack container itself, is dynamic, and must be determined [during run time](https://github.com/localstack/localstack/issues/2511).
2. When zipping the lambda file, remember to set the [chmod to allow execution](https://github.com/hashicorp/terraform-provider-archive/issues/10#issuecomment-1537081365), otherwise aws cannot execute the file.
