# code

```sh
git clone https://github.com/maddyonline/code
cd code/
./pull-docker-images.sh
./install_runner.sh
go install ./...
code-server -h
code-server -runner=runner
```
