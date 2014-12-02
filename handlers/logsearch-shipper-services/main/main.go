package main

import (
  "fmt"
  "flag"
  "encoding/json"
  "io/ioutil"
  "net/http"
  "strings"
  "errors"
  "time"
  "regexp"
  "github.com/dpb587/logsearch-checks/check"
)

var flagElasticsearch string

func init() {
    flag.StringVar(&flagElasticsearch, "es", "localhost:9200", "Elasticsearch endpoint")
}

func main() {
  flag.Parse()

  s := strings.NewReader(ESDATA)

  resp, err := http.Post(
    fmt.Sprintf("http://%s/logstash-%s/metric/_search", flagElasticsearch, time.Now().Format("2006.01.02")),
    "application/json",
    s,
  )

  if err != nil {
    panic(err)
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)

  if err != nil {
    panic(err)
  }

  var esres map[string]interface{}

  err = json.Unmarshal(body, &esres)

  if err != nil {
    panic(err)
  }

  rem := regexp.MustCompile("monit.([^\\.]+).status")

  for _, director := range esres["aggregations"].(map[string]interface{})["filtered"].(map[string]interface{})["director"].(map[string]interface{})["buckets"].([]interface{}) {
    directormap := director.(map[string]interface{})

    for _, deployment := range directormap["deployment"].(map[string]interface{})["buckets"].([]interface{}) {
      deploymentmap := deployment.(map[string]interface{})

      for _, job := range deploymentmap["job"].(map[string]interface{})["buckets"].([]interface{}) {
        jobmap := job.(map[string]interface{})

        for _, job := range jobmap["service"].(map[string]interface{})["buckets"].([]interface{}) {
          metricmap := job.(map[string]interface{})

          value, err := extractTopHitValue(metricmap)

          if nil != err {
            panic(err)
          }

          owner := fmt.Sprintf("%s/%s/%s", directormap["key"], deploymentmap["key"], jobmap["key"])
          srvname := rem.ReplaceAllString(metricmap["key"].(string), "$1")

          cs := check.CHECK_FAIL

          if value == 0 {
            cs = check.CHECK_OKAY
          }

          mc := check.Status{
            Check : check.Check{
              Owner : owner,
              Name : fmt.Sprintf("monit_status/%s", srvname),
              Status : cs,
            },
            CheckData : check.CheckData{
              Threshold : 1,
              Value : value,
              Units : "boolean",
            },
          }

          mcjs, err := json.Marshal(mc)

          if nil != err {
            panic(err)
          }

          fmt.Println(string(mcjs))
        }
      }
    }
  }
}

func extractTopHitValue(d interface{}) (val float64, err error) {
  sd := d.(map[string]interface{})["value"].(map[string]interface{})["hits"].(map[string]interface{})["hits"].([]interface{})

  if (0 == len(sd)) {
    err = errors.New("No top hit available")

    return
  }

  val = sd[0].(map[string]interface{})["_source"].(map[string]interface{})["value"].(float64)

  return
}

const ESDATA string = `
  {
    "size" : 0,
    "aggs" : {
      "filtered" : {
        "filter" :  {
          "regexp" : {
            "name" : "monit\\..*\\.status"
          }
        },
        "aggs" : {
          "director" : {
            "terms" : {
              "field" : "@source.bosh_director",
              "size" : 250
            },
            "aggs" : {
              "deployment" : {
                "terms" : {
                  "field" : "@source.bosh_deployment",
                  "size" : 250
                },
                "aggs" : {
                  "job" : {
                    "terms" : {
                      "field" : "@source.bosh_job",
                      "size" : 250
                    },
                    "aggs" : {
                      "service" : {
                        "terms" : {
                          "field" : "name",
                          "size" : 100
                        },
                        "aggs" : {
                          "value" : {
                            "top_hits" : {
                              "sort" : {
                                "@timestamp" : {
                                  "order" : "desc"
                                },
                                "name" : {
                                  "order" : "desc"
                                }
                              },
                              "_source" : {
                                "include" : [
                                  "value"
                                ]
                              },
                              "size" : 1
                            }
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
`
