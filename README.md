# Lucy
Helper go app


```sh
make build
make up
make console
```
inside console
```sh
make go
```
close console [ctrl c]
```sh
chmod +x ./bin/lucy [-c cluster_name(s)]
cp ./bin/lucy ~/bin/lucy
cp ./forwarder.dist.yaml ~/bin/forwarder.yaml
```
edit `~/bin/forwarder.yaml`

example:
```sh
lucy                   # all clusters from forwarder.yaml
lucy -c cluster1
lucy -c cluster1,cluster2
