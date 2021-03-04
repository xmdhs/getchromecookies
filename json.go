package getchromecookies

type cookie struct {
	ID     int          `json:"id"`
	Result cookieResult `json:"result"`
}

type cookieResult struct {
	Cookies []CookieResultCooky `json:"cookies"`
}

type CookieResultCooky struct {
	Domain   string  `json:"domain"`
	Expires  float64 `json:"expires"`
	HTTPOnly bool    `json:"httpOnly"`
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Priority string  `json:"priority"`
	SameSite string  `json:"sameSite"`
	Secure   bool    `json:"secure"`
	Session  bool    `json:"session"`
	Size     int     `json:"size"`
	Value    string  `json:"value"`
}
