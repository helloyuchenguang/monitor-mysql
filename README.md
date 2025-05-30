## meilisearch

docker

```shell
docker run -dit --name meilisearch \
--restart=always \
-p 7700:7700 \
-e MEILI_MASTER_KEY='root123456'\
-v $PWD/meilisearch/meili_data:/meili_data \
getmeili/meilisearch:v1.14
```
