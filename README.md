# code

### Prerequisites
1. A running docker daemon:

    On Mac, one can do the following:
    ```
    docker-machine start default
    docker-machine env default
    eval "$(docker-machine env default)"
    ```

### Running various commands

```sh
git clone https://github.com/maddyonline/code
cd code/
./pull-docker-images.sh
./install_runner.sh
go install ./...
code-server -h
code-server -runner=runner
```
