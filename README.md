# gostream

Simple & light Go-based cross-platform tool to execute commands from shell and from different types (json, yaml, toml) of files with the capabilities of multi-streaming each output to multiple & different destinations files including the standard console with the option to define an execution timeout value for each command.



## Table of contents
* [Description](#description)
* [Setup](#setup)
* [Usage](#usage)
* [Upcomings](#upcomings)
* [Contribution](#contribution)
* [License](#license)


## Description

Very simple. Please have a look at the [usage section](#usage) for examples.
This tool can be integrated into an automated script for more capabilities.


## Setup

On Windows, Linux macOS, and FreeBSD you will be able to download the pre-built binaries once available.
If your system has [Go even < 1.7](https://golang.org/dl/) you can pull the codebase and build from the source.

```
# build the gostream program on windows
git clone https://github.com/jeamon/gostream.git && cd gostream
go build -o gostream.exe gostream.go

# build the gostream program on linux and others
git clone https://github.com/jeamon/gostream.git && cd gostream
go build -o gostream gostream.go
```


## Usage


```Usage:
    
Option: define single task from shell

    gostream [-task <quoted-command-to-execute>] [-timeout <execution-deadline-in-seconds>] [-files <filenames-to-stream-output>] [-save] [-console]

Option: define multiple tasks from files

    gostream [-taskFile "xfile zfile morefiles"]
    gostream [-taskJson "xfile.json zfile.json more.json"]
    gostream [-taskYaml "xfile.yaml zfile.yaml more.yaml"]
    gostream [-taskToml "xfile.toml zfile.toml more.toml"]

Subcommands:
    version    Display the current version of this tool.
    help       Display the help - how to use this tool.


Options:
    -task      Specify the command with its arguments to run in quoted format.
    -timeout   Specify the number of seconds to allow the task to be running.
    -files     Specify all filenames to stream the output of the execution.
    -save      If present then execution output must also be saved in daily file.
    -console   If present then execution output will also be displayed on terminal.
    -tasksFile Load commands from the provided file (see "example.file" content).
    -tasksJson Load commands from the provided json file (see "example.json" content).
    -tasksToml Load commands from the provided toml file (see "example.toml" content).
    -tasksYaml Load commands from the provided yaml file (see "example.yaml" content).
    

Arguments:
    quoted-command-to-execute      complete command to execute from the shell.
    execution-deadline-in-seconds  number of seconds to have the task running.
    filenames-to-stream-output     space separed filenames where to stream output.

Using the shell option, you must provide at least one mandatory argument value [-task].
The list of files to stream the output must be in one word and space separed. To have
the output be saved into a daily filename pattern outputs.<year><month><day>,
just add -save flag when launching the program. Use options which support filenames
input if needed to run multiple tasks in parallel. Upcoming version will add capabilities
to mention multiple tasks (with its their own attributes) from the shell.
In the meantime please see below examples for current version 1.0 :


Examples:

    Option: single command execution
    
    $ gostream -task "netstat -n 2 | findstr ESTAB" -timeout 180 -files "a.txt b.txt" -save
    $ gostream -task "ping 127.0.0.1 -t" -timeout 3600 -files "ping.txt" --console
    $ gostream -task "journalctl -f | grep <xx>" -timeout 120 -files "proclog.txt" --save
    $ gostream -task "tail -f /var/log/syslog" -timeout 3600 -files "syslog.txt"

    Option: multiple commands from file(s)
    
    $ gostream -tasksFile "tasks.txt others.file"

    Option: multiple commands from json file(s)
    
    $ gostream -tasksJson "tasks.json more.json"

    Option: multiple commands from toml file(s)
    
    $ gostream -tasksToml "tasks.toml others.toml"

    Option: multiple commands from yaml file(s)
    
    $ gostream -tasksYaml "tasks.yaml more.yaml"
	
```


## Upcomings

* add capabilities to load mutiple commands fully defined from env variables.
* add capabilities to specify outputs display interval and for which commands.
* add capabilities to uniquely color each command outputs when displayed.
* add capabilities to pause & restart or stop output display on the terminal.


## Contribution

Pull requests are welcome. However, I would be glad to be contacted for discussion before.


## License

please check & read [the license details](https://github.com/jeamon/gostream/blob/master/LICENSE) or [reach out to me](https://blog.cloudmentor-scale.com/contact) before any action.