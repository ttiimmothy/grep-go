[![progress-banner](https://backend.codecrafters.io/progress/grep/28176ce0-63c3-4817-aa12-e6df9c6ea2f8)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

# Grep Go

This is a starting point for Go solutions to the
["Build Your Own grep" Challenge](https://app.codecrafters.io/courses/grep/overview).

[Regular expressions](https://en.wikipedia.org/wiki/Regular_expression)
(Regexes, for short) are patterns used to match character combinations in
strings. [`grep`](https://en.wikipedia.org/wiki/Grep) is a CLI tool for
searching using Regexes.

In this challenge you'll build your own implementation of `grep`. Along the way
we'll learn about Regex syntax, how parsers/lexers work, and how regular
expressions are evaluated.

**Note**: If you're viewing this repo on GitHub, head over to
[codecrafters.io](https://codecrafters.io) to try the challenge.

## Passing the first stage

The entry point for your `grep` implementation is in `cmd/mygrep/main.go`. Study
and uncomment the relevant code, and push your changes to pass the first stage:

```sh
git add .
git commit -m "pass 1st stage" # any msg
git push origin master
```

Time to move on to the next stage!

## Stage 2 & beyond

Note: This section is for stages 2 and beyond.

1. Ensure you have `go (1.19)` installed locally
2. Run `./your_grep.sh` to run your program, which is implemented in
   `cmd/mygrep/main.go`.
3. Commit your changes and run `git push origin master` to submit your solution
   to CodeCrafters. Test output will be streamed to your terminal.

## License

Grep Go is licensed under [GNU General Public License v3.0](LICENSE).