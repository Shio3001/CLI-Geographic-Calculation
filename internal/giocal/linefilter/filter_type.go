package linefilter

import "CLI-Geographic-Calculation/internal/giocal/giocaltype"

type FilterByProperties func(railroadSections *[]giocaltype.GiotypeRailroadSection, property string, value []string) []int
