# gpsa - A GPX Statistic Analysing tool
This is a simple comandline tool, that may helps to extract statistic data out of *.gpx files

# Usage
You might want to call ```-help```, to find out how to use the program.

```sh
./gpsa -help                                                     
./gpsa: Reads in GPS track files, and writes out basic statistic data found in the track

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

## Internals
The Geografic calculations are done with the  ```haversine formula```  as descripted [here](http://www.movable-type.co.uk/scripts/latlong.html. My Impelemtaion will ignore atitute difference for distance bigger then 33 km, by checking the agular distance to be bigger then 0.3°. For smaller distances the program will add the atitute difference using the ```pythagoras theorem```.

# License

gpsa is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)
May got some testfiles from [gpxgo](https://github.com/tkrajina/gpxgo/tree/master/test_files)