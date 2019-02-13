# OAuth2 demo server

for Resource Owner Password Credentials Grant

```
$ curl "http://localhost:9096/token?grant_type=password&client_id=APP01&client_secret=APPSEC&username=admin&password=123456" | jq
{
  "access_token": "DVWMV9ZQPTOE-BQUHZHCVA",
  "expires_in": 120,
  "token_type": "Bearer"
}
```

```
$ curl "http://localhost:9096/test?access_token=DVWMV9ZQPTOE-BQUHZHCVA" | jq
{
  "client_id": "APP01",
  "expires_in": 97,
  "user_id": "000000"
}
```
