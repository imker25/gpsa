# gpsa - A GPX Statistic Analysing tool
This is a simple comandline tool, that may helps to extract statistic data out of *.gpx files

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
The Geografic calculations are done with the  ```haversine formula```  as descripted [here](http://www.movable-type.co.uk/scripts/latlong.html)
A Go implementaion of this can also be found [here](https://github.com/tkrajina/gpxgo/blob/master/gpx/geo.go). In this example the distance function igorens the atitute difference for distance bigger then 22 km, by checking the agular distance to be bigger then 0.2°.

My Impelemtaion will ignore atitute difference for distance bigger then 44 km, by checking the agular distance to be bigger then 0.3°.

# License

gpsa is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)
May got some testfiles from [gpxgo](https://github.com/tkrajina/gpxgo/tree/master/test_files)