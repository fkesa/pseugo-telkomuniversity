package compiler

import "strings"

const (
	TypeInteger = "integer"
	TypeReal    = "real"
	TypeBoolean = "boolean"
	TypeString  = "string"
	TypeChar    = "char"
)

func normalizeType(name string) string {
	switch strings.ToLower(name) {
	case "int", "integer":
		return TypeInteger
	case "float", "float32", "float64", "real":
		return TypeReal
	case "bool", "boolean":
		return TypeBoolean
	case "string":
		return TypeString
	case "char", "rune":
		return TypeChar
	default:
		return goIdent(name)
	}
}
func isPrimitiveType(name string) bool {
	switch normalizeType(name) {
	case TypeInteger, TypeReal, TypeBoolean, TypeString, TypeChar:
		return true
	default:
		return false
	}
}
func goType(name string) string {
	switch normalizeType(name) {
	case TypeInteger:
		return "int"
	case TypeReal:
		return "float64"
	case TypeBoolean:
		return "bool"
	case TypeString:
		return "string"
	case TypeChar:
		return "rune"
	default:
		return goIdent(name)
	}
}
func isNumericType(name string) bool {
	t := normalizeType(name)
	return t == TypeInteger || t == TypeReal
}
