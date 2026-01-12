package giocaltype

type DatasetResourcePath struct {
	Rail string
	Station string
}

type DatasetResource struct {
	Rail *GiotypeRailroadSectionFeatureCollection
	Station *GiotypeStationFeatureCollection
}
