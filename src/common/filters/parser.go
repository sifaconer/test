package filters

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseFilters(filterParam string) (interface{}, error) {
	if filterParam == "" {
		return nil, nil
	}

	decoded, err := url.QueryUnescape(filterParam)
	if err != nil {
		decoded = filterParam
	}

	var filterData interface{}
	err = json.Unmarshal([]byte(decoded), &filterData)
	if err != nil {
		return nil, fmt.Errorf("error al parsear JSON de filtro: %w", err)
	}

	return filterData, nil
}

func (p *Parser) ParseSort(sortParam string) ([]map[string]map[string]string, error) {
	if sortParam == "" {
		return nil, nil
	}

	decoded, err := url.QueryUnescape(sortParam)
	if err != nil {
		decoded = sortParam
	}

	var sortData []map[string]map[string]string
	err = json.Unmarshal([]byte(decoded), &sortData)
	if err != nil {
		return nil, fmt.Errorf("error al parsear JSON de ordenamiento: %w", err)
	}

	return sortData, nil
}

func (p *Parser) ExtractRelations(filter interface{}) []string {
	relations := make(map[string]bool)

	var extract func(interface{})
	extract = func(f interface{}) {
		switch v := f.(type) {
		case map[string]interface{}:
			if andFilters, ok := v["AND"].([]interface{}); ok {
				for _, subFilter := range andFilters {
					extract(subFilter)
				}
			} else if orFilters, ok := v["OR"].([]interface{}); ok {
				for _, subFilter := range orFilters {
					extract(subFilter)
				}
			} else {
				for field := range v {
					if strings.Contains(field, ".") && !strings.Contains(field, "->") {
						parts := strings.Split(field, ".")
						if len(parts) > 1 {
							relations[parts[0]] = true
						}
					}
				}
			}
		}
	}

	extract(filter)

	result := make([]string, 0, len(relations))
	for relation := range relations {
		result = append(result, relation)
	}
	return result
}


func (p *Parser) ProcessFieldExpression(field string) string {
	if strings.Contains(field, "->") {
		return field
	}

	if strings.Contains(field, ".") {
		parts := strings.Split(field, ".")
		if len(parts) == 2 {
			return fmt.Sprintf("%s.%s", parts[0], parts[1])
		}
	}

	return field
}
