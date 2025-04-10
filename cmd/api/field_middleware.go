package api

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (r *Rest) FieldMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			return err
		}

		fieldsParam := c.Query("fields")
		if fieldsParam == "" {
			return nil
		}

		parsedFields := parseFieldsParam(fieldsParam)

		var originalData interface{}
		if err := json.Unmarshal(c.Response().Body(), &originalData); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error processing response",
			})
		}

		filteredData := filterFields(originalData, parsedFields)

		filteredJSON, err := json.Marshal(filteredData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error processing response",
			})
		}

		c.Response().SetBody(filteredJSON)
		c.Response().Header.Set("Content-Type", "application/json")

		return nil
	}
}

type fieldNode struct {
	Name      string
	Children  map[string]*fieldNode
	Requested bool
}

func parseFieldsParam(fieldsParam string) *fieldNode {
	return parseQuery(fieldsParam)
}

func parseQuery(query string) *fieldNode {
	root := &fieldNode{
		Name:      "root",
		Children:  make(map[string]*fieldNode),
		Requested: false,
	}
	if query == "" {
		return root
	}

	parseQuerySegment(root, query)

	return root
}

func parseQuerySegment(parentNode *fieldNode, querySegment string) {
	var currentPos int = 0
	var fieldName string
	var nestedContent string
	var depth int = 0
	var startNestedPos int = -1

	for currentPos <= len(querySegment) {
		if currentPos == len(querySegment) || (querySegment[currentPos] == ',' && depth == 0) {
			if startNestedPos != -1 {
				fieldName = strings.TrimSpace(querySegment[0:startNestedPos])
				nestedContent = querySegment[startNestedPos+1 : currentPos-1] // Excluir los {}
			} else {
				fieldName = strings.TrimSpace(querySegment[0:currentPos])
				nestedContent = ""
			}

			if fieldName != "" {
				child, exists := parentNode.Children[fieldName]
				if !exists {
					child = &fieldNode{
						Name:      fieldName,
						Children:  make(map[string]*fieldNode),
						Requested: true,
					}
					parentNode.Children[fieldName] = child
				} else {
					child.Requested = true
				}

				if nestedContent != "" {
					parseQuerySegment(child, nestedContent)
				}
			}

			if currentPos < len(querySegment) {
				querySegment = querySegment[currentPos+1:]
				currentPos = 0
				startNestedPos = -1
				depth = 0
				continue
			} else {
				break
			}
		}

		if querySegment[currentPos] == '{' {
			if depth == 0 {
				startNestedPos = currentPos
			}
			depth++
		} else if querySegment[currentPos] == '}' {
			depth--
		}

		currentPos++
	}
}

func filterFields(data interface{}, fields *fieldNode) interface{} {
	if data == nil {
		return nil
	}

	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return data
		}

		result := make(map[string]interface{})
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key().String()
			if node, exists := fields.Children[key]; exists {
				value := iter.Value().Interface()

				if len(node.Children) > 0 {
					result[key] = filterFields(value, node)
				} else {
					result[key] = value
				}
			}
		}
		return result

	case reflect.Slice, reflect.Array:
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = filterFields(v.Index(i).Interface(), fields)
		}
		return result

	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			return filterFields(v.Elem().Interface(), fields)
		}
		return nil

	default:
		return data
	}
}
