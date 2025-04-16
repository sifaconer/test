package filters

import (
	"errors"
	"fmt"
	"strings"

	"github.com/uptrace/bun"
)

type QueryBuilder struct {
	parser *Parser
}
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		parser: NewParser(),
	}
}

func (qb *QueryBuilder) ApplyFilters(query *bun.SelectQuery, filterData interface{}, modelName string) (*bun.SelectQuery, error) {
	if filterData == nil {
		return query, nil
	}

	relations := qb.parser.ExtractRelations(filterData)
	for _, relation := range relations {
		query = query.Relation(relation)
	}

	return qb.applyFilter(query, filterData, modelName)
}

func (qb *QueryBuilder) applyFilter(query *bun.SelectQuery, filter interface{}, modelName string) (*bun.SelectQuery, error) {
	if filter == nil {
		return query, nil
	}

	switch f := filter.(type) {
	case map[string]interface{}:
		if andFilters, ok := f["AND"].([]interface{}); ok {
			return qb.applyLogicalGroup(query, andFilters, "AND", modelName)
		}
		if orFilters, ok := f["OR"].([]interface{}); ok {
			return qb.applyLogicalGroup(query, orFilters, "OR", modelName)
		}

		for field, opValue := range f {
			opMap, ok := opValue.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("formato inválido para el campo %s", field)
			}

			for op, value := range opMap {
				var err error
				query, err = qb.applyOperator(query, field, op, value)
				if err != nil {
					return nil, err
				}
			}
		}

	default:
		return nil, errors.New("formato de filtro no reconocido")
	}

	return query, nil
}

func (qb *QueryBuilder) applyLogicalGroup(query *bun.SelectQuery, filters []interface{}, logicOp string, modelName string) (*bun.SelectQuery, error) {
	if len(filters) == 0 {
		return query, nil
	}

	if logicOp == "AND" {
		var err error
		for _, subFilter := range filters {
			query, err = qb.applyFilter(query, subFilter, modelName)
			if err != nil {
				return nil, err
			}
		}
		return query, nil
	} else if logicOp == "OR" {
		
		if len(filters) == 1 {
			return qb.applyFilter(query, filters[0], modelName)
		}
		
		builder := qb
		

		query = query.WhereGroup(" OR ", func(q *bun.SelectQuery) *bun.SelectQuery {
			for _, subFilter := range filters {

				currentFilter := subFilter
				

				q = q.WhereGroup(" AND ", func(subQ *bun.SelectQuery) *bun.SelectQuery {

					resultQ, err := builder.applyFilter(subQ, currentFilter, modelName)
					if err != nil || resultQ == nil {

						return subQ
					}
					return resultQ
				})
			}
			return q
		})
		
		return query, nil
	}

	return nil, fmt.Errorf("operador lógico no reconocido: %s", logicOp)
}


func (qb *QueryBuilder) applyOperator(query *bun.SelectQuery, field, operator string, value interface{}) (*bun.SelectQuery, error) {

	isJsonField := strings.Contains(field, "->")
	

	if isJsonField {
		jsonParts := strings.Split(field, "->")
		baseField := jsonParts[0]
		jsonKey := jsonParts[1]
		
	
		jsonExpr := fmt.Sprintf("%s::jsonb->>'%s'", baseField, jsonKey)
		
		switch strings.ToLower(operator) {
		case "eq":
			return query.Where("? = ?", bun.Safe(jsonExpr), value), nil
		case "neq":
			return query.Where("? != ?", bun.Safe(jsonExpr), value), nil
		case "gt":
			return query.Where("? > ?", bun.Safe(jsonExpr), value), nil
		case "gte":
			return query.Where("? >= ?", bun.Safe(jsonExpr), value), nil
		case "lt":
			return query.Where("? < ?", bun.Safe(jsonExpr), value), nil
		case "lte":
			return query.Where("? <= ?", bun.Safe(jsonExpr), value), nil
		case "startswith":
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("el operador 'startswith' requiere un valor string para el campo %s", field)
			}
			return query.Where("? LIKE ?", bun.Safe(jsonExpr), strValue+"%"), nil
		case "istartswith":
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("el operador 'istartswith' requiere un valor string para el campo %s", field)
			}
			return query.Where("? ILIKE ?", bun.Safe(jsonExpr), strValue+"%"), nil
		case "endswith":
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("el operador 'endswith' requiere un valor string para el campo %s", field)
			}
			return query.Where("? LIKE ?", bun.Safe(jsonExpr), "%"+strValue), nil
		case "iendswith":
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("el operador 'iendswith' requiere un valor string para el campo %s", field)
			}
			return query.Where("? ILIKE ?", bun.Safe(jsonExpr), "%"+strValue), nil
		case "contains":
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("el operador 'contains' requiere un valor string para el campo %s", field)
			}
			return query.Where("? LIKE ?", bun.Safe(jsonExpr), "%"+strValue+"%"), nil
		case "icontains":
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("el operador 'icontains' requiere un valor string para el campo %s", field)
			}
			return query.Where("? ILIKE ?", bun.Safe(jsonExpr), "%"+strValue+"%"), nil
		case "in":
			return query.Where("? IN (?)", bun.Safe(jsonExpr), bun.In(value)), nil
		case "notin":
			return query.Where("? NOT IN (?)", bun.Safe(jsonExpr), bun.In(value)), nil
		case "isnull":
			return query.Where("? IS NULL", bun.Safe(jsonExpr)), nil
		case "isnotnull":
			return query.Where("? IS NOT NULL", bun.Safe(jsonExpr)), nil
		default:
			return nil, fmt.Errorf("operador no reconocido para campo JSON: %s", operator)
		}
	}
	

	isRelationField := strings.Contains(field, ".") && !strings.Contains(field, "->")
	if isRelationField {
		parts := strings.Split(field, ".")
		if len(parts) == 2 {
			relationName := parts[0]
			relationField := parts[1]
			

			query = query.Relation(relationName)
			

			switch strings.ToLower(operator) {
			case "eq":
				return query.Where("?.? = ?", bun.Ident(relationName), bun.Ident(relationField), value), nil
			case "neq":
				return query.Where("?.? != ?", bun.Ident(relationName), bun.Ident(relationField), value), nil
			case "gt":
				return query.Where("?.? > ?", bun.Ident(relationName), bun.Ident(relationField), value), nil
			case "gte":
				return query.Where("?.? >= ?", bun.Ident(relationName), bun.Ident(relationField), value), nil
			case "lt":
				return query.Where("?.? < ?", bun.Ident(relationName), bun.Ident(relationField), value), nil
			case "lte":
				return query.Where("?.? <= ?", bun.Ident(relationName), bun.Ident(relationField), value), nil
			case "startswith":
				strValue, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("el operador 'startswith' requiere un valor string para el campo %s", field)
				}
				return query.Where("?.? LIKE ?", bun.Ident(relationName), bun.Ident(relationField), strValue+"%"), nil
			case "istartswith":
				strValue, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("el operador 'istartswith' requiere un valor string para el campo %s", field)
				}
				return query.Where("?.? ILIKE ?", bun.Ident(relationName), bun.Ident(relationField), strValue+"%"), nil
			case "endswith":
				strValue, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("el operador 'endswith' requiere un valor string para el campo %s", field)
				}
				return query.Where("?.? LIKE ?", bun.Ident(relationName), bun.Ident(relationField), "%"+strValue), nil
			case "iendswith":
				strValue, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("el operador 'iendswith' requiere un valor string para el campo %s", field)
				}
				return query.Where("?.? ILIKE ?", bun.Ident(relationName), bun.Ident(relationField), "%"+strValue), nil
			case "contains":
				strValue, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("el operador 'contains' requiere un valor string para el campo %s", field)
				}
				return query.Where("?.? LIKE ?", bun.Ident(relationName), bun.Ident(relationField), "%"+strValue+"%"), nil
			case "icontains":
				strValue, ok := value.(string)
				if !ok {
					return nil, fmt.Errorf("el operador 'icontains' requiere un valor string para el campo %s", field)
				}
				return query.Where("?.? ILIKE ?", bun.Ident(relationName), bun.Ident(relationField), "%"+strValue+"%"), nil
			case "in":
				return query.Where("?.? IN (?)", bun.Ident(relationName), bun.Ident(relationField), bun.In(value)), nil
			case "notin":
				return query.Where("?.? NOT IN (?)", bun.Ident(relationName), bun.Ident(relationField), bun.In(value)), nil
			case "isnull":
				return query.Where("?.? IS NULL", bun.Ident(relationName), bun.Ident(relationField)), nil
			case "isnotnull":
				return query.Where("?.? IS NOT NULL", bun.Ident(relationName), bun.Ident(relationField)), nil
			default:
				return nil, fmt.Errorf("operador no reconocido para campo de relación: %s", operator)
			}
		}
	}
	

	processedField := qb.parser.ProcessFieldExpression(field)
	
	switch strings.ToLower(operator) {
	case "eq":
		return query.Where("? = ?", bun.Ident(processedField), value), nil
	case "neq":
		return query.Where("? != ?", bun.Ident(processedField), value), nil
	case "gt":
		return query.Where("? > ?", bun.Ident(processedField), value), nil
	case "gte":
		return query.Where("? >= ?", bun.Ident(processedField), value), nil
	case "lt":
		return query.Where("? < ?", bun.Ident(processedField), value), nil
	case "lte":
		return query.Where("? <= ?", bun.Ident(processedField), value), nil
	case "startswith":
		strValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("el operador 'startswith' requiere un valor string para el campo %s", field)
		}
		return query.Where("? LIKE ?", bun.Ident(processedField), strValue+"%"), nil
	case "istartswith":
		strValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("el operador 'istartswith' requiere un valor string para el campo %s", field)
		}
		return query.Where("? ILIKE ?", bun.Ident(processedField), strValue+"%"), nil
	case "endswith":
		strValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("el operador 'endswith' requiere un valor string para el campo %s", field)
		}
		return query.Where("? LIKE ?", bun.Ident(processedField), "%"+strValue), nil
	case "iendswith":
		strValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("el operador 'iendswith' requiere un valor string para el campo %s", field)
		}
		return query.Where("? ILIKE ?", bun.Ident(processedField), "%"+strValue), nil
	case "contains":
		strValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("el operador 'contains' requiere un valor string para el campo %s", field)
		}
		return query.Where("? LIKE ?", bun.Ident(processedField), "%"+strValue+"%"), nil
	case "icontains":
		strValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("el operador 'icontains' requiere un valor string para el campo %s", field)
		}
		return query.Where("? ILIKE ?", bun.Ident(processedField), "%"+strValue+"%"), nil
	case "in":
		return query.Where("? IN (?)", bun.Ident(processedField), bun.In(value)), nil
	case "notin":
		return query.Where("? NOT IN (?)", bun.Ident(processedField), bun.In(value)), nil
	case "between":

		values, ok := value.([]interface{})
		if !ok || len(values) != 2 {
			return nil, fmt.Errorf("el operador 'between' requiere exactamente 2 valores")
		}
		return query.Where("? BETWEEN ? AND ?", bun.Ident(processedField), values[0], values[1]), nil
	case "isnull":
		return query.Where("? IS NULL", bun.Ident(processedField)), nil
	case "isnotnull":
		return query.Where("? IS NOT NULL", bun.Ident(processedField)), nil
	default:
		return nil, fmt.Errorf("operador no reconocido: %s", operator)
	}
}


func (qb *QueryBuilder) ApplySort(query *bun.SelectQuery, sortData []map[string]map[string]string) (*bun.SelectQuery, error) {
	if sortData == nil || len(sortData) == 0 {
		return query, nil
	}

	for _, sortItem := range sortData {
		for field, dirInfo := range sortItem {
		
			isJsonField := strings.Contains(field, "->")
			

			var direction string
			if dirInfo != nil {
				if dir, ok := dirInfo["dir"]; ok {
					direction = strings.ToUpper(dir)
				}
			}
			

			if isJsonField {

				jsonParts := strings.Split(field, "->")
				baseField := jsonParts[0]
				jsonKey := jsonParts[1]
				
			
				jsonExpr := fmt.Sprintf("%s::jsonb->>'%s'", baseField, jsonKey)
				

				if direction == "DESC" {
					query = query.OrderExpr("? DESC", bun.Safe(jsonExpr))
				} else {
					query = query.OrderExpr("? ASC", bun.Safe(jsonExpr))
				}
			} else if strings.Contains(field, ".") {

				parts := strings.Split(field, ".")
				if len(parts) == 2 {
					relationName := parts[0]
					relationField := parts[1]
					
		
					query = query.Relation(relationName)
					

					if direction == "DESC" {
						query = query.OrderExpr("?.? DESC", bun.Ident(relationName), bun.Ident(relationField))
					} else {
						query = query.OrderExpr("?.? ASC", bun.Ident(relationName), bun.Ident(relationField))
					}
				}
			} else {

				processedField := qb.parser.ProcessFieldExpression(field)
				

				if direction == "DESC" {
					query = query.OrderExpr("? DESC", bun.Ident(processedField))
				} else {
					query = query.OrderExpr("? ASC", bun.Ident(processedField))
				}
			}
		}
	}

	return query, nil
}


func (qb *QueryBuilder) ApplyPagination(query *bun.SelectQuery, page, size int) *bun.SelectQuery {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	offset := (page - 1) * size
	return query.Limit(size).Offset(offset)
}


func (qb *QueryBuilder) BuildQuery(query *bun.SelectQuery, params *FilterParams, modelName string) (*bun.SelectQuery, error) {
	var err error


	var filterData interface{}
	if params.Filters != "" {
		filterData, err = qb.parser.ParseFilters(params.Filters)
		if err != nil {
			return nil, fmt.Errorf("error al parsear filtros: %w", err)
		}
	}


	if filterData != nil {
		query, err = qb.ApplyFilters(query, filterData, modelName)
		if err != nil {
			return nil, fmt.Errorf("error al aplicar filtros: %w", err)
		}
	}

	var sortData []map[string]map[string]string
	if params.Sort != "" {
		sortData, err = qb.parser.ParseSort(params.Sort)
		if err != nil {
			return nil, fmt.Errorf("error al parsear ordenamiento: %w", err)
		}
	}

	if sortData != nil {
		query, err = qb.ApplySort(query, sortData)
		if err != nil {
			return nil, fmt.Errorf("error al aplicar ordenamiento: %w", err)
		}
	}


	if params.Pagination != nil {
		query = qb.ApplyPagination(query, params.Pagination.Page, params.Pagination.Size)
	}

	return query, nil
}
