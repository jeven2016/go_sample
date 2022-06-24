SSO集成

### grant_type

```shell
grant_types_supported: [
"authorization_code",
"implicit",
"refresh_token",
"password",
"client_credentials"
],
```

* 申请授权码

```
```

* 使用授权码申请token

```shell
curl -X POST \
  http://localhost:9999/auth/realms/zhongfu/protocol/openid-connect/token \
  -H 'Cache-Control: no-cache' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'cache-control: no-cache' \
  -d 'grant_type=password&username=login&password=pwd&client_id=my-app&response_type=code'


curl \
  -d "client_id=web1" \
#  -d "client_secret=d5141c8b-ea12-4320-a500-74a046083c08" \
  -d "grant_type=client_credentials" \
  "http://localhost:9080/auth/realms/myrealm/protocol/openid-connect/token"
```

```shell

```

* Token格式

```shell
{
   "access_token":"",
   "expires_in":300,
   "refresh_expires_in":1800,
   "refresh_token":"",
   "token_type":"Bearer",
   "id_token":"",
   "not-before-policy":0,
   "session_state":"1ae3d862-4057-4af7-b6d4-19c1c3309018",
   "scope":"openid profile email"
}
```

