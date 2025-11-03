package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Moyskladapiconfig
	RefGoconfig
}

type Moyskladapiconfig struct {
	APIKEY            string
	TimeSpan          time.Duration
	RequestCap        int
	Readystatehref    string
	Shipedstatehref   string
	SellTypehref      string
	SellTypeID        string
	SellTypeOtherhref string
	Storehref         string
	Orghref           string
	RefGoNumberhref   string
	RefGoNumberID     string
	Courierhref       string
	CouierID          string
	RefGoCourierhref  string
	TimeFormat        string
	URLstart          string
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

	selltypehref := os.Getenv("MSAPI_SELLTYPEHREF")
	if selltypehref == "" {
		panic("Selltype does not exist")
	}

	selltypeID := os.Getenv("MSAPI_SELLTYPEID")
	if selltypeID == "" {
		panic("SelltypeID does not exist")
	}

	selltypeOtherhref := os.Getenv("MSAPI_SELLTYPEOTHERHREF")
	if selltypeOtherhref == "" {
		panic("OtherSelltype does not exist")
	}

	storehref := os.Getenv("MSAPI_STOREHREF")
	if storehref == "" {
		panic("Storehref does not exist")
	}

	orghref := os.Getenv("MSAPI_ORGHREF")
	if orghref == "" {
		panic("Orghref does not exist")
	}

	refgonumberhref := os.Getenv("MSAPI_REFGONUMBERHREF")
	if refgonumberhref == "" {
		panic("RefGoNumberhref does not exist")
	}

	refgonumberid := os.Getenv("MSAPI_REFGONUMBERID")
	if refgonumberid == "" {
		panic("RefGoNumberID does not exist")
	}

	courierhref := os.Getenv("MSAPI_COURIERHREF")
	if courierhref == "" {
		panic("Courierhref does not exist")
	}

	courierid := os.Getenv("MSAPI_COURIERID")
	if courierid == "" {
		panic("CourierID does not exist")
	}

	refgocourierhref := os.Getenv("MSAPI_REFGOCOURIERHREF")
	if refgocourierhref == "" {
		panic("RefGoCourierhref does not exist")
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

	latestorder, err := strconv.Atoi(strings.Trim(os.Getenv("RG_LATESTORDER"), `"`))
	if err != nil {
		panic("Invalid RG_LATESTORDER")
	}
	fmt.Println(latestorder)

	return &Config{
		Moyskladapiconfig{
			APIKEY:            apiKey,
			TimeSpan:          tspn,
			RequestCap:        rqcap,
			Readystatehref:    readystatehref,
			Shipedstatehref:   shipedstatehref,
			SellTypehref:      selltypehref,
			SellTypeID:        selltypeID,
			SellTypeOtherhref: selltypeOtherhref,
			Storehref:         storehref,
			Orghref:           orghref,
			RefGoNumberhref:   refgonumberhref,
			RefGoNumberID:     refgonumberid,
			Courierhref:       courierhref,
			CouierID:          courierid,
			RefGoCourierhref:  refgocourierhref,
			TimeFormat:        timeFormat,
			URLstart:          urlstart,
		},
		RefGoconfig{
			RGLatestOrder: latestorder,
		},
	}
}

func ChangeRefGoLatest(latestOrder int) error {
	envFile := "../.env"
	content, err := os.ReadFile(envFile)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	found := false

	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "RG_LATESTORDER=") {
			lines[i] = fmt.Sprintf("RG_LATESTORDER=\"%d\"", latestOrder)
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("RG_LATESTORDER=\"%d\"", latestOrder))
	}

	if err := os.WriteFile(envFile, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}

	return nil
}
