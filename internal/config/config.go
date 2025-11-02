package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Moyskladapiconfig
	RefGoconfig
}

type Moyskladapiconfig struct {
	APIKEY          string
	TimeSpan        time.Duration
	RequestCap      int
	Readystatehref  string
	Shipedstatehref string
	Storehref       string
	Orghref         string
	TimeFormat      string
	URLstart        string
}

type RefGoconfig struct {
	RGLatestOrder int
}

func LoadConfig() *Config {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Cannot read config file")
	}

	apiKey := os.Getenv("MSAPI_KEY")
	if apiKey == "" {
		panic("API KEY does not exist")
	}

	tspnint, err := strconv.Atoi(os.Getenv("MSAPI_REQUESTCAPTIMESPAN"))
	if err != nil {
		panic("Timespan does not exist")
	}

	tspn := time.Duration(int64(tspnint)) * time.Second

	rqcap, err := strconv.Atoi(os.Getenv("MSAPI_REQUESTCAP"))
	if err != nil {
		panic("Requestcap does not exist")
	}

	readystatehref := os.Getenv("MSAPI_READYSTATEHREF")
	if readystatehref == "" {
		panic("Statehref does not exist")
	}

	shipedstatehref := os.Getenv("MSAPI_SHIPEDSTATEHREF")
	if shipedstatehref == "" {
		panic("Shipedstatehref does not exist")
	}

	storehref := os.Getenv("MSAPI_STOREHREF")
	if storehref == "" {
		panic("Storehref does not exist")
	}

	orghref := os.Getenv("MSAPI_ORGHREF")
	if orghref == "" {
		panic("Orghref does not exist")
	}

	timeFormat := os.Getenv("MSAPI_TIMEFORMAT")
	if timeFormat == "" {
		panic("Timeformat does not exist")
	}

	urlstart := os.Getenv("MSAPI_URLSTART")
	if urlstart == "" {
		panic("URLstart does not exist")
	}

	if os.Getenv("RG_LATESTORDER") == "" {
		panic("RG_LATESTORDER does not exist")
	}
	latestorder, err := strconv.Atoi(os.Getenv("RG_LATESTORDER"))
	if err != nil {
		panic("Invalid RG_LATESTORDER")
	}

	return &Config{
		Moyskladapiconfig{
			APIKEY:          apiKey,
			TimeSpan:        tspn,
			RequestCap:      rqcap,
			Readystatehref:  readystatehref,
			Shipedstatehref: shipedstatehref,
			Storehref:       storehref,
			Orghref:         orghref,
			TimeFormat:      timeFormat,
			URLstart:        urlstart,
		},
		RefGoconfig{
			RGLatestOrder: latestorder,
		},
	}
}
