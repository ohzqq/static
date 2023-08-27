package static

type Media struct {
	Img       string   `json:"img,omitempty"`
	Video     string   `json:"video,omitempty"`
	Caption   string   `json:"caption,omitempty"`
	Thumbnail string   `json:"thumbnail,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

func NewMedia(input string) *Media {
	var media Media

	mt := MimeType(input)
	switch {
	case mt.IsVideo():
		media.Video = input
	case mt.IsImage():
		media.Image = input
	}

	return &media
}

func (m *Media) WithThumb() *Media {
	var thumb []byte
	switch {
	case m.Video != "":
		thumb = VideoThumb(m.Video)
	case m.Img != "":
		thumb = ImageThumb(m.Img)
	}
	m.Thumbnail = ThumbToBase64(thumb)
	return m
}

func (m *Media) WithTags(tags ...string) *Media {
	m.Tags = tags
	return m
}

func (m *Media) WithCaption(caption string) *Media {
	m.Caption = caption
	return m
}
