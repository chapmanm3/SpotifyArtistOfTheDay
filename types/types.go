package types

import "gorm.io/gorm"

type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type ExplicitContentBody struct {
	FilterEnabled bool `json:"filter_enabled"`
	FilterLocked  bool `json:"filter_locked"`
}

type ExternalUrlsBody struct {
	Spotify string `json:"spotify"`
}

type FollowersBody struct {
	Href  string `json:"href"`
	Total int    `json:"total"`
}

type ImagesBody struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type UserProfileResponse struct {
	Country         string              `json:"country"`
	DisplayName     string              `json:"display_name"`
	Email           string              `json:"email"`
	ExplicitContent ExplicitContentBody `json:"explicit_content"`
	ExternalUrls    ExternalUrlsBody    `json:"external_urls"`
	Followers       FollowersBody       `json:"followers"`
	Href            string              `json:"href"`
	Images          []ImagesBody        `json:"images"`
	Product         string              `json:"product"`
	Type            string              `json:"type"`
	Uri             string              `json:"uri"`
}

type ArtistObject struct {
	ExternalUrls ExternalUrlsBody `json:"external_urls"`
	Followers    FollowersBody    `json:"followers"`
	Genres       []string         `json:"genres"`
	Href         string           `json:"string"`
	Id           string           `json:"id"`
	Images       []ImagesBody     `json:"images"`
	Name         string           `json:"name"`
	Popularity   int              `json:"popularity"`
	Type         string           `json:"type"`
	Uri          string           `json:"uri"`
}

type UsersTopArtistsResponse struct {
	Href     string         `json:"href"`
	Limit    int            `json:"limit"`
	Next     string         `json:"next"`
	Offset   int            `json:"offset"`
	Previous int            `json:"previous"`
	Total    int            `json:"total"`
	Items    []ArtistObject `json:"items"`
}

type UserInfo struct {
	gorm.Model
	Country         string
	DisplayName     string
	Email           string
	ExplicitContent bool
	Followers       int
	ImageUrl        string
	Uri             string
	AuthInfo        AuthInfo
	Artists         []*ArtistInfo `gorm:"many2many:user_artists;"`
}

type AuthInfo struct {
	gorm.Model
	UserInfoID   int
	AccessToken  string
	TokenType    string
	Scope        string
	ExpiresIn    int
	RefreshToken string
}

type ArtistInfo struct {
	gorm.Model
	SpotifyUrl string
	SpotifyId  string
	Image      string
	Name       string
	Uri        string
}
