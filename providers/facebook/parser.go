package facebook

import (
	fb "github.com/huandu/facebook"
)

type UserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Name      string `json:"name"`
	ID        string `json:"id"`
	Locale    string `json:"locale"`
	Link      string `json:"link"`

	Picture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`

	Location struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"location"`

	Timezone float32 `json:"timezone"`
	Friends  struct {
		Summary struct {
			TotalCount int `json:"total_count"`
		} `json:"summary"`
	} `json:"friends"`

	Email string `json:"email"`
}

type FbResponse struct {
	Data    []FbPost `json:"data"`
	Paging  Paging   `json:"paging"`
	FbError fb.Error `json:"error"`
	UserProfile
}

type Paging struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

type FbPost struct {
	ID string `json:"id" facebook:"id"`

	From struct {
		ID    string `json:"id" facebook:"id"`
		Name  string `json:"name" facebook:"name"`
		Likes int    `json:"likes" facebook:"likes"`
	} `json:"from" facebook:"from"`

	To struct {
		ID   string `json:"id" facebook:"id"`
		Name string `json:"name" facebook:"name"`
	} `json:"to" facebook:"to"`

	Story       string `json:"story" facebook:"story"`
	Name        string `json:"name" facebook:"name"`
	Message     string `json:"message" facebook:"message"`
	Description string `json:"description" facebook:"description"`
	Shares      struct {
		Count int `json:"count"`
	} `json:"shares" facebook:"shares"`

	CreatedTime string `json:"created_time" facebook:"created_time"`
	UpdatedTime string `json:"updated_time" facebook:"updated_time"`

	Type       string `json:"type" facebook:"type"`
	StatusType string `json:"status_type" facebook:"status_type"`

	Picture string `json:"full_picture" facebook:"full_picture"`
	Link    string `json:"link" facebook:"link"`
	Source  string `json:"source" facebook:"source"`
	Icon    string `json:"icon" facebook:"icon"`

	Attachments struct {
		Data []struct {
			Description string `json:"description" facebook:"description"`
			Title       string `json:"title" facebook:"title"`
			Type        string `json:"type" facebook:"type"`
			URL         string `json:"url" facebook:"url"`
			Media       struct {
				Image struct {
					Width  int    `json:"width" facebook:"width"`
					Height int    `json:"height" facebook:"height"`
					Src    string `json:"src" facebook:"src"`
				} `json:"image" facebook:"image"`
			} `json:"media" facebook:"media"`
			SubAttachments struct {
				Data []struct {
					Description string `json:"description" facebook:"description"`
					Title       string `json:"title" facebook:"title"`
					Type        string `json:"type" facebook:"type"`
					URL         string `json:"url" facebook:"url"`
					Media       struct {
						Image struct {
							Width  int    `json:"width" facebook:"width"`
							Height int    `json:"height" facebook:"height"`
							Src    string `json:"src" facebook:"src"`
						} `json:"image" facebook:"image"`
					} `json:"media" facebook:"media"`
				} `json:"data" facebook:"data"`
			} `json:"subattachments" facebook:"subattachments"`
		} `json:"data" facebook:"data"`
	} `json:"attachments" facebook:"attachments"`
}

type FbResponseAccounts struct {
	Data    []FbAccount `json:"data"`
	Paging  Paging      `json:"paging"`
	FbError fb.Error    `json:"error"`
}

type FbAccount struct {
	ID          string   `json:"id" facebook:"id"`
	Name        string   `json:"name" facebook:"name"`
	AccessToken string   `json:"access_token" facebook:"access_token"`
	Category    string   `json:"category" facebook:"category"`
	Perms       []string `json:"perms" facebook:"perms"`
}
