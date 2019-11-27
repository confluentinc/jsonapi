package jsonapi

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type BadModel struct {
	ID int `jsonapi:"primary"`
}

type ModelBadTypes struct {
	ID           string     `jsonapi:"primary,badtypes"`
	StringField  string     `jsonapi:"attr,string_field"`
	FloatField   float64    `jsonapi:"attr,float_field"`
	TimeField    time.Time  `jsonapi:"attr,time_field"`
	TimePtrField *time.Time `jsonapi:"attr,time_ptr_field"`
}

type WithPointer struct {
	ID       *uint64  `jsonapi:"primary,with-pointers"`
	Name     *string  `jsonapi:"attr,name"`
	IsActive *bool    `jsonapi:"attr,is-active"`
	IntVal   *int     `jsonapi:"attr,int-val"`
	FloatVal *float32 `jsonapi:"attr,float-val"`
}

type Timestamp struct {
	ID   int        `jsonapi:"primary,timestamps"`
	Time time.Time  `jsonapi:"attr,timestamp,iso8601"`
	Next *time.Time `jsonapi:"attr,next,iso8601"`
}

type Car struct {
	ID    *string `jsonapi:"primary,cars"`
	Make  *string `jsonapi:"attr,make,omitempty"`
	Model *string `jsonapi:"attr,model,omitempty"`
	Year  *uint   `jsonapi:"attr,year,omitempty"`
}

type Post struct {
	Blog
	ID            uint64     `jsonapi:"primary,posts"             json:"id"`
	BlogID        int        `jsonapi:"attr,blog_id"              json:"-"`
	ClientID      string     `jsonapi:"client-id"                 json:"-"`
	Title         string     `jsonapi:"attr,title"                json:"title"`
	Body          string     `jsonapi:"attr,body"                 json:"body"`
	PublishedAt   *EpochTime `jsonapi:"attr,published_at,jsonify" json:"published_at"`
	Comments      []*Comment `jsonapi:"relation,comments"         json:"-"`
	LatestComment *Comment   `jsonapi:"relation,latest_comment"   json:"-"`
}

// Obviously you're unlikely to do this, but protobuf does... so test it here
// See https://github.com/golang/protobuf/blob/master/ptypes/timestamp/timestamp.pb.go
type EpochTime struct {
	Seconds int64
	Nanos   int32
}

func (p *Post) MarshalJSON() ([]byte, error) {
	publishedAt := time.Unix(p.PublishedAt.Seconds, int64(p.PublishedAt.Nanos))
	type Alias Post
	return json.Marshal(&struct {
		PublishedAt string `json:"published_at"`
		*Alias
	}{
		PublishedAt: publishedAt.Format(time.RFC3339),
		Alias:       (*Alias)(p),
	})
}

func (p *Post) UnmarshalJSON(data []byte) error {
	type Alias Post
	aux := &struct {
		PublishedAt string `json:"published_at"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, aux.PublishedAt)
	nanos := t.UnixNano()
	if err != nil {
		return err
	}
	p.PublishedAt = &EpochTime{Seconds: nanos/int64(time.Second), Nanos: int32(math.Mod(float64(nanos), float64(time.Second)))}
	return nil
}

type Comment struct {
	ID       int    `jsonapi:"primary,comments"`
	ClientID string `jsonapi:"client-id"`
	PostID   int    `jsonapi:"attr,post_id"`
	Body     string `jsonapi:"attr,body"`
}

type Book struct {
	ID          uint64  `jsonapi:"primary,books"`
	Author      string  `jsonapi:"attr,author"`
	ISBN        string  `jsonapi:"attr,isbn"`
	Title       string  `jsonapi:"attr,title,omitempty"`
	Description *string `jsonapi:"attr,description"`
	Pages       *uint   `jsonapi:"attr,pages,omitempty"`
	PublishedAt time.Time
	Tags        []string `jsonapi:"attr,tags"`
}

type Blog struct {
	ID            int       `jsonapi:"primary,blogs"`
	ClientID      string    `jsonapi:"client-id"`
	Title         string    `jsonapi:"attr,title"`
	Posts         []*Post   `jsonapi:"relation,posts"`
	CurrentPost   *Post     `jsonapi:"relation,current_post"`
	CurrentPostID int       `jsonapi:"attr,current_post_id"`
	CreatedAt     time.Time `jsonapi:"attr,created_at"`
	ViewCount     int       `jsonapi:"attr,view_count"`
}

func (b *Blog) JSONAPILinks() *Links {
	return &Links{
		"self": fmt.Sprintf("https://example.com/api/blogs/%d", b.ID),
		"comments": Link{
			Href: fmt.Sprintf("https://example.com/api/blogs/%d/comments", b.ID),
			Meta: Meta{
				"counts": map[string]uint{
					"likes":    4,
					"comments": 20,
				},
			},
		},
	}
}

func (b *Blog) JSONAPIRelationshipLinks(relation string) *Links {
	if relation == "posts" {
		return &Links{
			"related": Link{
				Href: fmt.Sprintf("https://example.com/api/blogs/%d/posts", b.ID),
				Meta: Meta{
					"count": len(b.Posts),
				},
			},
		}
	}
	if relation == "current_post" {
		return &Links{
			"self": fmt.Sprintf("https://example.com/api/posts/%s", "3"),
			"related": Link{
				Href: fmt.Sprintf("https://example.com/api/blogs/%d/current_post", b.ID),
			},
		}
	}
	return nil
}

func (b *Blog) JSONAPIMeta() *Meta {
	return &Meta{
		"detail": "extra details regarding the blog",
	}
}

func (b *Blog) JSONAPIRelationshipMeta(relation string) *Meta {
	if relation == "posts" {
		return &Meta{
			"this": map[string]interface{}{
				"can": map[string]interface{}{
					"go": []interface{}{
						"as",
						"deep",
						map[string]interface{}{
							"as": "required",
						},
					},
				},
			},
		}
	}
	if relation == "current_post" {
		return &Meta{
			"detail": "extra current_post detail",
		}
	}
	return nil
}

type BadComment struct {
	ID   uint64 `jsonapi:"primary,bad-comment"`
	Body string `jsonapi:"attr,body"`
}

func (bc *BadComment) JSONAPILinks() *Links {
	return &Links{
		"self": []string{"invalid", "should error"},
	}
}

type Company struct {
	ID        string    `jsonapi:"primary,companies"`
	Name      string    `jsonapi:"attr,name"`
	Boss      Employee  `jsonapi:"attr,boss"`
	Teams     []Team    `jsonapi:"attr,teams"`
	FoundedAt time.Time `jsonapi:"attr,founded-at,iso8601"`
}

type Team struct {
	Name    string     `jsonapi:"attr,name"`
	Leader  *Employee  `jsonapi:"attr,leader"`
	Members []Employee `jsonapi:"attr,members"`
}

type Employee struct {
	Firstname string     `jsonapi:"attr,firstname"`
	Surname   string     `jsonapi:"attr,surname"`
	Age       int        `jsonapi:"attr,age"`
	HiredAt   *time.Time `jsonapi:"attr,hired-at,iso8601"`
}

type CustomIntType int
type CustomFloatType float64
type CustomStringType string

type CustomAttributeTypes struct {
	ID string `jsonapi:"primary,customtypes"`

	Int        CustomIntType  `jsonapi:"attr,int"`
	IntPtr     *CustomIntType `jsonapi:"attr,intptr"`
	IntPtrNull *CustomIntType `jsonapi:"attr,intptrnull"`

	Float  CustomFloatType  `jsonapi:"attr,float"`
	String CustomStringType `jsonapi:"attr,string"`
}
