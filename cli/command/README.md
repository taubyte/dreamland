# Command utils

This module is for adding functionality to commands

Say you want a function to take a name, you define the command
```go
someCommand := &cli.Command{
    Name: "some-command"
}
```
Now before returning if you want to attach name argument parsing ability to this command

```go
command.Name(someCommand)
```

This empty wrapping will attribute args[0] or the name flag to `name`, and return an
error if the name is not found
