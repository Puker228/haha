# haha

Terminal CLI that prints a character face and can animate it as rain or an explosion.

## Install

Install with Go:

```sh
go install github.com/Puker228/haha@latest
```

## Usage

Print one trollface:

```sh
go run .
```

Choose a character with `-c` or `--character`:

```sh
go run . -c joker
```

Run falling trollface rain:

```sh
go run . rain
```

Run falling rain with another character:

```sh
go run . rain -c joker
```

Explode trollface art across the terminal:

```sh
go run . boom
```

Explode another character:

```sh
go run . boom -c joker
```

Available characters:

- `trollface`
- `joker`

In animation modes, press `q` or `ctrl+c` to quit.
