
# Docker Toolkit

```
gvm use go1.18.1

```

```
./docker/tookit/create-toolkit-image.sh
./docker/tookit/build-indexer.sh
./docker/tookit/create-indexer-image.sh

```


```
serving on :8980
{"level":"info","msg":"serving on :8980","time":"2022-07-04T00:44:43Z"}
{"level":"info","msg":"Running 0 migrations.","time":"2022-07-04T00:44:43Z"}
⇨ http server started on [::]:8980
{"level":"info","msg":"Setting status: Migrations Complete","time":"2022-07-04T00:44:43Z"}
{"level":"info","msg":"Migration finished successfully.","time":"2022-07-04T00:44:43Z"}
{"level":"info","msg":"loading genesis file genesis.json","time":"2022-07-04T00:44:43Z"}
{"error":"error decoding genesis, json decode error [pos 652]: no matching struct field found when decoding stream map with key stprf","level":"error","msg":"genesis.json: could not load genesis json, error decoding genesis, json decode error [pos 652]: no matching struct field found when decoding stream map with key stprf","time":"2022-07-04T00:44:43Z"}

git pull upstream master
rm third_party/ .gopath/ -Rf # Clean & Build 

```

```
indexer data directory was not providedStarting indexer against algod.

add --data-dir /tmp on start.sh

```
