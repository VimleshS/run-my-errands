package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"github.com/VimleshS/run-my-errands/models"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	log "github.com/Sirupsen/logrus"

	"github.com/VimleshS/run-my-errands/setup"
	"github.com/bgentry/que-go"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

var (
	qc      *que.Client
	pgxpool *pgx.ConnPool
)

func main() {
	var err error
	pgxpool, qc, err = setup.SetUp.PoolAndQueueConnection()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer pgxpool.Close()

	log.Info("Starting dispatcher ...")

	wm := que.WorkMap{
		models.IndexRequestJob: indexURLJob,
	}

	// 2 worker go routines
	workers := que.NewWorkerPool(qc, wm, 2)

	// Catch signal so we can shutdown gracefully
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go workers.Start()

	// Wait for a signal
	sig := <-sigCh
	log.WithField("signal", sig).Info("Signal received. Shutting down.")

	workers.Shutdown()
}

const mailtemplate = `
Hi {{.Email}},

   Thankyou for shopping on RunAErrand, for {{.Total}}
   
Have a nice day,

**Auto generate email please do not reply.**
`

func indexURLJob(j *que.Job) error {
	var groceries models.Groceries
	if err := json.Unmarshal(j.Args, &groceries); err != nil {
		return errors.Wrap(err, "Unable to unmarshal job arguments into IndexRequest: "+string(j.Args))
	}

	groceries.Message = "Thankyou for shopping with us, Have a nice day."
	log.WithField("IndexRequest", groceries).Info(groceries.Email)

	from := mail.NewEmail("From-RunAErrand-App", "runerrandapp@test.com")
	subject := fmt.Sprintf("Your Grocery shopping order %d is dispatched", groceries.ID)
	to := mail.NewEmail("Example User", groceries.Email)

	tmpl := template.Must(template.New("template").Parse(mailtemplate))
	buf := bytes.Buffer{}
	tmpl.Execute(&buf, groceries)

	content := mail.NewContent("text/plain", buf.String())
	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		log.Errorln(err.Error())
	} else {
		log.Info(response.StatusCode)
	}

	return nil
}
