mkdir tmp_work_dir
git clone https://github.com/maddyonline/glot-code-runner tmp_work_dir
docker run --rm -it -v $PWD/tmp_work_dir:/go/src/glot -v $PWD:/dest -w /go/src/glot  golang go build -o /dest/runner runner.go
rm -rf tmp_work_dir