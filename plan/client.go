// Package plan provides the /plans APIs
package plan

import (
	"net/url"
	"strconv"

	stripe "github.com/stripe/stripe-go"
)

const (
	Day   stripe.PlanInterval = "day"
	Week  stripe.PlanInterval = "week"
	Month stripe.PlanInterval = "month"
	Year  stripe.PlanInterval = "year"
)

// Client is used to invoke /plans APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New POSTs a new plan.
// For more details see https://stripe.com/docs/api#create_plan.
func New(params *stripe.PlanParams) (*stripe.Plan, error) {
	return getC().New(params)
}

func (c Client) New(params *stripe.PlanParams) (*stripe.Plan, error) {
	body := &url.Values{
		"id":       {params.ID},
		"name":     {params.Name},
		"amount":   {strconv.FormatUint(params.Amount, 10)},
		"currency": {string(params.Currency)},
		"interval": {string(params.Interval)},
	}

	if params.IntervalCount > 0 {
		body.Add("interval_count", strconv.FormatUint(params.IntervalCount, 10))
	}

	if params.TrialPeriod > 0 {
		body.Add("trial_period_days", strconv.FormatUint(params.TrialPeriod, 10))
	}

	if len(params.Statement) > 0 {
		body.Add("statement_description", params.Statement)
	}

	params.AppendTo(body)

	plan := &stripe.Plan{}
	err := c.B.Call("POST", "/plans", c.Key, body, plan)

	return plan, err
}

// Get returns the details of a plan.
// For more details see https://stripe.com/docs/api#retrieve_plan.
func Get(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	return getC().Get(id, params)
}

func (c Client) Get(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	var body *url.Values

	if params != nil {
		body = &url.Values{}
		params.AppendTo(body)
	}

	plan := &stripe.Plan{}
	err := c.B.Call("GET", "/plans/"+id, c.Key, body, plan)

	return plan, err
}

// Update updates a plan's properties.
// For more details see https://stripe.com/docs/api#update_plan.
func Update(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	return getC().Update(id, params)
}

func (c Client) Update(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	var body *url.Values

	if params != nil {
		body = &url.Values{}

		if len(params.Name) > 0 {
			body.Add("name", params.Name)
		}

		if len(params.Statement) > 0 {
			body.Add("statement_description", params.Statement)
		}

		params.AppendTo(body)
	}

	plan := &stripe.Plan{}
	err := c.B.Call("POST", "/plans/"+id, c.Key, body, plan)

	return plan, err
}

// Del removes a plan.
// For more details see https://stripe.com/docs/api#delete_plan.
func Del(id string) error {
	return getC().Del(id)
}

func (c Client) Del(id string) error {
	return c.B.Call("DELETE", "/plans/"+id, c.Key, nil, nil)
}

// List returns a list of plans.
// For more details see https://stripe.com/docs/api#list_plans.
func List(params *stripe.PlanListParams) *Iter {
	return getC().List(params)
}

func (c Client) List(params *stripe.PlanListParams) *Iter {
	type planList struct {
		stripe.ListMeta
		Values []*stripe.Plan `json:"data"`
	}

	var body *url.Values
	var lp *stripe.ListParams

	if params != nil {
		body = &url.Values{}

		params.AppendTo(body)
		lp = &params.ListParams
	}

	return &Iter{stripe.GetIter(lp, body, func(b url.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &planList{}
		err := c.B.Call("GET", "/plans", c.Key, &b, list)

		ret := make([]interface{}, len(list.Values))
		for i, v := range list.Values {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is a iterator for list responses.
type Iter struct {
	Iter *stripe.Iter
}

// Next returns the next value in the list.
func (i *Iter) Next() (*stripe.Plan, error) {
	p, err := i.Iter.Next()
	if err != nil {
		return nil, err
	}

	return p.(*stripe.Plan), err
}

// Stop returns true if there are no more iterations to be performed.
func (i *Iter) Stop() bool {
	return i.Iter.Stop()
}

// Meta returns the list metadata.
func (i *Iter) Meta() *stripe.ListMeta {
	return i.Iter.Meta()
}
func getC() Client {
	return Client{stripe.GetBackend(), stripe.Key}
}
