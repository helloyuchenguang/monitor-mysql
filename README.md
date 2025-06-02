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
# grpc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```