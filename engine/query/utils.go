package query

import "strings"

// ============================================================================
// VALIDATION UTILITIES
// ============================================================================

// IsValidIdentifier valida se uma string é um identificador SQL válido
// Aceita apenas: letras, números e underscore
// Aceita apenas: letras, números e underscore
func IsValidIdentifier(s string) bool {
	if s == "" {
		return false
	}
	
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '_') {
			return false
		}
	}
	return true
}

// IsValidTableName valida nome de tabela
func IsValidTableName(name string) bool {
	return IsValidIdentifier(name)
}

// IsValidColumnName valida nome de coluna
// Aceita tanto "coluna" quanto "alias.coluna"
func IsValidColumnName(name string) bool {
	if name == "" {
		return false
	}

	parts := strings.Split(name, ".")
	if len(parts) > 2 {
		return false
	}

	for _, part := range parts {
		if !IsValidIdentifier(part) {
			return false
		}
	}

	return true
}

// ============================================================================
// STRING UTILITIES
// ============================================================================

// BuildPlaceholders gera string de placeholders (?, ?, ?)
func BuildPlaceholders(count int) string {
	if count <= 0 {
		return ""
	}
	return "(" + strings.Repeat("?,", count-1) + "?)"
}

// NormalizeJoinType normaliza tipo de JOIN
func NormalizeJoinType(joinType string) string {
	joinType = strings.ToUpper(strings.TrimSpace(joinType))
	switch joinType {
	case "LEFT", "RIGHT", "INNER", "OUTER", "CROSS":
		return joinType
	default:
		return "INNER"
	}
}

// NormalizeOperation normaliza operação de agregação
func NormalizeOperation(op string) string {
	return strings.ToUpper(strings.TrimSpace(op))
}
