package domain

type Redirect struct {
	Code      string
	URL       string `validate:"empty=false & format=url`
	CreatedAt int64
}
