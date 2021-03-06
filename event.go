package stripe

// Event is the resource representing a Stripe event.
// For more details see https://stripe.com/docs/api#events.
type Event struct {
	ID       string     `json:"id"`
	Live     bool       `json:"livemode"`
	Created  int64      `json:"created"`
	Data     *EventData `json:"data"`
	Webhooks uint64     `json:"pending_webhooks"`
	Type     string     `json:"type"`
	Req      string     `json:"request"`
}

// EventData is the unmarshalled object as a map.
type EventData struct {
	Obj  map[string]interface{} `json:"object"`
	Prev map[string]interface{} `json:"previous_attributes"`
}

// EventListParams is the set of parameters that can be used when listing events.
// For more details see https://stripe.com/docs/api#list_events.
type EventListParams struct {
	ListParams
	Created int64
	// Type is one of the values documented at https://stripe.com/docs/api#event_types.
	Type string
}

// GetObjValue returns the value from the e.Data.Obj bag based on the keys hierarchy.
func (e *Event) GetObjValue(keys ...string) string {
	return getValue(e.Data.Obj, keys)
}

// GetPrevValue returns the value from the e.Data.Prev bag based on the keys hierarchy.
func (e *Event) GetPrevValue(keys ...string) string {
	return getValue(e.Data.Prev, keys)
}

// getValue returns the value from the m map based on the keys.
func getValue(m map[string]interface{}, keys []string) string {
	node := m[keys[0]]

	for i := 1; i < len(keys); i++ {
		node = node.(map[string]interface{})[keys[i]]
	}

	if node == nil {
		return ""
	}

	return node.(string)
}
