package uuid

// UUID ...
type UUID struct {
	Data   string
	Config string
	Period string
}

var IRTemperature = UUID{
	Data:   "F000AA01-0451-4000-B000-000000000000",
	Config: "F000AA02-0451-4000-B000-000000000000",
	Period: "F000AA03-0451-4000-B000-000000000000",
}

var Humidity = UUID{
	Data:   "F000AA21-0451-4000-B000-000000000000",
	Config: "F000AA22-0451-4000-B000-000000000000",
	Period: "F000AA23-0451-4000-B000-000000000000",
}

var Barometer = UUID{
	Data:   "F000AA41-0451-4000-b000-000000000000",
	Config: "F000AA42-0451-4000-B000-000000000000",
	Period: "F000AA44-0451-4000-B000-000000000000",
}

var Optical = UUID{
	Data:   "F000AA71-0451-4000-B000-000000000000",
	Config: "F000AA72-0451-4000-B000-000000000000",
	Period: "F000AA73-0451-4000-B000-000000000000",
}

var IO = UUID{
	Data:   "F000AA65-0451-4000-b000-000000000000",
	Config: "F000AA66-0451-4000-B000-000000000000",
	Period: "", // i/o has no period
}

var Movement = UUID{
	Data:   "F000AA81-0451-4000-b000-000000000000",
	Config: "F000AA82-0451-4000-b000-000000000000",
	Period: "F000AA83-0451-4000-b000-000000000000",
}
