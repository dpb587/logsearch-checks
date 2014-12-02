package main

import (
    "bytes"
    "flag"
    "os"
    "fmt"
    "regexp"
    "text/template"
    "github.com/dpb587/logsearch-checks/check"
    "github.com/dpb587/logsearch-checks/handler"
    //"net/smtp"
)

var notifyFrom string
var notifyTo string

func init() {
    flag.StringVar(&notifyFrom, "from", "nobody@example.com", "From Address")
    flag.StringVar(&notifyTo, "to", "somebody@example.com", "To Address")
}

func mhandle(cs check.Status, emit handler.Emitter) {
    var body bytes.Buffer
    var tmplbody string
    var tmplbodybuf bytes.Buffer

    body.WriteString(fmt.Sprintf("From: %s\r\n", notifyFrom))
    body.WriteString(fmt.Sprintf("To: %s\r\n", notifyTo))
    
    tmpl, err := template.New("test").Parse(EMAIL_TXT)

    if nil != err {
        panic(err)
    }

    err = tmpl.Execute(&tmplbodybuf, cs)

    if nil != err {
        panic(err)
    }

    tmplbody = tmplbodybuf.String()

    resub := regexp.MustCompile("(?m)^([^\n]+)\n(.*)$")

    body.WriteString(fmt.Sprintf("Subject: %s", resub.ReplaceAllString(tmplbody, "$1\r\n\r\n$2")))

    fmt.Printf(body.String())
    fmt.Printf("\n")

    // err := smtp.SendMail(
    //     fmt.Sprintf("%s:%s", os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT")),
    //     smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST")),
    //     notifyFrom,
    //     []string{notifyTo},
    //     []byte(body)
    // )

    // if err != nil {
    //     panic(err)
    // }

}

func main() {
    flag.Parse()

    v := handler.Handler{os.Stdin, mhandle}
    v.Run()
}

const EMAIL_TXT string = `{{ .Check.Status }}: {{ .Check.Name }}
{{ .Check.Owner }}
Actual is {{ .CheckData.Value }}, threshold is {{ .CheckData.Threshold }}
`
