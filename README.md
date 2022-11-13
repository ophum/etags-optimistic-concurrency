# etags-optimistic-concurrency

etag を用いた楽観的排他制御を雑に実装。

`GET /api/v1/pets/:id` で取得する際に Response Header に ETag がセットされる。

`PUT /api/v1/pets/:id` で更新する際に Request Header の If-Match に取得の際にセットされていた ETag をセットする。

If-Match の値と現在の entity の hash 値が異なる場合は、`412 Precondition Failed` エラーになる。

## 参考

- [mdn web docs ETag](https://developer.mozilla.org/ja/docs/Web/HTTP/Headers/ETag)
