package listRegions

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-rest")

const (
	methodGET = "GET"

	ivMethod  = "method"
	ivURI     = "uri"
	ivHeader  = "header"
	ivContent = "content"

	ovResult = "result"
	ovStatus = "status"
)

// RESTActivity is an Activity that is used to invoke a REST Operation
// inputs : {method,uri,params}
// outputs: {result}
type RESTActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new RESTActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &RESTActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *RESTActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *RESTActivity) Eval(context activity.Context) (done bool, err error) {

	method := strings.ToUpper(context.GetInput(ivMethod).(string))
	uri := context.GetInput(ivURI).(string)

	log.Debugf("REST Call: [%s] %s\n", method, uri)

	var reqBody io.Reader
	reqBody = nil

	req, err := http.NewRequest(method, uri, reqBody)

	if err != nil {
		return false, err
	}

	// Set headers
	log.Debug("Setting HTTP request headers...")
	if headers, ok := context.GetInput(ivHeader).(map[string]string); ok && len(headers) > 0 {
		for key, value := range headers {
			log.Debugf("%s: %s", key, value)
			req.Header.Set(key, value)
		}
	}

	httpTransportSettings := &http.Transport{}

	var client *http.Client

	client = &http.Client{Transport: httpTransportSettings}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return false, err
	}

	log.Debug("response Status:", resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)

	var result interface{}

	d := json.NewDecoder(bytes.NewReader(respBody))
	d.UseNumber()
	err = d.Decode(&result)

	//json.Unmarshal(respBody, &result)

	log.Debug("response Body:", result)

	context.SetOutput(ovResult, result)
	context.SetOutput(ovStatus, resp.StatusCode)

	return true, nil
}
