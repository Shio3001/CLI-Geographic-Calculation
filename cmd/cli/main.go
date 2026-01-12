package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"CLI-Geographic-Calculation/internal/giocal"
	"CLI-Geographic-Calculation/internal/giocal/giocaltype"
	giocal_load "CLI-Geographic-Calculation/internal/giocal/load"
)

//現在のミリ秒取得
func currentMillis() int64 {
	return int64(float64(time.Now().UnixNano()) / 1e6)
}

func main() {
	var (
		stationPath    = flag.String("station", "", "path to station GeoJSON (N02 station)")
		sectionPath    = flag.String("section", "", "path to railroad section GeoJSON (N02 railroad section)")
		passengersPath = flag.String("passengers", "", "path to passengers GeoJSON (S12 passengers)")
		history = flag.String("history", "", "path to railroad history GeoJSON (N05 railroad history)")

		company  = flag.String("company", "", "filter by company name (e.g. \"東日本旅客鉄道\")")
		linesRaw = flag.String("lines", "", "comma-separated target line names (e.g. \"山手線,中央線\")")

		pretty = flag.Bool("pretty", false, "pretty-print JSON output")
		out    = flag.String("out", "", "output file path (optional). if empty, print to stdout")
	)
	flag.Parse()

	if *stationPath == "" || *sectionPath == "" || *passengersPath == "" || *history == "" {
		fmt.Fprintln(os.Stderr, "ERROR: -station, -section, -passengers and -history are required")
		flag.Usage()
		os.Exit(2)
	}

	

	//実行開始時間出力(ミリ秒単位)
	startTime := currentMillis()
	fmt.Fprintf(os.Stderr, "Start time: %d ms\n",startTime)

	targetLines := parseCSV(*linesRaw)
	targetCompany := strings.TrimSpace(*company)

	var (
		stFC         *giocaltype.GiotypeStationFeatureCollection
		rrFC         *giocaltype.GiotypeRailroadSectionFeatureCollection
		passengersFC *giocaltype.GiotypePassengersFeatureCollection
		historyFC    *giocaltype.GiotypeN05RailroadSectionFeatureCollection
		err          error
	)

	// ---- Load stations / sections with optional filters ----
	if len(targetLines) > 0 || targetCompany != "" {
		stFC, err = giocal_load.LoadGiotypeStationForCompanyAndLines(*stationPath, targetCompany, targetLines)
		if err != nil {
			die(err)
		}

		rrFC, err = giocal_load.LoadGiotypeRailroadSectionForCompanyAndLines(*sectionPath, targetCompany, targetLines)
		if err != nil {
			die(err)
		}
	} else {
		stFC, err = giocal.LoadGiotypeStation(*stationPath)
		if err != nil {
			die(err)
		}

		rrFC, err = giocal.LoadGiotypeRailroadSection(*sectionPath)
		if err != nil {
			die(err)
		}
	}

	passengersFC, err = giocal_load.LoadGiotypePassengersForCompanyAndLines(*passengersPath , targetCompany, targetLines)
	if err != nil {
		die(err)
	}
	historyFC, err = giocal_load.LoadGiotypeRailHistoryForCompanyAndLines(*history, targetCompany, targetLines)
	if err != nil {
		die(err)
	}

	g := giocal.ConvertGiotypeStationToGraph(stFC, rrFC, passengersFC, historyFC)

	convertTime := currentMillis()
	fmt.Fprintf(os.Stderr, "Conversion time: %d ms\n", convertTime - startTime)



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

	//実行終了時間出力(ミリ秒単位)
	endTime := currentMillis()
	
	fmt.Fprintf(os.Stderr, "End time: %d ms\n",endTime)

	//計測時刻出力
	fmt.Fprintf(os.Stderr, "Total execution time: %d ms\n", endTime - startTime)
	fmt.Fprintf(os.Stderr, "Graph conversion time: %d ms\n", convertTime - startTime)
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
