# cli-streamer

Simple & light Go-based cross-platform tool to execute commands from shell and from different types (json, yaml, toml etc) of files with the capabilities of multi-streaming the output to multiple destinations files including the standard console with the posibility of adding an execution timeout value for each command. 



## Table of contents
* [Description](#description)
* [Technologies](#technologies)
* [Setup](#setup)
* [Usage](#usage)
* [Upcomings](#upcomings)
* [Contribution](#contribution)
* [License](#license)


## Description

Very simple. Please have a look at the [usage section](#usage) for examples.
This tool can be integrated into an automated script for more capabilities.



## Technologies

This project is developed with:
* Golang version: 1.13
* Native libraries only


## Setup

On Windows, Linux macOS, and FreeBSD you will be able to download the pre-built binaries once available.
If your system has [Go even < 1.7](https://golang.org/dl/) you can pull the codebase and build from the source.

```
# build the cli-streamer program on windows
git clone https://github.com/jeamon/cli-streamer.git && cd cli-streamer
go build -o cli-streamer.exe cli-streamer.go

# build the cli-streamer program on linux and others
git clone https://github.com/jeamon/cli-streamer.git && cd cli-streamer
go build -o cli-streamer cli-streamer.go
```


## Usage


```Usage:
    
cli-streamer [-task <quoted-command-to-execute>] [-timeout <execution-deadline-in-seconds>] [-files <filenames-to-stream-output>] [-save] [-console]

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

You have to provide at least one mandatory argument value [-task]. The list of files
where to have the output duplicated (in real-time) must be in one word and space
separed. To have the output be saved into a daily file named outputs-<year><month><day>,
just add -save flag when launching the program. Upcoming version will add capabilities
to mention multiple tasks and each task with its own destinations output filenames.
In the meantime please see below examples for current version 1.0 :


Examples:

    Option: single command execution
    
    $ cli-streamer -task "netstat -n 2 | findstr ESTAB" -timeout 180 -files "a.txt b.txt" -save
    $ cli-streamer -task "ping 127.0.0.1 -t" -timeout 3600 -files "ping.txt" --console
    $ cli-streamer -task "journalctl -f | grep <xx>" -timeout 120 -files "proclog.txt" --save
    $ cli-streamer -task "tail -f /var/log/syslog" -timeout 3600 -files "syslog.txt"

    Option: multiple commands from file
    
    $ cli-streamer -tasksFile "tasks.txt"

    Option: multiple commands from json file
    
    $ cli-streamer -tasksJson "tasks.txt"

    Option: multiple commands from toml file
    
    $ cli-streamer -tasksToml "tasks.txt"

    Option: multiple commands from yaml file
    
    $ cli-streamer -tasksYaml "tasks.txt"
	
```


## Upcomings

* add capabilities to specify mutiple commands and each with its own destinations files as options and timeout value.
* add capabilities to specify interval at which we want to view output on the terminal and for which command.
* add capabilities to pause & restart or stop output display at console screen.
* add capabilities to specify multiple files containing a list of task to execute.
* add capabilities to specify if wanted to add file to a daily folder to current folder.


## Contribution

Pull requests are welcome. However, I would be glad to be contacted for discussion before.


## License

please check & read [the license details](https://github.com/jeamon/cli-streamer/blob/master/LICENSE) or [reach out to me](https://blog.cloudmentor-scale.com/contact) before any action.