Dashboard UI at http://localhost:9000 (username: admin, password: admin)

APISIX gateway at http://localhost:9080

Admin API at http://localhost:9180

Setup api gateway
```
docker exec -it --user root sad-apisix-1 /bin/sh -c "./setup.sh"
```

```
curl -i "http://127.0.0.1:9180/apisix/admin/routes" -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1" -X PUT -d '   
{
  "id": "getting-started-ip",
  "uri": "/ip",
  "upstream": {
    "type": "roundrobin",
    "nodes": {
      "httpbin.org:80": 1
    }
  }
}'
```

```
curl "http://127.0.0.1:9080/ip"
```

List existing upstreams
```
curl -X GET "http://127.0.0.1:9180/apisix/admin/upstreams" \
  -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1"
```

List existing routes
```
curl -X GET "http://127.0.0.1:9180/apisix/admin/routes" \
  -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1"
```

