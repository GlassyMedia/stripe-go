// Package sub provides the /subscriptions APIs
package sub

import (
	"fmt"
	"net/url"
	"strconv"

	stripe "github.com/stripe/stripe-go"
)

const (
	Trialing stripe.SubStatus = "trialing"
	Active   stripe.SubStatus = "active"
	PastDue  stripe.SubStatus = "past_due"
	Canceled stripe.SubStatus = "canceled"
	Unpaid   stripe.SubStatus = "unpaid"
)

// Client is used to invoke /subscriptions APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New POSTS a new subscription for a customer.
// For more details see https://stripe.com/docs/api#create_subscription.
func New(params *stripe.SubParams) (*stripe.Sub, error) {
	return getC().New(params)
}

func (c Client) New(params *stripe.SubParams) (*stripe.Sub, error) {
	body := &url.Values{
		"plan": {params.Plan},
	}

	if len(params.Token) > 0 {
		body.Add("card", params.Token)
	} else if params.Card != nil {
		params.Card.AppendDetails(body, true)
	}

	if len(params.Coupon) > 0 {
		body.Add("coupon", params.Coupon)
	}

	if params.TrialEnd > 0 {
		body.Add("trial_end", strconv.FormatInt(params.TrialEnd, 10))
	}

	if params.Quantity > 0 {
		body.Add("quantity", strconv.FormatUint(params.Quantity, 10))
	}

	token := c.Key
	if params.FeePercent > 0 {
		body.Add("application_fee_percent", strconv.FormatFloat(params.FeePercent, 'f', 2, 64))
	}

	params.AppendTo(body)

	sub := &stripe.Sub{}
	err := c.B.Call("POST", fmt.Sprintf("/customers/%v/subscriptions", params.Customer), token, body, sub)

	return sub, err
}

// Get returns the details of a subscription.
// For more details see https://stripe.com/docs/api#retrieve_subscription.
func Get(id string, params *stripe.SubParams) (*stripe.Sub, error) {
	return getC().Get(id, params)
}

func (c Client) Get(id string, params *stripe.SubParams) (*stripe.Sub, error) {
	body := &url.Values{}

	params.AppendTo(body)

	sub := &stripe.Sub{}
	err := c.B.Call("GET", fmt.Sprintf("/customers/%v/subscriptions/%v", params.Customer, id), c.Key, body, sub)

	return sub, err
}

// Update updates a subscription's properties.
// For more details see https://stripe.com/docs/api#update_subscription.
func Update(id string, params *stripe.SubParams) (*stripe.Sub, error) {
	return getC().Update(id, params)
}

func (c Client) Update(id string, params *stripe.SubParams) (*stripe.Sub, error) {
	body := &url.Values{}

	if len(params.Plan) > 0 {
		body.Add("plan", params.Plan)
	}

	if params.NoProrate {
		body.Add("prorate", strconv.FormatBool(false))
	}

	if len(params.Token) > 0 {
		body.Add("card", params.Token)
	} else if params.Card != nil {
		if len(params.Card.Token) > 0 {
			body.Add("card", params.Card.Token)
		} else {
			params.Card.AppendDetails(body, true)
		}
	}

	if len(params.Coupon) > 0 {
		body.Add("coupon", params.Coupon)
	}

	if params.TrialEnd > 0 {
		body.Add("trial_end", strconv.FormatInt(params.TrialEnd, 10))
	}

	if params.Quantity > 0 {
		body.Add("quantity", strconv.FormatUint(params.Quantity, 10))
	}

	token := c.Key
	if params.FeePercent > 0 {
		body.Add("application_fee_percent", strconv.FormatFloat(params.FeePercent, 'f', 2, 64))
	}

	params.AppendTo(body)

	sub := &stripe.Sub{}
	err := c.B.Call("POST", fmt.Sprintf("/customers/%v/subscriptions/%v", params.Customer, id), token, body, sub)

	return sub, err
}

// Cancel removes a subscription.
// For more details see https://stripe.com/docs/api#cancel_subscription.
func Cancel(id string, params *stripe.SubParams) error {
	return getC().Cancel(id, params)
}

func (c Client) Cancel(id string, params *stripe.SubParams) error {
	body := &url.Values{}

	if params.EndCancel {
		body.Add("at_period_end", strconv.FormatBool(true))
	}

	params.AppendTo(body)

	return c.B.Call("DELETE", fmt.Sprintf("/customers/%v/subscriptions/%v", params.Customer, id), c.Key, body, nil)
}

// List returns a list of subscriptions.
// For more details see https://stripe.com/docs/api#list_subscriptions.
func List(params *stripe.SubListParams) *Iter {
	return getC().List(params)
}

func (c Client) List(params *stripe.SubListParams) *Iter {
	body := &url.Values{}
	var lp *stripe.ListParams

	params.AppendTo(body)
	lp = &params.ListParams

	return &Iter{stripe.GetIter(lp, body, func(b url.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.SubList{}
		err := c.B.Call("GET", fmt.Sprintf("/customers/%v/subscriptions", params.Customer), c.Key, &b, list)

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
func (i *Iter) Next() (*stripe.Sub, error) {
	s, err := i.Iter.Next()
	if err != nil {
		return nil, err
	}

	return s.(*stripe.Sub), err
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
