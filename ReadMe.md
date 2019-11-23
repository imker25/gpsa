# gpsa - A GPX Statistic Analysing tool
This is a simple comandline tool, that may helps to extract statistic data out of *.gpx files
- [gpsa - A GPX Statistic Analysing tool](#gpsa---a-gpx-statistic-analysing-tool)
- [User Documentaion](#user-documentaion)
  - [Installation](#installation)
  - [Usage](#usage)
- [Development](#development)
  - [Build](#build)
  - [Hints for VSCode Users](#hints-for-vscode-users)
- [Geo Math Internals](#geo-math-internals)
- [License](#license)
# User Documentaion

## Installation
On Linux do the following in a empty folder:
```sh
wget https://homer.tobi.backfrak.de/jenkins/job/GPSA/job/master/lastSuccessfulBuild/artifact/bin/gpsa && chmod 770 ./gpsa

```

On Windows download [gpsa.exe](https://homer.tobi.backfrak.de/jenkins/job/GPSA/job/master/lastSuccessfulBuild/artifact/bin/gpsa.exe) and execute it on comandline.

## Usage
You might want to call ```-help```, to find out how to use the program.

```
~$ ./gpsa -help                                                     
./gpsa: Reads in GPS track files, and writes out basic statistic data found in the track as a CSV style report
Program Version: 0.2.1+afb46fa

Usage: ./bin/gpsa [options] [files]
  files
        One or more track files (only *.gpx) supported at the moment
Options:
  -depth string
        Tell how depth the program should analyse the files. Possible values are [segment file track ] (default "track")
  -dont-panic
        Tell if the prgramm will exit with panic, or with negiatv exit code in error cases (default true)
  -help
        Prints this help message
  -license
        Print the license information of the program
  -out-file string
        Tell where to write the output. StdOut is used when not set
  -print-csv-header
        Print out a csv header line (default true)
  -skip-error-exit
        Don't exit the program on track file processing errors
  -verbose
        Run the program with verbose output
  -version
        Print the version of the program
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

# Geo Math Internals 
The Geografic calculations are done with the  ```haversine formula```  as descripted [here](http://www.movable-type.co.uk/scripts/latlong.html). My Impelemtaion will ignore atitute difference for distance bigger then 33 km, by checking the agular distance to be bigger then 0.3°. For smaller distances the program will add the atitute difference using the ```pythagoras theorem```.

# License

gpsa is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)
May got some testfiles and ideas from [gpxgo](https://github.com/tkrajina/gpxgo/tree/master/test_files)

