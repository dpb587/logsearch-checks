Example

    ( \
      go run handlers/logsearch-shipper-services/main/main.go -es=api.meta.logsearch.io:9200 ; \
      go run handlers/logsearch-shipper-diskusage/main/main.go -es=api.meta.logsearch.io:9200 -persistent=70 \
    ) \
      | go run handlers/drop-okay-checks/main/main.go \
      | go run handlers/notify-via-email/main/main.go
