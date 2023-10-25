# `ht` -- Hash Tree

`ht` computes a hash tree ([Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree)) using `sha256`
for a given directory and all of its descendants.

## `ht` CLI Tool Installation Guide

### Prerequisites

Before you begin, make sure you have the following prerequisites:

- Go programming language installed on your system. You can download and install Go from [the official Go website](https://golang.org/dl/).

### Installation

#### Using brew

1. Open your terminal.

2. `brew install aaraney/tap/ht`

3. Run `ht`.

    ```shell
    ./ht --help
    ```

#### Using go install

1. Open your terminal.

2. `go install github.com/aaraney/ht@latest`

3. Run `ht`.

    ```shell
    ./ht --help
    ```

#### For Development

1. Open your terminal.

2. Clone the repo.

    ```shell
    git clone git@github.com:aaraney/ht.git
    # or
    git clone https://github.com/aaraney/ht.git

    # then cd
    cd ht
    ```

3. Build `ht`.

    ```shell
    go build
    ```

4. Run `ht`.

    ```shell
    ./ht --help
    ```

#### Once the repo is public the following will work

Follow these steps to install `ht`:

1. Open your terminal.

2. Use the `go get` command to download and install `ht`:

   ```shell
   go get -u github.com/aaraney/ht
   ```

3. Verify the installation by running the following command:

   ```bash
   ht --help
   ```

### Usage

```shell
ht --help

Usage of ht:
  -n int
        Maximum number of workers. Defaults to number of cpus. (default 10)

# running ht on the repo
ht
0c74421ad3e6aca9eefaa02cb8b50772f32179a88c098f8a2e6e7288d51426d9 ./
254161b1da36140336d88a41b31d4d0aff3803e55599b03d3445174d8b06cbd3 ./node_test.go
3032e21626454f9914cb863b47f054899726c30682f35cda782055b3404b1cba ./LICENSE
30bd74413565e2bd817f7e4565f8ede6288a5c872c5a1b94f7990d81f6c2d8a1 ./ht
42e9942c4b41c70b04cad8db29c3537f67aefc3c987aa0fdfb6f1f161bf01bff ./go.mod
4a1055e0a39f836fcc6eec5b93b980172d1c02f589241f85869b4746ea69dc94 ./hash_files.go
a4d2a24977591fcd4d982556c83daa10f44d25718d436d30c3e9ea362199ce4b ./main.go
b002061e0ef1ac4b7e6f777463e49033e5d457a81ca0708fc2a968bf263ea90b ./merkle_tree.go
cf237c7aff44efbe6e502e645c3e06da03a69d7bdeb43392108ef3348143417e ./.gitignore
```

### License

This Go CLI tool is distributed under the MIT [license](LICENSE). Please refer to the project's
repository for more details on licensing.
