# Battle of Bits

[![Build Status](https://travis-ci.org/battleofbits/arena.png?branch=master)](https://travis-ci.org/battleofbits/arena)

Battle of Bits is a hosted service for competitive board game AI programs. We
host daily tournaments and continual round-robin matches with your robot.

## How it Works

After you sign up on the website, create a profile for your bot. For each bot
your create, you'll need to associate a URL with that bot so that we can talk
to it. Each game a bot plays will need a different URL.

Should a bot only exist for a single type of game? Or should a bot work across
multiple games?


## Developing Battle of Bits

If you wish to work on Battle of Bits itself, you'll first need
[Go](http://golang.org) installed (version 1.2+ is _required_).
Make sure you have Go properly installed, including setting up your
[GOPATH](http://golang.org/doc/code.html#GOPATH).

You'll also need [`lib/pq`](https://github.com/lib/pq) to compile packer. You
can install that with:

```
$ go get -u github.com/lib/pq
```

Next, clone this repository into `$GOPATH/src/github.com/battleofbits/arena`
and then just type `make`. In a few moments, you'll have a working `arena`
executable:

```
$ make
...
$ bin/arena
...
```

You can run tests by typing `make test`.

If you make any changes to the code, run `make format` in order to
automatically format the code according to Go standards.

## Adding Games

See [the readme in the `games` subfolder][games-readme].

 [games-readme]: https://github.com/battleofbits/arena/tree/master/games
