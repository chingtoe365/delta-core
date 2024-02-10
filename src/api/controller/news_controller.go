package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/robtec/newsapi/api"
)

type NewsController struct {
}

type NewsQueryParams struct {
	Keywords string `form:"keywords" binding:"required"`
	Start    string `form:"start" default:"2023-02-08T00:00:00Z"`
	End      string `form:"end" default:"2023-02-09T00:00:00Z"`
}

type NewsArtical struct {
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PublishedAt string `json:"publishedAt"`
}

// GetQuote godoc
// @Summary get news
// @Schemes
// @Description get news via keywords and timespan
// @Tags News
// @Param keywords query string true "Keywords"
// @Param start query string false "Start time"
// @Param end query string false "End time"
// @Accept json
// @Produce json
// @Success 200
// @Router /get-news [get]
func (nc *NewsController) Get(c *gin.Context) {
	var params NewsQueryParams
	var result []NewsArtical
	if c.ShouldBindQuery(&params) == nil {
		log.Println(params)
		log.Println("Error in parsing query parameters")
		c.JSON(http.StatusBadRequest, result)
	}
	// keywords := c.Query("keywords")
	// stas
	query := params.Keywords

	httpClient := http.Client{}
	key := "55669d67fa794458b77130bb3b71aeb4"
	url := "https://newsapi.org"

	// Create a client, passing in the above
	client, err := api.New(&httpClient, key, url)

	if err != nil {
		c.JSON(http.StatusInternalServerError, "News search error")
	}

	// Create options for Ireland and Business
	// opts := api.Options{Country: "ie", Category: "business"}

	// Get Top Headlines with options from above
	// topHeadlines, err := client.TopHeadlines(opts)

	// Different options
	moreOpts := api.Options{
		Language: "en",
		Q:        query,
		SortBy:   "popularity",
		From:     params.Start,
		To:       params.End,
	}

	// Get Everything with options from above
	everything, err := client.Everything(moreOpts)

	if err != nil {
		log.Fatalf("Something wrong when fetching news")
	}

	for _, artical := range everything.News.Articles {
		a := NewsArtical{
			Author:      artical.Author,
			Title:       artical.Title,
			Description: artical.Description,
			URL:         artical.URL,
			PublishedAt: artical.PublishedAt,
		}
		result = append(result, a)
	}

	c.JSON(http.StatusOK, result)
}
