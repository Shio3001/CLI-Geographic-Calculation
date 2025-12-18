package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Shio3001/CLI-Geographic-Calculation/internal/giocal"
	"github.com/Shio3001/CLI-Geographic-Calculation/internal/giocal/giocaltype"
)

func main() {
	var (
		stationPath = flag.String("station", "", "path to station GeoJSON (N02 station)")
		sectionPath = flag.String("section", "", "path to railroad section GeoJSON (N02 railroad section)")
		linesRaw    = flag.String("lines", "", "comma-separated target line names (e.g. \"山手線,中央線\")")
		pretty      = flag.Bool("pretty", false, "pretty-print JSON output")
		out         = flag.String("out", "", "output file path (optional). if empty, print to stdout")
	)
	flag.Parse()

	if *stationPath == "" || *sectionPath == "" {
		fmt.Fprintln(os.Stderr, "ERROR: -station and -section are required")
		flag.Usage()
		os.Exit(2)
	}

	targetLines := parseCSV(*linesRaw)

	var (
		stFC *giocaltype.GiotypeStationFeatureCollection
		rrFC *giocaltype.GiotypeRailroadSectionFeatureCollection
		err  error
	)

	if len(targetLines) > 0 {
		stFC, err = giocal.ReadGiotypeStationForLines(*stationPath, targetLines)
		if err != nil {
			die(err)
		}
		rrFC, err = giocal.ReadGiotypeRailroadSectionForLines(*sectionPath, targetLines)
		if err != nil {
			die(err)
		}
	} else {
		stFC, err = giocal.ReadGiotypeStation(*stationPath)
		if err != nil {
			die(err)
		}
		rrFC, err = giocal.ReadGiotypeRailroadSection(*sectionPath)
		if err != nil {
			die(err)
		}
	}

	g := giocal.ConvertGiotypeStationToGraph(stFC, rrFC)

	var b []byte
	if *pretty {
		b, err = json.MarshalIndent(g, "", "  ")
	} else {
		b, err = json.Marshal(g)
	}
	if err != nil {
		die(err)
	}

	if *out == "" {
		fmt.Println(string(b))
		return
	}

	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		die(err)
	}
	if err := os.WriteFile(*out, b, 0o644); err != nil {
		die(err)
	}
	fmt.Fprintf(os.Stderr, "Wrote: %s\n", *out)
}

func parseCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "ERROR:", err)
	os.Exit(1)
}


