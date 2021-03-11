package kibana

import (
	"encoding/json"
)

const (
	Ascending SortOrder = iota
	Descending
)

// Enums for searchReferencesType
const (
	SearchReferencesTypeIndexPattern searchReferencesType = "index-pattern"
)

type searchReferencesType string

func (r searchReferencesType) String() string {
	return string(r)
}

type SortOrder int

type SearchSourceBuilderFactory interface {
	NewSearchSource() SearchSourceBuilder
}

type SearchClient interface {
	Create(request *CreateSearchRequest) (*Search, error)
	Update(id string, request *UpdateSearchRequest) (*Search, error)
	GetById(id string) (*Search, error)
	List() ([]*Search, error)
	Delete(id string) error
	NewSearchSource() SearchSourceBuilder
	Version() string
}

type CreateSearchRequest struct {
	Attributes *SearchAttributes   `json:"attributes"`
	References []*SearchReferences `json:"references,omitempty"`
}

type UpdateSearchRequest struct {
	Attributes *SearchAttributes   `json:"attributes"`
	References []*SearchReferences `json:"references,omitempty"`
}

type Search struct {
	Id         string              `json:"id"`
	Type       string              `json:"type"`
	Version    version             `json:"version"`
	Attributes *SearchAttributes   `json:"attributes"`
	References []*SearchReferences `json:"references,omitempty"`
}

type SearchReferences struct {
	Name string               `json:"name"`
	Type searchReferencesType `json:"type"`
	Id   string               `json:"id"`
}

type SearchAttributes struct {
	Title                 string                       `json:"title"`
	Description           string                       `json:"description"`
	Hits                  int                          `json:"hits"`
	Columns               []string                     `json:"columns"`
	Sort                  Sort                         `json:"sort"`
	Version               int                          `json:"version"`
	KibanaSavedObjectMeta *SearchKibanaSavedObjectMeta `json:"kibanaSavedObjectMeta"`
}

// Sort allows unmarshalling different json structure for the sort field. In
// newer version of Kibana this can be a nested JSON array
// (https://github.com/elastic/kibana/pull/41918/), while in the older versions
// it is a flat JSON array.
type Sort []string

// UnmarshalJSON tries to unmarshal the json data first into a nested array and if that fails it will try to unmarshal into a slice of string
func (s *Sort) UnmarshalJSON(data []byte) error {
	var nestedSlice [][]string

	err := json.Unmarshal(data, &nestedSlice)
	if err != nil {
		var slice []string
		err = json.Unmarshal(data, &slice)
		if err != nil {
			return err
		}

		*s = slice
	}

	// Successfully unmarshalled into a nested slice, which will now need to be
	// flatten into a slice of string
	if len(nestedSlice) > 0 && len(nestedSlice[0]) == 2 {
		// The sort order for the old format will be
		// taken from the first element of the new format
		sortOrder := nestedSlice[0][1]
		for _, sort := range nestedSlice {
			*s = append(*s, sort[0])
		}
		*s = append(*s, sortOrder)
	}

	return nil
}

type searchReadResult553 struct {
	Id      string            `json:"_id"`
	Type    string            `json:"_type"`
	Version version           `json:"_version"`
	Source  *SearchAttributes `json:"_source"`
}

type SearchKibanaSavedObjectMeta struct {
	SearchSourceJSON string `json:"searchSourceJSON"`
}

type SearchSource struct {
	IndexId      string          `json:"index"`
	HighlightAll bool            `json:"highlightAll"`
	Version      bool            `json:"version"`
	Query        interface{}     `json:"query,omitempty"`
	Filter       []*SearchFilter `json:"filter"`
}

type SearchQuery600 struct {
	Query    string `json:"query"`
	Language string `json:"language"`
}

type SearchQuery553 struct {
	QueryString *searchQueryString `json:"query_string"`
}

type searchQueryString struct {
	Query   string `json:"query"`
	Analyze bool   `json:"analyze_wildcard"`
}

type SearchFilter struct {
	Query  *SearchFilterQuery    `json:"query"`
	Exists *SearchFilterExists   `json:"exists"`
	Meta   *SearchFilterMetaData `json:"meta,omitempty"`
}

type SearchFilterQuery struct {
	Match map[string]*SearchFilterQueryAttributes `json:"match"`
}

type SearchFilterExists struct {
	Field string `json:"field"`
}

type SearchFilterMetaData struct {
	Index    string                       `json:"index"`
	Negate   bool                         `json:"negate"`
	Disabled bool                         `json:"disabled"`
	Alias    string                       `json:"alias"`
	Type     string                       `json:"type"`
	Key      string                       `json:"key"`
	Value    string                       `json:"value"`
	Params   *SearchFilterQueryAttributes `json:"params"`
}

type SearchFilterQueryAttributes struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}

type SearchRequestBuilder struct {
	title          string
	description    string
	displayColumns []string
	sortColumns    []string
	searchSource   *SearchSource
	references     []*SearchReferences
}

type SearchSourceBuilder interface {
	WithIndexId(indexId string) SearchSourceBuilder
	WithQuery(query string) SearchSourceBuilder
	WithFilter(filter *SearchFilter) SearchSourceBuilder
	Build() (*SearchSource, error)
}

func NewSearchRequestBuilder() *SearchRequestBuilder {
	return &SearchRequestBuilder{}
}

func (builder *SearchRequestBuilder) WithTitle(title string) *SearchRequestBuilder {
	builder.title = title
	return builder
}

func (builder *SearchRequestBuilder) WithDescription(description string) *SearchRequestBuilder {
	builder.description = description
	return builder
}

func (builder *SearchRequestBuilder) WithDisplayColumns(columns []string) *SearchRequestBuilder {
	builder.displayColumns = columns
	return builder
}

func (builder *SearchRequestBuilder) WithSortColumns(columns []string, order SortOrder) *SearchRequestBuilder {
	var sortOrder = ""
	if order == Descending {
		sortOrder = "desc"
	} else {
		sortOrder = "asc"
	}

	builder.sortColumns = append(columns, sortOrder)
	return builder
}

func (builder *SearchRequestBuilder) WithSearchSource(searchSource *SearchSource) *SearchRequestBuilder {
	builder.searchSource = searchSource
	return builder
}

func (builder *SearchRequestBuilder) WithReferences(refs []*SearchReferences) *SearchRequestBuilder {
	builder.references = refs
	return builder
}

func (builder *SearchRequestBuilder) Build() (*CreateSearchRequest, error) {
	searchSourceBytes, err := json.Marshal(builder.searchSource)
	if err != nil {
		return nil, err
	}

	request := &CreateSearchRequest{
		Attributes: &SearchAttributes{
			Title:       builder.title,
			Description: builder.description,
			Hits:        0,
			Columns:     builder.displayColumns,
			Sort:        builder.sortColumns,
			Version:     1,
			KibanaSavedObjectMeta: &SearchKibanaSavedObjectMeta{
				SearchSourceJSON: string(searchSourceBytes),
			},
		},
		References: builder.references,
	}
	return request, nil
}
