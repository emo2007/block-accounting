package hal

type Link struct {
	Href   string `json:"href"`
	Title  string `json:"title,omitempty"`
	Method string `json:"method,omitempty"`
}

type Resource struct {
	Type     string          `json:"_type,omitempty"`
	Links    map[string]Link `json:"_links"`
	Embedded map[string]any  `json:"_embedded,omitempty"`
}

type Payload interface{}

type NewResourceOption func(r *Resource)

func NewResource(selfLink string, opts ...NewResourceOption) *Resource {
	r := &Resource{
		Links: map[string]Link{
			"self": {
				Href: selfLink,
			},
		},
	}

	for _, o := range opts {
		o(r)
	}

	return r
}

func WithType(t string) NewResourceOption {
	return func(r *Resource) {
		r.Type = t
	}
}

func WithSelfTitle(t string) NewResourceOption {
	return func(r *Resource) {
		if l, ok := r.Links["self"]; ok {
			l.Title = t

			r.Links["self"] = l
		}
	}
}

func (r *Resource) Embed(relation string, resource *Resource) {
	if r.Embedded == nil {
		r.Embedded = make(map[string]any, 1)
	}

	r.Embedded[relation] = resource
}

func (r *Resource) AddLink(relation string, link Link) {
	if r.Links == nil {
		r.Links = make(map[string]Link, 1)
	}

	r.Links[relation] = link
}

func (r *Resource) SetType(t string) {
	r.Type = t
}
