module "tobi.backfrak.de/gpsa"



require "tobi.backfrak.de/internal/gpsabl" v0.0.0
replace  "tobi.backfrak.de/internal/gpsabl" v0.0.0 => "../../internal/gpsabl"
require "tobi.backfrak.de/internal/csvbl" v0.0.0
replace  "tobi.backfrak.de/internal/csvbl" v0.0.0 => "../../internal/csvbl"
require "tobi.backfrak.de/internal/gpxbl" v0.0.0
replace  "tobi.backfrak.de/internal/gpxbl" v0.0.0 => "../../internal/gpxbl"
require "tobi.backfrak.de/internal/jsonbl" v0.0.0
replace  "tobi.backfrak.de/internal/jsonbl" v0.0.0 => "../../internal/jsonbl"
require "tobi.backfrak.de/internal/tcxbl" v0.0.0
replace  "tobi.backfrak.de/internal/tcxbl" v0.0.0 => "../../internal/tcxbl"

require "tobi.backfrak.de/internal/testhelper" v0.0.0
replace  "tobi.backfrak.de/internal/testhelper" v0.0.0 => "../../internal/testhelper"


