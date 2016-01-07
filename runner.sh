git clone https://github.com/maddyonline/glot-code-runner
docker pull rsmmr/clang
docker run --rm -it -v $PWD/glot-code-runner:/go/src/glot -w /go/src/glot  golang go build runner.go
cat glot-code-runner/test_with_stdin.txt | docker run --rm -i -v "$(pwd)/glot-code-runner":/app -w /app rsmmr/clang ./runner