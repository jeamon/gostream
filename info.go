package main

const version = "This tool is <gostream> • version 1.0 By Jerome AMON"

const usage = `Usage:
    
Option: define single task from shell

    gostream [-task <quoted-command-to-execute>] [-timeout <execution-deadline-in-seconds>] [-files <filenames-to-stream-output>] [-save] [-console]

Option: define multiple tasks from files

    gostream [-taskFile <xfile> <zfile> <morefiles>]
    gostream [-taskJson <xfile.json> <zfile.json> <more.json>]
    gostream [-taskYaml <xfile.yaml> <zfile.yaml> <more.yaml>]
    gostream [-taskToml <xfile.toml> <zfile.toml> <more.toml>]

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

    $ gostream version
    $ gostream --help
    
    Option: single command execution
    
    $ gostream -task "netstat -n 2 | findstr ESTAB" -timeout 180 -files "a.txt b.txt" -save
    $ gostream -task "ping 127.0.0.1 -t" -timeout 3600 -files "ping.txt" --console
    $ gostream -task "journalctl -f | grep <xx>" -timeout 120 -files "proclog.txt" --save
    $ gostream -task "tail -f /var/log/syslog" -timeout 3600 -files "syslog.txt"

    Option: multiple commands from file
    
    $ gostream -tasksFile "tasks.txt"

    Option: multiple commands from json file
    
    $ gostream -tasksJson "tasks.txt"

    Option: multiple commands from toml file
    
    $ gostream -tasksToml "tasks.txt"

    Option: multiple commands from yaml file
    
    $ gostream -tasksYaml "tasks.txt"`
