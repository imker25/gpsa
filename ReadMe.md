# gpsa - A GPX Statistic extracting tool

This is a simple command line tool that helps to extract data for statistical analysis out of `*.gpx` and `*.tcx` files. You might want to use this program to extract data like `Distance`, `ElevationGain` or `AverageSpeed` from a bunch of `*.gpx` or `*.tcx` files and store this data in a *.csv file for further analysis.

- [gpsa - A GPX Statistic extracting tool](#gpsa---a-gpx-statistic-extracting-tool)
  - [User Documentation](#user-documentation)
    - [Installation](#installation)
    - [Usage](#usage)
      - [Help Command](#help-command)
      - [Examples](#examples)
      - [Output Values explained](#output-values-explained)
  - [Development](#development)
    - [Build](#build)
    - [Hints for VSCode Users](#hints-for-vscode-users)
  - [Geo Math Internals](#geo-math-internals)
  - [License](#license)

## User Documentation

Some documentation and examples about gpsa usage.

### Installation

Download the latest release executable from [GitHub Release](https://github.com/imker25/gpsa/releases/latest) according to your operation system and copy the file to a empty folder.

**On Linux**, open a command line and set the file executable:

```sh
chmod 770 ./gpsa

# You may want to get the programs help
./gpsa -help
```

**On Windows**, open the command window in this folder and execute the program (may with the -help flag):

```batch
.\gpsa.exe -help
```

### Usage

Below you can find some kind of user manual for this program.

#### Help Command

You might want to call ```-help``` to find out how to use the program.

```txt
~$ ./gpsa -help
./gpsa: Reads in GPS track files, and writes out basic statistic data found in the track as a CSV style report
Program Version: 0.9.0+7b79520

Usage: ./gpsa [options] [files]
  files
        One or more track files (only  *.gpx and *.tcx) supported at the moment
Options:
  -correction string
        Define how to correct the elevation data read in from the track. Possible values are [steps linear none ] (default "steps")
  -depth string
        Define the way the program should analyse the files. Possible values are [segment file track ] (default "track")
  -dont-panic
        Define if the programm will exit with panic or with a negativ exit code in error cases. Possible values are [true false] (default true).
  -help
        Prints this help message
  -license
        Print the license information of the program 
  -minimal-moving-speed float
        The minimal speed. Distances traveled with less speed are not counted. In [m/s] (default 0.3)
  -minimal-step-hight float
        The minimal step hight. Only in use when "steps"  elevation correction is used. In [m] (default 10)
  -out-file string
        Define where to write the output. (default "StdOut" if not explicitly set)
  -print-csv-header
        Print out a csv header line. Possible values are [true false] (default true).
  -print-elevation-over-distance
        Tell if "ElevationOverDistance.csv" should be created for each track. The files will be locate in tmp dir.
  -skip-error-exit
        Use this flag if you don't want to abort the program during track file processing errors
  -summary string
        Tell if you want to get a summary report. Possible values are [only additional none ] (default "none")
 -suppress-duplicate-out-put
        Suppress the output of duplicate output lines. Duplicates are detected by timestamps. Output with non valid time data may still contains duplicates.
  -verbose
        Run the program with verbose output
  -version
        Print the version of the program
```

#### Examples

Simple call with one file:

```sh
~$  ./gpsa my/test/file.gpx
Name; StartTime; EndTime; TrackTime (xxhxxmxxs); Distance (km); HorizontalDistance (km); AltitudeRange (m); MinimumAltitude (m); MaximumAltitude (m); ElevationGain (m); ElevationLose (m); UpwardsDistance (km); DownwardsDistance (km); MovingTime (xxhxxmxxs); UpwardsTime (xxhxxmxxs); DownwardsTime (xxhxxmxxs); AverageSpeed (km/h); UpwardsSpeed (km/h); DownwardsSpeed (km/h);
my/test/file.gpx: 2020-01-29 09:28:06;2020-01-29T08:28:10Z; 2020-01-29T13:48:07Z; 5h19m53s; 94.750000; 90.250000; 1188.370000; 821.610000; 2009.980000; 10659.340000; -10884.500000; 43.470000; 51.000000; 4h1m44s; 2h8m49s; 1h52m55s; 23.520000; 20.250000; 27.100000; 

```

Simple call with multiple files:

```sh
~$  ./gpsa my/test/01.gpx my/test/02.tcx my/test/03.gpx
Name; StartTime; EndTime; TrackTime (xxhxxmxxs); Distance (km); HorizontalDistance (km); AltitudeRange (m); MinimumAltitude (m); MaximumAltitude (m); ElevationGain (m); ElevationLose (m); UpwardsDistance (km); DownwardsDistance (km); MovingTime (xxhxxmxxs); UpwardsTime (xxhxxmxxs); DownwardsTime (xxhxxmxxs); AverageSpeed (km/h); UpwardsSpeed (km/h); DownwardsSpeed (km/h);
03.gpx: Tulln - Wien; not valid; not valid; not valid; 37.640000; 36.120000; 48.000000; 158.000000; 206.000000; 52.000000; -26.000000; 17.520000; 14.060000; not valid; not valid; not valid; not valid; not valid; not valid;
02.gpx: 2019-08-18 11:07:40; 2019-08-18T09:11:01Z; 2019-08-18T15:47:34Z; 1h35m40s; 37.820000; 37.230000; 104.090000; 347.020000; 451.110000; 263.880000; -251.430000;17.860000; 19.760000; 1h33m20s; 47m54s; 44m56s; 24.320000; 22.370000; 26.390000; 

```

Get statistics for a number of files into a csv output:

```sh
~$  ./gpsa -out-file=gps-statistics.csv my/test/*.gpx

```

Get statistics for a number of files into a csv output, with verbose comandline output:

```sh
~$  ./gpsa -verbose -out-file=gps-statistics.csv my/test/*.gpx
Call:  ./gpsa -verbose -out-file=gps-statistics.csv my/test/01.gpx my/test/02.gpx my/test/03.gpx
Version: 0.9.0+7b79520
Read file: my/test/01.gpx
Read file: my/test/02.gpx
Read file: my/test/03.gpx
3 of 3 files process successfull

```

#### Output Values explained

Below is a list of the output values and what they mean:

- `Name`: The name of the output line. Either read from the file, or if no name tag is set in the file it will be calculated out of the filename and the file hierarchy.
- `StartTime`: The time the track started. *not valid* in case we detect no or invalid time data.
- `EndTime`: The time the track ended. *not valid* in case we detect no or invalid time data.
- `TrackTime`: The time between `StartTime` and `EndTime`. Formatted as `xxhxxmxxs`. *not valid* in case we detect no or invalid time data.
- `Distance`: The distance of the track measured in `km`.
- `HorizontalDistance`: The horizontal distance of the track measured in `km`. This value ignores the vertical component of the distance like most GPS tools do.
- `AltitudeRange`: The range between the highest and the lowest point. Measured in `m`.
- `MinimumAltitude`: The altitude of the lowest point. Measured in `m`.
- `MaximumAltitude`: The altitude of the highest point. Measured in `m`.
- `ElevationGain`: The total sum of all upwards vertical distance.  Measured in `m`.
- `ElevationLose`: The total sum of all downwards vertical distance.  Measured in `m`.
- `UpwardsDistance`: The total sum of all distance moved upwards.  Measured in `km`.
- `DownwardsDistance`: The total sum of all distance moved downwards.  Measured in `km`.
- `MovingTime`: The time spend moving. Formatted as `xxhxxmxxs`. *not valid* in case we detect no or invalid time data.
- `UpwardsTime`: The time spend moving upwards. Formatted as `xxhxxmxxs`. *not valid* in case we detect no or invalid time data.
- `DownwardsTime`: The time spend downwards. Formatted as `xxhxxmxxs`. *not valid* in case we detect no or invalid time data.
- `AverageSpeed`: The average speed. Calculated from `Distance` and `MovingTime`. Measured in  `km/h`. *not valid* in case we detect no or invalid time data.
- `UpwardsSpeed`: The average speed during upwards movement . Measured in  `km/h`. *not valid* in case we detect no or invalid time data.
- `DownwardsSpeed`: The average speed during downwards movement . Measured in  `km/h`. *not valid* in case we detect no or invalid time data.

## Development

To develop this software install [Go](https://golang.org/) and [Gradle](https://gradle.org/) on your machine.

### Build

Use [Gradle](https://gradle.org/) to build and test the project.

```sh
gradle build        # build the project
gradle build test   # build and run the tests for the project
gradle test         # test the project
```

### Hints for VSCode Users

If you use [VS Code](https://code.visualstudio.com/) for GO development, you might find the following example settings useful.

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
                "${workspaceRoot}/testdata/valid-gpx/02.gpx"
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

## Geo Math Internals

The Geografic calculations are done with the  ```haversine formula```  as described [here](http://www.movable-type.co.uk/scripts/latlong.html). My implementation will ignore altitude differences for distances bigger than 33km by checking the angular distance which in this case then is bigger than 0.3°. For smaller distances the program will add the altitude difference using the ```pythagoras theorem```.

## License

gpsa is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).
May got some testfiles and ideas from [gpxgo](https://github.com/tkrajina/gpxgo/tree/master/test_files)
