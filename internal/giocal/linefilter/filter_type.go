package linefilter

import "CLI-Geographic-Calculation/internal/giocal/giocaltype"

type FilterByProperties[T giocaltype.GiotypeFeatureConstraint] func(feature *[]T, property string, value []string) []int
