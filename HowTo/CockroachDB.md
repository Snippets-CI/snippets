# CockroachDB Steps

## CockroachDB for Windows

### Installation

* Download [CockroachDB for Windows](https://binaries.cockroachdb.com/cockroach-v20.1.0.windows-6.2-amd64.zip)
* Download [Zone info](https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip)
  * 1. Set `ZONEINFO` environment variable to point to the zip file

### Usage

**Start node:**  

```
cockroach start \
--insecure \
--store=node1 \
--listen-addr=localhost:26257 \
--http-addr=localhost:8080 \
--join=localhost:26257,localhost:26258,localhost:26259
```

Join will only be needed of 2 ore more nodes will be up.

**Explanation:**  
`--insecure` communication is unencrypted  
`--store` indicated the location where the node's data and logs are stored  
`--listen-addr` listen only on localhost with port 26257 used for internal and client traffic  
`--http-addr` listen for HTTP request from the Admin UI

**Use the built-in SQL client:**  

```
cockroach sql --insecure --host=localhost:26257
```


cockroach start --insecure --store=node1 --listen-addr=localhost:26257 --http-addr=localhost:8080
