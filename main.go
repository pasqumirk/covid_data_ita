package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type DayData struct {
	Data                      string `csv:"data" json:"data"`
	Stato                     string `csv:"stato" json:"stato"`
	RicoveratiConSintomi      int    `csv:"ricoverati_con_sintomi" json:"ricoverati_con_sintomi"`
	TerapiaIntensiva          int    `csv:"terapia_intensiva" json:"terapia_intensiva"`
	TotaleOspedalizzati       int    `csv:"totale_ospedalizzati" json:"totale_ospedalizzati"`
	IsolamentoDomiciliare     int    `csv:"isolamento_domiciliare" json:"isolamento_domiciliare"`
	TotalePositivi            int    `csv:"totale_positivi" json:"totale_positivi"`
	VariazioneTotalePositivi  int    `csv:"variazione_totale_positivi" json:"variazione_totale_positivi"`
	NuoviPositivi             int    `csv:"nuovi_positivi" json:"nuovi_positivi"`
	DimessiGuariti            int    `csv:"dimessi_guariti" json:"dimessi_guariti"`
	Deceduti                  int    `csv:"deceduti" json:"deceduti"`
	CasiDaSospettoDiagnostico int    `csv:"casi_da_sospetto_diagnostico" json:"casi_da_sospetto_diagnostico"`
	CasiDaScreening           int    `csv:"casi_da_screening" json:"casi_da_screening"`
	TotaleCasi                int    `csv:"totale_casi" json:"totale_casi"`
	Tamponi                   int    `csv:"tamponi" json:"tamponi"`
	CasiTestati               int    `csv:"casi_testati" json:"casi_testati"`
	Note                      string `csv:"note" json:"note"`
	IngressiTerapiaIntensiva  int    `csv:"ingressi_terapia_intensiva" json:"ingressi_terapia_intensiva"`
	NoteTest                  string `csv:"note_test" json:"note_test"`
	NoteCasi                  string `csv:"note_casi" json:"note_casi"`
}

func (data *DayData) Adjust() {
	if data.Data != "" {
		data.Data = data.Data[:10]
	}
}

func main() {
	os.Mkdir("out/", os.ModePerm)

	fmt.Println("Dati acquisiti da https://github.com/pcm-dpc/COVID-19/")

	now := time.Now()
	location := now.Location()

	// 1) Grab files
	day := time.Date(2020, 2, 24, 23, 59, 59, 0, location)
	for day.Before(now) {
		day = day.AddDate(0, 0, 1)
		file := fmt.Sprintf("dpc-covid19-ita-andamento-nazionale-%s.csv", day.Format("20060102"))
		outfile := "out/" + file
		if _, err := os.Stat(outfile); !os.IsNotExist(err) {
			continue
		}
		fmt.Printf("Processa %s\n", day.Format("2006-01-02"))

		url := fmt.Sprintf("https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-andamento-nazionale/%s", file)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Errore di rete: %s", err.Error())
			continue
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Errore di lettura da rete: %s", err.Error())
			continue
		}
		err = ioutil.WriteFile(outfile, body, 0644)
		if err != nil {
			fmt.Printf("Errore di apertura file: %s", err.Error())
			continue
		}
	}

	// 2) Process files
	allDays := []*DayData{}
	day = time.Date(2020, 2, 24, 23, 59, 59, 0, location)
	for day.Before(now) {
		day = day.AddDate(0, 0, 1)
		csvFileName := fmt.Sprintf("out/dpc-covid19-ita-andamento-nazionale-%s.csv", day.Format("20060102"))
		csvFile, err := os.Open(csvFileName)
		if err != nil {
			fmt.Printf("Errore di apertura file %s: %s", csvFileName, err.Error())
			continue
		}
		defer csvFile.Close()

		lines := []*DayData{}
		if err := gocsv.UnmarshalFile(csvFile, &lines); err != nil {
			fmt.Printf("Errore di lettura file %s: %s", csvFileName, err.Error())
			continue
		}
		csvFile.Close()
		allDays = append(allDays, lines...)
	}

	for _, day := range allDays {
		day.Adjust()
	}

	csvAllName := "out/dati-totali.csv"
	csvAll, err := os.Create(csvAllName)
	if err != nil {
		fmt.Printf("Errore di apertura file globale: %s", err.Error())
		os.Exit(1)
	}
	defer csvAll.Close()
	err = gocsv.MarshalFile(&allDays, csvAll)
	if err != nil {
		fmt.Printf("Errore di scrittura file globale: %s", err.Error())
		os.Exit(1)
	}

	json, err := json.Marshal(allDays)
	if err != nil {
		fmt.Printf("Errore di costruzione json: %s", err.Error())
		os.Exit(1)
	}
	ioutil.WriteFile("out/dati-totali.json", json, fs.ModePerm)

	os.Exit(0)
}
