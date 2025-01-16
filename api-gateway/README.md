Dashboard UI at http://localhost:9000 (username: admin, password: admin)

APISIX gateway at http://localhost:9080 and http://localhost

Admin API at http://localhost:9180

Setup API Gateway routes. You might need to change container name in the script.
```
docker exec -it --user root sad-apisix-1 /bin/sh -c "./setup.sh"
```

Validate/List existing routes (run from within the container):
```
curl localhost:9180/apisix/admin/routes -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1"
```

