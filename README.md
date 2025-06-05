## meilisearch

docker

```shell
docker run -dit --name meilisearch --restart=always -p 7700:7700 -e MEILI_MASTER_KEY='root123456' -v $PWD/meilisearch/ meili_data:/meili_data getmeili/meilisearch:v1.14
```


```shell
# wire
go install github.com/google/wire/cmd/wire

# protobuf
winget install --id=Google.Protobuf --proxy http://127.0.0.1:10808

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# grpc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```


## Meilisearch

```shell
# 开启 contains 预览特性
# containsFilter: 允许在过滤器中使用 contains 操作符(同时启用start with)
# editDocumentsByFunction: 允许通过函数编辑文档
curl -X PATCH 'MEILISEARCH_URL/experimental-features/' -H 'Content-Type: application/json' \
--data-binary '{
"containsFilter": true,
"editDocumentsByFunction": true,
}'
```