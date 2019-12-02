module github.com/ardacetinkaya/FirstGO

require (
	github.com/Azure/azure-storage-queue-go v0.0.0-20191125232315-636801874cdd
	github.com/ardacetinkaya/FirstGO/config v0.0.0-00010101000000-000000000000
	github.com/ardacetinkaya/FirstGO/token v0.0.0-00010101000000-000000000000
)

replace github.com/ardacetinkaya/FirstGO/config => ./config

replace github.com/ardacetinkaya/FirstGO/token => ./token

go 1.13
