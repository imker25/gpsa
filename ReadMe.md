# gpsa - A GPX Statistic extracting tool

This is a simple command line tool that helps to extract data for statistical analysis out of `*.gpx` and `*.tcx` files. You might want to use this program to extract data like `Distance`, `ElevationGain` or `AverageSpeed` from a bunch of `*.gpx` or `*.tcx` files and store this data in a *.csv or *.json file for further analysis.

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
./bin/gpsa: Reads in GPS track files, and writes out basic statistic data found in the track as a report
Program Version: 2.3.6-4142b3b

Usage: ./bin/gpsa [options] [files]
  files
        One or more track files of the following type: *.tcx, *.gpx, 
Options:
  -correction string
    	Define how to correct the elevation data read in from the track. Possible values are [steps linear none ] (default "steps")
  -depth string
    	Define the way the program should analyse the files. Possible values are [segment file track ] (default "track")
  -dont-panic
    	Decide if the program will exit with panic or with negative exit code in error cases. Possible values are [true false] (default true)
  -help
    	Print help message and exit
  -license
    	Print license information of the program and exit
  -minimal-moving-speed float
    	The minimal speed. Distances traveled with less speed are not counted. In [m/s] (default 0.3)
  -max-start-time string
        The maximum StartTime for a track to be added to the output. Formatted in "YYYY-MMM-dd HH:mm:ss", may without seconds or just a date
  -min-start-time string
        The minimum StartTime for a track to be added to the output. Formatted in "YYYY-MMM-dd HH:mm:ss", may without seconds or just a date      
  -minimal-step-hight float
    	The minimal step hight. Only in use when "steps"  elevation correction is used. In [m] (default 10)
  -out-file string
    	Decide where to write the output. StdOut is used when not explicitly set. Supported file endings are: *.json, *.csv, . The format will be set according the given ending.
  -print-csv-header
    	Print out a csv header line. Possible values are [true false] (default true)
  -print-elevation-over-distance
    	Tell if "ElevationOverDistance.csv" should be created for each track. The files will be locate in tmp dir.
  -skip-error-exit
    	Don't exit the program on track file processing errors
  -std-out-format string
    	The output format when stdout is the used output. Ignored when out-file is given. Possible values are [JSON CSV ] (default "CSV")
  -summary string
    	Tell if you want to get a summary report. Possible values are [only additional none ] (default "none")
  -suppress-duplicate-out-put
    	Suppress the output of duplicate lines. Duplicates are detected by timestamps. Output with non valid time data may still contains duplicates.
  -time-format string
    	Tell how the csv output formater should format times. Possible values are ["Mon Jan _2 15:04:05 MST 2006" "Monday, 02-Jan-06 15:04:05 MST" "2006-01-02T15:04:05Z07:00" ] (default "Monday, 02-Jan-06 15:04:05 MST")
  -verbose
    	Run the program with verbose output
  -version
    	Print version of the program and exit

It is also possible to pipe track file names or track file content into

Examples:
./gpsa my/test/file.gpx
./gpsa -verbose -out-file=gps-statistics.csv my/test/*.gpx
find ./testdata/valid-gpx -name "*.gpx" | ./bin/gpsa -summary=additional -out-file=./test.json
cat  01.gpx 01.tcx 03.tcx 02.gpx | ./bin/gpsa -out-file=./test.json
```

#### Examples

Simple call with one file:

```sh
~$  ./gpsa my/test/file.gpx
Name; StartTime; EndTime; TrackTime (xxhxxmxxs); Distance (km); HorizontalDistance (km); AltitudeRange (m); MinimumAltitude (m); MaximumAltitude (m); ElevationGain (m); ElevationLose (m); UpwardsDistance (km); DownwardsDistance (km); MovingTime (xxhxxmxxs); UpwardsTime (xxhxxmxxs); DownwardsTime (xxhxxmxxs); AverageSpeed (km/h); UpwardsSpeed (km/h); DownwardsSpeed (km/h);
my/test/file.gpx: 2020-01-29 09:28:06;2020-01-29T08:28:10Z; 2020-01-29T13:48:07Z; 5h19m53s; 94.75; 90.25; 1188.37; 821.61; 2009.98; 10659.34; -10884.50; 43.47; 51.00; 4h1m44s; 2h8m49s; 1h52m55s; 23.52; 20.25; 27.10; 

```

Simple call with multiple files:

```sh
~$  ./gpsa my/test/01.gpx my/test/02.tcx my/test/03.gpx
Name; StartTime; EndTime; TrackTime (xxhxxmxxs); Distance (km); HorizontalDistance (km); AltitudeRange (m); MinimumAltitude (m); MaximumAltitude (m); ElevationGain (m); ElevationLose (m); UpwardsDistance (km); DownwardsDistance (km); MovingTime (xxhxxmxxs); UpwardsTime (xxhxxmxxs); DownwardsTime (xxhxxmxxs); AverageSpeed (km/h); UpwardsSpeed (km/h); DownwardsSpeed (km/h);
03.gpx: Tulln - Wien; not valid; not valid; not valid; 37.64; 36.12; 48.00; 158.00; 206.00; 52.00; -26.00; 17.52; 14.06; not valid; not valid; not valid; not valid; not valid; not valid;
02.gpx: 2019-08-18 11:07:40; 2019-08-18T09:11:01Z; 2019-08-18T15:47:34Z; 1h35m40s; 37.82; 37.23; 104.09; 347.02; 451.11; 263.88; -251.43;17.86; 19.76; 1h33m20s; 47m54s; 44m56s; 24.32; 22.37; 26.39; 

```

Get statistics for a number of files into a csv output:

```sh
~$  ./gpsa -out-file=gps-statistics.csv my/test/*.gpx

```

Get statistics for a number of files into a csv output, with verbose comandline output:

```sh
~$  ./gpsa -verbose -out-file=gps-statistics.csv my/test/*.gpx
Call:  ./gpsa -verbose -out-file=gps-statistics.csv my/test/01.gpx my/test/02.gpx my/test/03.gpx
Version: 2.0.1+7b79520
Read file: my/test/01.gpx
Read file: my/test/02.gpx
Read file: my/test/03.gpx
3 of 3 files process successfull

```

Get only the summary report out of a bunch of files

```sh
./bin/gpsa -summary=only my/test/*.gpx
Name; StartTime; EndTime; TrackTime (hh:mm:ss); Distance (km); HorizontalDistance (km); AltitudeRange (m); MinimumAltitude (m); MaximumAltitude (m); ElevationGain (m); ElevationLose (m); UpwardsDistance (km); DownwardsDistance (km); MovingTime (hh:mm:ss); UpwardsTime (hh:mm:ss); DownwardsTime (hh:mm:ss); AverageSpeed (km/h); UpwardsSpeed (km/h); DownwardsSpeed (km/h); 
Sum:; -; -; 54:8:43.442; 775.00; 770.48; -; -; -; 7033.35; -6989.57; 309.34; 359.75; 32:47:16.514; 14:19:21.328; 13:33:19.262; -; -; -; 
Average:; -; -; 2:15:21.810083333; 32.29; 32.10; 134.81; -; -; 293.06; -291.23; 12.89; 14.99; 1:21:58.188083333; 35:48.388666666; 33:53.302583333; 23.78; 21.86; 26.65; 
Minimum:; Sunday, 01-Mar-20 09:27:27 UTC; Sunday, 01-Mar-20 15:11:19 UTC; 53:3; 21.13; 21.07; 63.72; 287.46; 359.73; 159.59; -141.38; 7.92; 8.79; 51:18.444; 21:56.828; 18:22.126; 21.28; 17.73; 23.75; 
Maximum:; Saturday, 31-Oct-20 13:01:16 UTC; Saturday, 31-Oct-20 14:31:42 UTC; 6:4:25; 68.74; 66.54; 801.99; 486.48; 1288.47; 1747.38; -1754.86; 26.45; 34.71; 3:13:47.838; 1:29:31.416; 1:20:40.686; 25.58; 23.78; 28.73;

```

Get statistics from a file as json on stdout

```sh
./bin/gpsa  -std-out-format=json my/test/02.gpx 
{
 "Statistics": [
  {
   "Name": "my/test/02.gpx: 2019-08-18 11:07:40",
   "Data": {
    "Distance": 37823.344979382266,
    "HorizontalDistance": 37741.53944560436,
    "MinimumAltitude": 347.02,
    "MaximumAltitude": 451.11,
    "ElevationGain": 263.88007,
    "ElevationLose": -251.43008,
    "UpwardsDistance": 17858.070360985712,
    "DownwardsDistance": 19761.009730234404,
    "TimeDataValid": true,
    "StartTime": "2019-08-18T09:11:01Z",
    "EndTime": "2019-08-18T15:47:34Z",
    "MovingTime": 5600000000000,
    "UpwardsTime": 2874000000000,
    "DownwardsTime": 2696000000000,
    "Duration": 23793000000000,
    "AverageSpeed": 6.754168746318261,
    "UpwardsSpeed": 6.213664008693707,
    "DownwardsSpeed": 7.3297513836181025,
    "AltitudeRange": 104.08999633789062
   }
  }
 ],
 "Summary": null
```

It is also possible to pipe in some file names instead of using the file names as input parameter

```sh
find ./testdata/valid-gpx -name "*.gpx" | ./bin/gpsa -summary=additional -out-file=./test.json
```

And you can pipe in file contents as well

```sh
cat  01.gpx 01.tcx 03.tcx 02.gpx | ./bin/gpsa -out-file=./test.json
```


#### Output Values explained

Below is a list of the output values and what they mean:

- `Name`: The name of the output line. Either read from the file, or if no name tag is set in the file it will be calculated out of the filename and the file hierarchy.
- `StartTime`: The time the track started. *not valid* in case we detect no or invalid time data.
- `EndTime`: The time the track ended. *not valid* in case we detect no or invalid time data.
- `TrackTime`: The time between `StartTime` and `EndTime`. *not valid* in case we detect no or invalid time data.
  - Formatted as `hh:mm:ss` in case of csv output
  - Measured in `ns` (nano seconds) in case of json output
- `Distance`: The distance of the track measured in `km`.
- `HorizontalDistance`: The horizontal distance of the track measured in `km`. This value ignores the vertical component of the distance like most GPS tools do.
- `AltitudeRange`: The range between the highest and the lowest point. Measured in `m`.
- `MinimumAltitude`: The altitude of the lowest point. Measured in `m`.
- `MaximumAltitude`: The altitude of the highest point. Measured in `m`.
- `ElevationGain`: The total sum of all upwards vertical distance.  Measured in `m`.
- `ElevationLose`: The total sum of all downwards vertical distance.  Measured in `m`.
- `UpwardsDistance`: The total sum of all distance moved upwards.  Measured in `km`.
- `DownwardsDistance`: The total sum of all distance moved downwards.  Measured in `km`.
- `MovingTime`: The time spend moving. *not valid* in case we detect no or invalid time data.
  - Formatted as `hh:mm:ss` in case of csv output
  - Measured in `ns` (nano seconds) in case of json output
- `UpwardsTime`: The time spend moving upwards. *not valid* in case we detect no or invalid time data.
  - Formatted as `hh:mm:ss` in case of csv output
  - Measured in `ns` (nano seconds) in case of json output
- `DownwardsTime`: The time spend downwards. *not valid* in case we detect no or invalid time data.
  - Formatted as `hh:mm:ss` in case of csv output
  - Measured in `ns` (nano seconds) in case of json output
- `AverageSpeed`: The average speed. Calculated from `Distance` and `MovingTime`.  *not valid* in case we detect no or invalid time data.
  - Measured in  `km/h`in case of csv output
  - Measured in  `m/s`in case of json output  
- `UpwardsSpeed`: The average speed during upwards movement . *not valid* in case we detect no or invalid time data.
  - Measured in  `km/h`in case of csv output
  - Measured in  `m/s`in case of json output  
- `DownwardsSpeed`: The average speed during downwards movement . *not valid* in case we detect no or invalid time data.
  - Measured in  `km/h`in case of csv output
  - Measured in  `m/s`in case of json output  

The statistic summary report will include `-` (in case of csv output) or `0.0000` (in case of json output) for statistic values that make no sense. For example it will not calculate a sum out of speed values or the average out of time stamps.

## Development

To develop this software install [Go v1.17](https://golang.org/) on your machine.

### Build

Use the [mage](https://magefile.org/) base build script to build and test the project.

```sh
./build.sh build        # build the project
./build.sh build test   # build and run the tests for the project
./build.sh test         # test the project
```

**Remark:** On Windows replace `./build.sh` with `.\build.bat`

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
            "command": "${workspaceRoot}/build.sh build test",
            "group": {
                "kind": "build",
                "isDefault": true
            }
        },
        {
            "label": "Test",
            "type": "shell",
            "command": "${workspaceRoot}/build.sh test",
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
