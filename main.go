package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type MetarRequest struct {
	Raw string `json:"raw"`
}

type MetarResponse struct {
	Readable string `json:"readable"`
}

func parseMetar(metar string) string {
	r := regexp.MustCompile(`(?i)^(?P<airport>[A-Z]{4})\s*(?P<time>\d{6}Z)\s*(?P<winddir>\d{3})(?P<windspeed>\d{2})(?:G(?P<gusts>\d{2}))?KT\s*(?P<visibility>\d{4})\s*(?P<clouds>.*)\s*(?P<temperature>\d{2})/(?P<dewpoint>\d{2})\s*(?P<pressure>\d{4})(?P<tempo>.*)?$`)
	matches := r.FindStringSubmatch(metar)

	if matches == nil {
		return "METAR解析に失敗しました"
	}

	airport := matches[1]
	time := matches[2]
	winddir := matches[3]
	windspeed := matches[4]
	gusts := matches[5]
	visibility := matches[6]
	clouds := matches[7]
	temperature := matches[8]
	dewpoint := matches[9]
	pressure := matches[10]
	tempo := matches[11]

	cloudInfo := "雲なし"
	if clouds != "" {
		cloudInfo = clouds
	}

	windGustInfo := ""
	if gusts != "" {
		windGustInfo = "突風: " + gusts + "KT"
	}

	if strings.HasPrefix(temperature, "M") {
		temperature = temperature[1:] + "℃"
	} else {
		temperature += "℃"
	}

	if strings.HasPrefix(dewpoint, "M") {
		dewpoint = dewpoint[1:] + "℃"
	} else {
		dewpoint += "℃"
	}

	tempoInfo := ""
	if tempo != "" {
		tempoInfo = "一時的な天候: " + tempo
	}

	return strings.Join([]string{
		"空港コード: " + airport,
		"観測時刻(UTC): " + time,
		"風向: " + winddir + "度",
		"風速: " + windspeed + "KT",
		windGustInfo,
		"視程: " + visibility + "メートル",
		"雲情報: " + cloudInfo,
		"気温: " + temperature,
		"露点温度: " + dewpoint,
		"気圧: " + pressure + "hPa",
		tempoInfo,
	}, "\n")
}

func metarHandler(ctx *gin.Context) {
	var request MetarRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(request.Raw)
	result := parseMetar(request.Raw)
	ctx.JSON(http.StatusOK, MetarResponse{Readable: result})
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
