package linefilter

import "CLI-Geographic-Calculation/pkg/giocal/giocaltype"

type FilterByProperties[T giocaltype.GiotypeFeatureConstraint] func(feature *[]T, property string, value []string) []int
