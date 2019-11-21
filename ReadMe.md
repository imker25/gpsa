# gpsa - A GPX Statistic Analysing tool
This is a simple comandline tool, that may helps to extract statistic data out of *.gpx files

# Usage
You might want to call ```-help```, to find out how to use the program.

```sh
./gpsa -help                                                     
./gpsa: Reads in GPS track files, and writes out basic statistic data found in the track as *.csv

Usage: ./gpsa [options] [files]
  files
        One or more track files (only *.gpx) supported at the moment
Options:
  -depth string
        Tell how depth the program should analyse the files. Possible values are [segment file track ] (default "track")
  -dont-panic
        Tell if the prgramm will exit with panic, or with negiatv exit code in error cases (default true)
  -help
        Prints this message
  -out-file string
        Tell where to write the output. StdOut is used when not set
  -print-csv-header
        Print out a csv header line (default true)
  -skip-error-exit
        Don't exit the programm on track file processing errors
  -verbose
        Run the programm with verbose output
```

# Development
To develop this software install [Go](https://golang.org/) and [Gradle](https://gradle.org/) on your machine. 

## Build
Use [Gradle](https://gradle.org/) to build and test the porject

```sh
gradle build        # build the project
gradle build test   # build and run the tests for the project
gradle test         # test the project
```

## Hints for VSCode Users
If you use [VS Code](https://code.visualstudio.com/) for GO development, you might find the following example settings usefull.

The ```tasks.json```:
```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Build",
            "type": "shell",
            "command": "gradle build test",
            "group": {
                "kind": "build",
                "isDefault": true
            }
        },
        {
            "label": "Test",
            "type": "shell",
            "command": "gradle test",
            "group": {
                "kind": "test",
                "isDefault": true
            }
        }
    ]
} 
```

The ```launch.json```:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "debug",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceRoot}/src/tobi.backfrak.de/cmd/gpsa/main.go",
            "cwd": "${workspaceRoot}",
            "args": [
          //      "-dont-panic=false",
                "-out-file=/dev/shm/test.csv",
                "${workspaceRoot}/testdata/valide-gpx/02.gpx"
            ]
        }
    ]
}
```

The ```settings.json```:
```json
{
      "go.gopath": "${env:GOPATH}:${workspaceFolder}",
}
```

## Internals
The Geografic calculations are done with the  ```haversine formula```  as descripted [here](http://www.movable-type.co.uk/scripts/latlong.html. My Impelemtaion will ignore atitute difference for distance bigger then 33 km, by checking the agular distance to be bigger then 0.3°. For smaller distances the program will add the atitute difference using the ```pythagoras theorem```.

# License

gpsa is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)
May got some testfiles from [gpxgo](https://github.com/tkrajina/gpxgo/tree/master/test_files)