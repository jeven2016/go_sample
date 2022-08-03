package common

type Assets struct {
	Items             *[]Item `json:"items"`
	ContinuationToken string  `json:"continuationToken"`
}

type Npm struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Asset struct {
	DownloadUrl string `json:"downloadUrl"`
	Path        string `json:"path"`
	Npm         Npm    `json:"npm"`
}

type Item struct {
	Repository string `json:"repository"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Assets     *[]Asset
}
