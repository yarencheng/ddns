```
docker build -t yarencheng/ddns .
```

```
docker run -d --restart always \
    --name ddns \
    -v $PWD/sa.json:$PWD/sa.json \
    -e GOOGLE_SERVICE_ACCOUNT_KEY_JSON_BASE64=$(cat sa.json | base64 -w 0) \
    -e PROJECT_ID=my_gcp_project_id \
    -e MANAGED_ZONE=my_zone \
    -e NAME=example.local \
    yarencheng/ddns
```
