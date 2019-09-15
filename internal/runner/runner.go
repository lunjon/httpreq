package runner

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lunjon/httpreq/internal/rest"
)

type Runner struct {
	Spec   *Spec
	client *rest.Client
	hasRun bool
}

func NewRunner(spec *Spec, client *rest.Client) *Runner {
	return &Runner{
		Spec:   spec,
		client: client,
		hasRun: false}
}

func (runner *Runner) SetBaseURL(url string) (err error) {
	for _, req := range runner.Spec.Requests {
		err = req.SetBaseURL(url)
		if err != nil {
			return
		}
	}

	return
}

func (runner *Runner) Run(targets ...string) ([]*rest.Result, error) {
	if runner.hasRun {
		panic("a runner may only runner once")
	}

	var requests []*RequestTarget
	if len(targets) > 0 {
		for _, t := range targets {
			req, err := runner.findTarget(t)
			if err != nil {
				return nil, err
			}
			requests = append(requests, req)
		}

	} else {
		requests = runner.Spec.Requests
	}

	var results []*rest.Result
	for _, req := range requests {
		res, err := run(req, runner.client)
		if err != nil {
			return nil, err
		}

		results = append(results, res)
	}

	runner.hasRun = true

	return results, nil
}

func (runner *Runner) findTarget(id string) (*RequestTarget, error) {
	for _, req := range runner.Spec.Requests {
		if req.ID == id {
			return req, nil
		}
	}

	return nil, fmt.Errorf("unknown request ID: %s", id)
}

func run(req *RequestTarget, client *rest.Client) (res *rest.Result, err error) {
	header := http.Header{}
	for k, v := range req.Headers {
		header.Add(k, v)
	}

	var body []byte
	if req.Method == http.MethodPost {
		body, err = json.Marshal(req.Body)
		if err != nil {
			return
		}
		header.Add("Content-type", "application/json")
	}

	r, err := client.BuildRequest(req.Method, req.URL, body, header)
	if err != nil {
		return
	}

	if req.AWS != nil {
		aws := req.GetAWSSign()
		err = client.SignRequest(r, body, aws.Region, aws.Profile)
		if err != nil {
			return
		}
	}

	res = client.SendRequest(r)
	return
}
