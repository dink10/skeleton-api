package filter

import (
  "context"
  "encoding/json"
  "errors"
  "fmt"
  "net/http"
  "net/url"
  "strings"
  "time"
  
  "bitbucket.org/gismart/{{Name}}/app/models"
  
  dbx "github.com/go-ozzo/ozzo-dbx"
)

const contextQuery = models.ContextKey("query")
const defaultOrderField = "id"
const defaultPaginationLimit = -1 // A negative limit means no limit.

func GetQueryFromContext(ctx context.Context) (*models.QParams, error) {
  query, ok := ctx.Value(contextQuery).(models.QParams)
  if ok {
    return &query, nil
  }
  
  return nil, errors.New("can not cast to *models.QParams")
}

func GetDefaultQuery() *models.QParams {
  defaultQuery := &models.QParams{}
  defaultQuery.SetOrder(defaultOrderField)
  defaultQuery.SetLimit(defaultPaginationLimit)
  
  return defaultQuery
}

type QValues struct {
  url.Values
}

func (qvs QValues) ParseQuery(model models.ModelQuery) (models.QParams, error) {
  var err error
  q := models.QParams{}
  
  order, err := qvs.getOrder(model)
  if err != nil {
    return q, err
  }
  q.SetOrder(order)
  
  filter, err := qvs.getFilter(model)
  if err != nil {
    return q, err
  }
  q.SetFilter(filter)
  
  pagination, err := qvs.getPagination()
  if err != nil {
    return q, err
  }
  q.SetPagination(pagination)
  
  rangeFilter, err := qvs.getRangeFilter(model)
  if err != nil {
    return q, err
  }
  q.SetRangeFilter(rangeFilter)
  
  return q, nil
}

func (qvs QValues) getFilter(model models.ModelQuery) (dbx.HashExp, error) {
  filteredFields := model.FilteredFields()
  filter := dbx.HashExp{}
  f := make(map[string]interface{})
  
  values := qvs.Get("filter")
  if values == "" {
    return filter, nil
  }
  
  if err := json.Unmarshal([]byte(values), &f); err != nil {
    return filter, err
  }
  
  for k, v := range f {
    if !isAllowedField(k, filteredFields) {
      continue
    }
    
    if err := validateFilter(k, v); err != nil {
      return filter, err
    }
    
    filter[k] = v
  }
  
  return filter, nil
}

func validateFilter(key string, value interface{}) error {
  _, okFloat := value.(float64)
  _, okString := value.(string)
  
  switch key {
  case "id":
    if !okFloat {
      return fmt.Errorf("%s: Invalid integer value", key)
    }
  default:
    if !okString && !okFloat {
      return fmt.Errorf("%s: Invalid filter value", key)
    }
  }
  
  return nil
}

func (qvs QValues) getOrder(model models.ModelQuery) (string, error) {
  var field string
  filteredFields := model.FilteredFields()
  tableName := model.TableName()
  defaultOrder := fmt.Sprintf(`"%s"."%+v"`, tableName, defaultOrderField)
  
  value := qvs.Get("order")
  
  if value == "" {
    return defaultOrder, nil
  }
  
  valueSlice := strings.Split(value, ":")
  if len(valueSlice) > 2 {
    return field, fmt.Errorf("invalid order value")
  }
  
  if !isAllowedField(valueSlice[0], filteredFields) {
    return defaultOrder, nil
  }
  
  if len(valueSlice) == 2 {
    switch valueSlice[1] {
    case "asc", "desc":
      break
    default:
      return field, fmt.Errorf("invalid order direction")
    }
  }
  
  return strings.Join(valueSlice, " "), nil
}

func (qvs QValues) getPagination() (models.Pagination, error) {
  var err error
  var pg = models.Pagination{}
  pg.Limit = defaultPaginationLimit
  
  // When using pagination it's important to use an ORDER BY clause
  // that constrains the result rows into a unique order.
  // So don't forget to add ordering in handlers.
  
  values := qvs.Get("pagination")
  if values == "" {
    return pg, nil
  }
  
  if err := json.Unmarshal([]byte(values), &pg); err != nil {
    return pg, err
  }
  
  if pg.Limit == -1 {
    return pg, err
  }
  
  pg.Offset = pg.Page * pg.Limit
  
  return pg, err
}

func (qvs QValues) getRangeFilter(model models.ModelQuery) (map[string]models.DateRange, error) {
  filteredFields := model.FilteredFields()
  dr := make(map[string]models.DateRange)
  rangeFilter := make(map[string]models.DateRange)
  
  values := qvs.Get("range")
  if values == "" {
    return dr, nil
  }
  
  if err := json.Unmarshal([]byte(values), &rangeFilter); err != nil {
    return dr, err
  }
  
  for fName, fValue := range rangeFilter {
    if !isAllowedField(fName, filteredFields) {
      continue
    }
    
    if fValue.Lt.IsZero() {
      fValue.Lt = time.Now().UTC()
    }
    
    dr[fName] = fValue
  }
  
  return dr, nil
}

func isAllowedField(field string, filteredFields []string) bool {
  for _, allowedField := range filteredFields {
    if field == allowedField {
      return true
    }
  }
  return false
}

func GetView(r *http.Request) string {
  return r.URL.Query().Get("view")
}

func GetResults(r *http.Request) string {
  return r.URL.Query().Get("results")
}

type Exclude struct {
  Exclude []string `json:"exclude"`
}

type Include struct {
  Include []string `json:"include"`
}

type IncludeExclude struct {
  *Exclude
  *Include
}

func (cs IncludeExclude) Validate() error {
  if (cs.Include == nil) == (cs.Exclude == nil) {
    return fmt.Errorf("countries: one of 'include', 'exclude' is allowed")
  }
  
  return nil
}

func GetCountries(r *http.Request) (IncludeExclude, error) {
  countries := r.URL.Query().Get("countries")
  if countries == "" {
    return IncludeExclude{Include: &Include{make([]string, 0)}}, nil
  }
  
  var cs IncludeExclude
  
  if err := json.Unmarshal([]byte(countries), &cs); err != nil {
    return cs, fmt.Errorf("invalid countries format")
  }
  
  return cs, cs.Validate()
}
