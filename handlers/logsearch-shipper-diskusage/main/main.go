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
  "strconv"
  "github.com/dpb587/logsearch-checks/check"
  "github.com/dpb587/logsearch-checks/handlers/logsearch-shipper-diskusage"
)

var flagElasticsearch string
var flagEphemeral string
var flagPersistent string
var flagSystem string

func init() {
    flag.StringVar(&flagElasticsearch, "es", "localhost:9200", "Elasticsearch endpoint")
    flag.StringVar(&flagEphemeral, "ephemeral", "90.00", "Ephemeral Usage Threshold")
    flag.StringVar(&flagSystem, "system", "90.00", "Systemt Usage Threshold")
    flag.StringVar(&flagPersistent, "persistent", "80.00", "Persistent Usage Threshold")
}

func main() {
  flag.Parse()

  // ignoring err :(
  flagEphemeralFloat, _ := strconv.ParseFloat(flagEphemeral, 64)
  flagSystemFloat, _ := strconv.ParseFloat(flagSystem, 64)
  flagPersistentFloat, _ := strconv.ParseFloat(flagPersistent, 64)

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

  for _, director := range esres["aggregations"].(map[string]interface{})["director"].(map[string]interface{})["buckets"].([]interface{}) {
    directormap := director.(map[string]interface{})

    for _, deployment := range directormap["deployment"].(map[string]interface{})["buckets"].([]interface{}) {
      deploymentmap := deployment.(map[string]interface{})

      for _, job := range deploymentmap["job"].(map[string]interface{})["buckets"].([]interface{}) {
        jobmap := job.(map[string]interface{})

        owner := fmt.Sprintf("%s/%s/%s", directormap["key"], deploymentmap["key"], jobmap["key"])

        generateCheckStatust(owner, "system", extractDiskUsage(jobmap, "system"), flagSystemFloat)
        generateCheckStatust(owner, "ephemeral", extractDiskUsage(jobmap, "ephemeral"), flagEphemeralFloat)
        generateCheckStatust(owner, "persistent", extractDiskUsage(jobmap, "persistent"), flagPersistentFloat)
      }
    }
  }
}

func generateCheckStatust(owner string, disk string, du logsearchshipperdiskusage.DiskUsage, threshold float64) {
  if du.IsMissingData() {
    return
  }

  value := du.GetUsedPct()
  cs := check.CHECK_FAIL

  if value < threshold {
    cs = check.CHECK_OKAY
  }

  mc := check.Status{
    Check : check.Check{
      Owner : owner,
      Name : fmt.Sprintf("disk_usage/%s", disk),
      Status : cs,
    },
    CheckData : check.CheckData{
      Threshold : threshold,
      Value : value,
      Units : "percent",
      Extra : map[string]float64{
        "free_bytes" : du.GetFree(),
        "used_bytes" : du.GetUsed(),
      },
    },
  }

  mcjs, err := json.Marshal(mc)

  if nil != err {
    panic(err)
  }

  fmt.Println(string(mcjs))
}

func extractDiskUsage(d map[string]interface{}, name string) (du logsearchshipperdiskusage.DiskUsage) {
  used, err1 := extractTopHitValue(d[name + "_used"])
  free, err2 := extractTopHitValue(d[name + "_free"])

  return logsearchshipperdiskusage.DiskUsage{ err1 != nil || err2 != nil, used, free }
}

func extractTopHitValue(d interface{}) (val float64, err error) {
  sd := d.(map[string]interface{})["value"].(map[string]interface{})["hits"].(map[string]interface{})["hits"].([]interface{})

  if 0 == len(sd) {
    err = errors.New("No top hit available")

    return
  }

  val = sd[0].(map[string]interface{})["_source"].(map[string]interface{})["value"].(float64)

  return
}

const ESDATA string = `
    {
      "aggregations" : {
        "director" : {
          "terms" : {
            "field" : "@source.bosh_director",
            "size" : 0,
            "order" : {
              "_term" : "asc"
            }
          },
          "aggregations" : {
            "deployment" : {
              "terms" : {
                "field" : "@source.bosh_deployment",
                "size" : 0,
                "order" : {
                  "_term" : "asc"
                }
              },
              "aggregations" : {
                "job" : {
                  "terms" : {
                    "field" : "@source.bosh_job",
                    "size" : 0,
                    "order" : {
                      "_term" : "asc"
                    }
                  },
                  "aggregations" : {
                    "system_free" : {
                      "filter" : {
                        "term" : {
                          "name" : "host.df_xvda1.df_complex_free"
                        }
                      },
                      "aggregations" : {
                        "value" : {
                          "top_hits" : {
                            "sort" : {
                              "@timestamp" : "desc"
                            },
                            "size" : 1,
                            "_source" : [
                              "value"
                            ]
                          }
                        }
                      }
                    },
                    "system_used" : {
                      "filter" : {
                        "term" : {
                          "name" : "host.df_xvda1.df_complex_used"
                        }
                      },
                      "aggregations" : {
                        "value" : {
                          "top_hits" : {
                            "sort" : {
                              "@timestamp" : "desc"
                            },
                            "size" : 1,
                            "_source" : [
                              "value"
                            ]
                          }
                        }
                      }
                    },
                    "ephemeral_free" : {
                      "filter" : {
                        "term" : {
                          "name" : "host.df_xvdb2.df_complex_free"
                        }
                      },
                      "aggregations" : {
                        "value" : {
                          "top_hits" : {
                            "sort" : {
                              "@timestamp" : "desc"
                            },
                            "size" : 1,
                            "_source" : [
                              "value"
                            ]
                          }
                        }
                      }
                    },
                    "ephemeral_used" : {
                      "filter" : {
                        "term" : {
                          "name" : "host.df_xvdb2.df_complex_used"
                        }
                      },
                      "aggregations" : {
                        "value" : {
                          "top_hits" : {
                            "sort" : {
                              "@timestamp" : "desc"
                            },
                            "size" : 1,
                            "_source" : [
                              "value"
                            ]
                          }
                        }
                      }
                    },
                    "persistent_free" : {
                      "filter" : {
                        "term" : {
                          "name" : "host.df_xvdf1.df_complex_free"
                        }
                      },
                      "aggregations" : {
                        "value" : {
                          "top_hits" : {
                            "sort" : {
                              "@timestamp" : "desc"
                            },
                            "size" : 1,
                            "_source" : [
                              "value"
                            ]
                          }
                        }
                      }
                    },
                    "persistent_used" : {
                      "filter" : {
                        "term" : {
                          "name" : "host.df_xvdf1.df_complex_used"
                        }
                      },
                      "aggregations" : {
                        "value" : {
                          "top_hits" : {
                            "sort" : {
                              "@timestamp" : "desc"
                            },
                            "size" : 1,
                            "_source" : [
                              "value"
                            ]
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
      },
      "size" : 0
    }
`
