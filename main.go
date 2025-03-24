package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type MetarRequest struct {
	Raw string `json:"raw"`
}

// 構造化されたレスポンスの定義
type MetarResponse struct {
	Airport       string      `json:"airport"`
	Time          string      `json:"time"`
	WindDirection string      `json:"windDirection"`
	WindSpeed     string      `json:"windSpeed"`
	Visibility    string      `json:"visibility"`
	Clouds        []CloudInfo `json:"clouds"`
	Temperature   string      `json:"temperature"`
	DewPoint      string      `json:"dewPoint"`
	Pressure      string      `json:"pressure"`
	TempoInfo     string      `json:"tempoInfo"`
	Remarks       string      `json:"remarks"`
}

type CloudInfo struct {
	Type   string `json:"type"`
	Height string `json:"height"`
}

func parseMetar(metar string) (MetarResponse, error) {
	response := MetarResponse{
		Clouds: []CloudInfo{},
	}

	basicInfo := regexp.MustCompile(`([A-Z]{4})\s+(\d{2})(\d{2})(\d{2})Z`)
	basicMatches := basicInfo.FindStringSubmatch(metar)

	if basicMatches == nil {
		return response, fmt.Errorf("METAR解析に失敗しました：基本情報の取得ができません")
	}

	response.Airport = basicMatches[1]

	hour := basicMatches[3]
	minute := basicMatches[4]
	response.Time = hour + ":" + minute

	windInfo := regexp.MustCompile(`(\d{3})(\d{2})KT`)
	windMatches := windInfo.FindStringSubmatch(metar)

	if windMatches != nil {
		response.WindDirection = windMatches[1]
		response.WindSpeed = windMatches[2]
	}

	// 視程を抽出
	visibilityInfo := regexp.MustCompile(`KT\s+(9999|\d{4})`)
	visibilityMatches := visibilityInfo.FindStringSubmatch(metar)

	if visibilityMatches != nil {
		if visibilityMatches[1] == "9999" {
			response.Visibility = "10km以上"
		} else {
			response.Visibility = visibilityMatches[1] + "メートル"
		}
	}

	// 雲情報を抽出
	cloudInfo := regexp.MustCompile(`(FEW|SCT|BKN|OVC)(\d{3}|///)`)
	cloudMatches := cloudInfo.FindAllStringSubmatch(metar, -1)

	for _, cloud := range cloudMatches {
		cloudType := cloud[1]
		cloudHeight := cloud[2]

		heightDesc := cloudHeight
		if cloudHeight == "///" {
			heightDesc = "不明"
		} else {
			heightDesc += "フィート"
		}

		response.Clouds = append(response.Clouds, CloudInfo{
			Type:   cloudType,
			Height: heightDesc,
		})
	}

	// 気温と露点を抽出
	tempInfo := regexp.MustCompile(`\s+(M?\d{1,2})/(M?\d{1,2})\s+`)
	tempMatches := tempInfo.FindStringSubmatch(metar)

	if tempMatches != nil {
		temperature := tempMatches[1]
		dewpoint := tempMatches[2]

		if strings.HasPrefix(temperature, "M") {
			temperature = "-" + temperature[1:]
		}
		if strings.HasPrefix(dewpoint, "M") {
			dewpoint = "-" + dewpoint[1:]
		}

		response.Temperature = temperature + "℃"
		response.DewPoint = dewpoint + "℃"
	}

	// 気圧を抽出
	pressureInfo := regexp.MustCompile(`\s+(Q\d{4})`)
	pressureMatches := pressureInfo.FindStringSubmatch(metar)

	if pressureMatches != nil {
		response.Pressure = pressureMatches[1]
	}

	// 一時的な天候情報を抽出
	if tempoParts := regexp.MustCompile(`TEMPO\s+(.+?)(?:\s+RMK|\s*$)`).FindStringSubmatch(metar); tempoParts != nil {
		response.TempoInfo = tempoParts[1]
	}

	// 備考情報を抽出
	if rmkParts := regexp.MustCompile(`RMK\s+(.+?)(?:\s*$)`).FindStringSubmatch(metar); rmkParts != nil {
		response.Remarks = rmkParts[1]
	}

	return response, nil
}

func metarHandler(ctx *gin.Context) {
	var request MetarRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(request.Raw)
	result, err := parseMetar(request.Raw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("src/*.html")
	router.Static("/js", "src/js")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.POST("/api/metar", metarHandler)

	router.Run()
}
