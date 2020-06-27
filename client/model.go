package client

type Product struct {
	Category string  `json:"category"`
	Title    string  `json:"title"`
	Weight   float64 `json:"weight"`
	Size     Size    `json:"size"`
}

type Size struct {
	Width  float64 `json:"width"`
	Length float64 `json:"length"`
	Height float64 `json:"height"`
}
