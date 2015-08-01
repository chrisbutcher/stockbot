package bingclient

import (
  "encoding/json"
  "io/ioutil"
  "math/rand"
  "net/http"
  "time"
)

// Sandbox: https://datamarket.azure.com/dataset/explore/bing/search
var bingImagesBaseUrl string = "https://api.datamarket.azure.com/Bing/Search/v1/Image"

type BingResult struct {
  D struct {
    _Next   string `json:"__next"`
    Results []struct {
      ContentType string `json:"ContentType"`
      DisplayURL  string `json:"DisplayUrl"`
      FileSize    string `json:"FileSize"`
      Height      string `json:"Height"`
      ID          string `json:"ID"`
      MediaURL    string `json:"MediaUrl"`
      SourceURL   string `json:"SourceUrl"`
      Thumbnail   struct {
        ContentType string `json:"ContentType"`
        FileSize    string `json:"FileSize"`
        Height      string `json:"Height"`
        MediaURL    string `json:"MediaUrl"`
        Width       string `json:"Width"`
        _Metadata   struct {
          Type string `json:"type"`
        } `json:"__metadata"`
      } `json:"Thumbnail"`
      Title     string `json:"Title"`
      Width     string `json:"Width"`
      _Metadata struct {
        Type string `json:"type"`
        URI  string `json:"uri"`
      } `json:"__metadata"`
    } `json:"results"`
  } `json:"d"`
}

type BingClient struct {
  ApiToken string
  Market   string
  Adult    string
  Format   string
}

func (bc *BingClient) FetchImageRandomImage(searchTerms string) (string, error) {
  timeout := time.Duration(10 * time.Second)
  client := &http.Client{Timeout: timeout}

  req, err := bc.buildBingImageSearchRequest(searchTerms)
  resp, err := client.Do(req)

  defer resp.Body.Close()
  bytes, err := ioutil.ReadAll(resp.Body)

  var result BingResult
  if err := json.Unmarshal(bytes, &result); err != nil {
    panic(err)
  }

  if err != nil {
    return "Could not fetch image from bing", err
  }

  rand.Seed(time.Now().UTC().UnixNano())
  randomBingResult := result.D.Results[rand.Intn(len(result.D.Results))]
  thumbnailUrl := randomBingResult.Thumbnail.MediaURL

  return thumbnailUrl, err
}

func (bc *BingClient) buildBingImageSearchRequest(searchTerms string) (*http.Request, error) {
  req, err := http.NewRequest("GET", bingImagesBaseUrl, nil)

  values := req.URL.Query()
  values.Add("Query", "'stock photo "+searchTerms+"'")
  values.Add("Market", "'"+bc.Market+"'")
  values.Add("Adult", "'"+bc.Adult+"'")
  values.Add("$format", bc.Format)
  req.URL.RawQuery = values.Encode()

  req.Header.Set("Authorization", "Basic "+bc.ApiToken)

  return req, err
}
