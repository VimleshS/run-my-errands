package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"

	"github.com/Sirupsen/logrus"
	"github.com/VimleshS/run-my-errands/models"
	que "github.com/bgentry/que-go"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var qc *que.Client

//InjectQc A hookfunc to set que Client
func InjectQc(_qc *que.Client) {
	qc = _qc
}

// GroceryUploadList ...
func GroceryUploadList(w http.ResponseWriter, req *http.Request) {
	decoded := context.Get(req, "user")
	var user models.User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)
	logrus.WithField("from context", user).Info(user.Email)

	var groceries models.Groceries
	err := json.NewDecoder(req.Body).Decode(&groceries)
	if err != nil {
		logrus.WithField("NewDecoder", "req.body").Fatal(err.Error())
	}
	groceries.Email = user.Email
	groceries.Message = "Order processed and forwarded for dispatching"

	err = queueGroceries(groceries)
	if err != nil {
		json.NewEncoder(w).Encode(models.Exception{Message: err.Error()})
		logrus.WithField("EnqueueGroceryOrder", "req.body").Fatal(err.Error())
	}

	logrus.WithField("grocery", groceries).Info("queued for dispatch")
	json.NewEncoder(w).Encode(groceries)
}

// queueGroceries into the que as an encoded JSON object
func queueGroceries(groceries models.Groceries) error {
	enc, err := json.Marshal(groceries)
	if err != nil {
		return errors.Wrap(err, "Marshalling the IndexRequest")
	}

	j := que.Job{
		Type: models.IndexRequestJob,
		Args: enc,
	}

	return errors.Wrap(qc.Enqueue(&j), "Enqueueing Job")
}
