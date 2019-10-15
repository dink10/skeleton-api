package filter

import (
  "encoding/json"
  "fmt"
  "net/http"
  "net/http/httptest"
  "net/url"
  "strings"
  "testing"
  
  "bitbucket.org/gismart/{{Name}}/app/models"
  
  dbx "github.com/go-ozzo/ozzo-dbx"
  "github.com/google/go-cmp/cmp"
  "github.com/tidwall/pretty"
)

type testModel struct{}

func (testModel) FilteredFields() []string {
  return []string{"id", "cardinality", "name", "created_at"}
}

func (testModel) TableName() string {
  return "test"
}

func getURL(query string) (*url.URL, error) {
  return url.Parse(fmt.Sprintf(`http://google.com?%s`, pretty.Ugly([]byte(query))))
}

var (
  validOrder            = "name"
  invalidOrder          = "slug"
  validOrderDirection   = "desc"
  invalidOrderDirection = "smth"
  
  invalidFilters = []string{
    `{
			"id": false
		}`, `{
			"cardinality": []byte("1")
		}`, `{
			"name": false
		}`, `{
			"created_at": []byte("1234")
		}`,
    "someAnotherFilter",
  }
  validFilter = `{
		"name":           "test",
		"type_id":        1.0,
		"application_id": 1.0,
		"cardinality":    1.0,
		"created_at":     "123456",
		"id":             1.0
	}`
  
  validView = "references"
  
  validRange = `{
		"created_at": {
			"gt": "2006-01-02T15:04:05Z",
			"lt": "2006-01-02T15:04:05Z"
		}
	}`
  
  invalidRanges = []string{
    `{
			"created_at": {
				"gt": "2006-01-24",
				"lt": "2019-04-18"
			}
		}`,
    `{
			"created_at": {
				"gt": true,
				"lt": "1"
			}
		}`,
    `{
			"created_at": 1
		}`,
    `{
			"slug": {
				"gt": true
			}
		}`,
    `{
			"created_at": {
				"smth": "test"
			}
		}`,
    "someanotherRange",
    `{
			"created_at": {
				"gt": "2006-01-02T15:04:05Z",
				"lt": "2006-01-02T15:04:05Z"
			},
			"deleted_at": {
				"gt": "2006-01-02T15:04:05Z"
			},
			"updated_at": {
				"lt": "2006-01-02T15:04:05Z"
			}
		}`,
  }
  
  invalidCountries = []string{
    `{"include":["en", "ru"],"exclude":["en", "ru"]}`,
    "someAnotherFilter",
  }
  validCountries = []string{
    `{"include":["en", "ru"]}`,
    `{"exclude":["en", "ru"]}`,
  }
)

func TestQuery(t *testing.T) {
  t.Run("View testing...", testView(validView))
  t.Run("Valid filters testing...", testValidFilter(validFilter))
  
  for _, f := range invalidFilters {
    t.Run("Invalid filters testing...", testInvalidFilter(f))
  }
  
  t.Run("Valid order testing...", testValidOrder(validOrder, validOrderDirection))
  t.Run("Valid order with empty direction testing...", testValidOrder(validOrder, ""))
  
  t.Run("Invalid order testing...", testInvalidOrder(invalidOrder, invalidOrderDirection))
  t.Run("Invalid order direction testing...", testInvalidOrder(validOrder, invalidOrderDirection))
  t.Run("Invalid order field testing...", testInvalidOrder(invalidOrder, validOrderDirection))
  t.Run("Invalid order with empty field testing...", testInvalidOrder("", validOrderDirection))
  
  t.Run("Valid range testing...", testValidRange(validRange))
  
  for _, r := range invalidRanges {
    t.Run("Invalid range testing...", testInvalidRange(r))
  }
  
  for _, c := range validCountries {
    t.Run("Valid countries testing...", testCountries(c, validCsHandler))
  }
  
  for _, c := range invalidCountries {
    t.Run("Invalid countries testing...", testCountries(c, invalidCsHandler))
  }
}

func testValidFilter(expr string) func(t *testing.T) {
  return func(t *testing.T) {
    f, err := getFilter(expr)
    if err != nil {
      t.Errorf("Cannot parse the query: \n\t %v", err)
      return
    }
    
    var testFilter dbx.HashExp
    if err := json.Unmarshal([]byte(expr), &testFilter); err != nil {
      t.Errorf("Cannot parse the valid filter: \n\t %v", err)
      return
    }
    
    if !cmp.Equal(f, testFilter) {
      t.Errorf("Valid filter doesn't match: \n\t %s", cmp.Diff(f, testFilter))
    }
  }
}

func testInvalidFilter(expr string) func(t *testing.T) {
  return func(t *testing.T) {
    f, err := getFilter(expr)
    if err != nil {
      return
    }
    
    var testFilter dbx.HashExp
    if err := json.Unmarshal([]byte(expr), &testFilter); err != nil {
      return
    }
    
    if cmp.Equal(f, testFilter) {
      t.Errorf("Invalid filter match: \n\t %s \n\t %s \n\t %s", f, testFilter, cmp.Diff(f, testFilter))
    }
  }
}

func getFilter(filterString string) (dbx.HashExp, error) {
  u, err := getURL(fmt.Sprintf(`filter=%s`, filterString))
  if err != nil {
    return dbx.HashExp{}, err
  }
  qvs := QValues{Values: u.Query()}
  
  qp, err := qvs.ParseQuery(testModel{})
  if err != nil {
    return dbx.HashExp{}, err
  }
  
  return qp.GetFilter(), nil
}

func testView(viewValue string) func(t *testing.T) {
  return func(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(
      func(w http.ResponseWriter, vr *http.Request) {
        qViewValue := GetView(vr)
        fmt.Printf("view=%s", qViewValue)
        
        if qViewValue != viewValue {
          t.Errorf("View doesn't match: \n\t %s \n\t %s", viewValue, qViewValue)
        }
      },
    ))
    
    if _, err := http.Get(fmt.Sprintf(`%s?view=%s`, ts.URL, viewValue)); err != nil {
      t.Errorf("Test view error")
    }
    
    defer ts.Close()
    
  }
}

func testValidOrder(field, direction string) func(t *testing.T) {
  return func(t *testing.T) {
    orderString := field
    if direction != "" {
      orderString += fmt.Sprintf(`:%s`, direction)
    }
    order, err := getOrder(orderString)
    if err != nil {
      t.Errorf("Cannot parse the query: \n\t %v", err)
      return
    }
    
    o := strings.Replace(orderString, ":", " ", 1)
    
    if order != o {
      t.Errorf("Valid order doesn't match: \n\t %s \n\t %s", order, o)
    }
  }
}

func testInvalidOrder(field, direction string) func(t *testing.T) {
  return func(t *testing.T) {
    orderString := field
    if direction != "" {
      orderString += ":"
      orderString += direction
    }
    order, err := getOrder(orderString)
    if err != nil {
      return
    }
    
    o := strings.Replace(orderString, ":", " ", 1)
    
    if order == o {
      t.Errorf("Valid order doesn't match: \n\t %s \n\t %s", order, o)
    }
  }
}

func getOrder(order string) (string, error) {
  u, err := getURL(fmt.Sprintf(`order=%s`, order))
  if err != nil {
    return "", err
  }
  qvs := QValues{Values: u.Query()}
  
  qp, err := qvs.ParseQuery(testModel{})
  if err != nil {
    return "", err
  }
  
  return qp.GetOrder(), nil
}

func testValidRange(rg string) func(t *testing.T) {
  return func(t *testing.T) {
    r, err := getRange(rg)
    if err != nil {
      t.Errorf("Cannot parse the query: \n\t %v", err)
      return
    }
    
    data := make(map[string]models.DateRange)
    if err := json.Unmarshal([]byte(rg), &data); err != nil {
      t.Errorf("Cannot parse the valid range: \n\t %v", err)
      return
    }
    
    if !cmp.Equal(data, r) {
      t.Errorf("Valid date range doesn't match: \n\t %s", cmp.Diff(data, r))
    }
  }
}

func testInvalidRange(rg string) func(t *testing.T) {
  return func(t *testing.T) {
    r, err := getRange(rg)
    if err != nil {
      return
    }
    
    data := make(map[string]models.DateRange)
    if err := json.Unmarshal([]byte(rg), &data); err != nil {
      return
    }
    
    if cmp.Equal(data, r) {
      t.Errorf("Invalid date range match: \n\t %s", cmp.Diff(data, r))
    }
  }
}

func getRange(rangeString string) (map[string]models.DateRange, error) {
  u, err := getURL(fmt.Sprintf(`range=%s`, rangeString))
  if err != nil {
    return nil, err
  }
  qvs := QValues{Values: u.Query()}
  
  qp, err := qvs.ParseQuery(testModel{})
  if err != nil {
    return nil, err
  }
  
  return qp.GetRange(), nil
}

func invalidCsHandler(t *testing.T) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, vr *http.Request) {
    qCsValue, err := GetCountries(vr)
    if err == nil {
      t.Errorf("Invalid countries are valid: %+v", qCsValue)
    }
  }
}

func validCsHandler(t *testing.T) func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, vr *http.Request) {
    _, err := GetCountries(vr)
    if err != nil {
      t.Errorf("Invalid countries: %+v", err)
    }
  }
}

func testCountries(cs string, handler func(*testing.T) func(http.ResponseWriter, *http.Request)) func(t *testing.T) {
  return func(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(handler(t)))
    
    u, _ := url.Parse(ts.URL)
    q := u.Query()
    q.Set("countries", cs)
    u.RawQuery = q.Encode()
    
    if _, err := http.Get(u.String()); err != nil {
      t.Errorf("Test error: %s", err)
    }
    
    defer ts.Close()
  }
}
