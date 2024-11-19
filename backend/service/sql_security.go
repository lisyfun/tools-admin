package service

import (
	"fmt"
	"regexp"
	"strings"
	"tools-admin/backend/model"
)

// SQLSecurityService SQL安全服务
type SQLSecurityService struct{}

// NewSQLSecurityService 创建SQL安全服务实例
func NewSQLSecurityService() *SQLSecurityService {
	return &SQLSecurityService{}
}

// 危险SQL关键字
var dangerousKeywords = []string{
	"DELETE FROM", "DROP TABLE", "DROP DATABASE", "TRUNCATE TABLE",
	"ALTER TABLE", "UPDATE.*SET", "INSERT INTO",
}

// 危险函数
var dangerousFunctions = []string{
	"SLEEP", "BENCHMARK", "LOAD_FILE", "INTO OUTFILE",
	"INTO DUMPFILE", "INFORMATION_SCHEMA",
}

// 注入特征
var injectionPatterns = []string{
	"'.*OR.*'", "'.*AND.*'", "--.*", "/\\*.*\\*/",
	"#.*", "UNION.*SELECT", "CONCAT\\(.*\\)",
}

// AnalyzeSQL 分析SQL语句安全性
func (s *SQLSecurityService) AnalyzeSQL(sql string) *model.SQLAnalysisResult {
	sql = strings.TrimSpace(strings.ToUpper(sql))
	result := &model.SQLAnalysisResult{
		Risk:        model.RiskLow,
		Description: "SQL语句安全",
		Suggestion:  "",
	}

	// 检查危险关键字
	for _, keyword := range dangerousKeywords {
		if matched, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", keyword), sql); matched {
			result.Risk = model.RiskHigh
			result.Description = fmt.Sprintf("包含危险操作: %s", keyword)
			result.Suggestion = "请确认是否需要执行数据修改操作，建议添加 WHERE 条件限制范围"
			return result
		}
	}

	// 检查危险函数
	for _, function := range dangerousFunctions {
		if matched, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", function), sql); matched {
			result.Risk = model.RiskHigh
			result.Description = fmt.Sprintf("包含危险函数: %s", function)
			result.Suggestion = "请避免使用系统敏感函数"
			return result
		}
	}

	// 检查注入特征
	for _, pattern := range injectionPatterns {
		if matched, _ := regexp.MatchString(fmt.Sprintf("(?i)%s", pattern), sql); matched {
			result.Risk = model.RiskHigh
			result.Description = "疑似SQL注入攻击"
			result.Suggestion = "请使用参数化查询，避免拼接SQL"
			return result
		}
	}

	// 检查查询性能
	if strings.Contains(sql, "SELECT") {
		// 检查是否有 WHERE 条件
		if !strings.Contains(sql, "WHERE") {
			result.Risk = model.RiskMedium
			result.Description = "全表扫描风险"
			result.Suggestion = "建议添加 WHERE 条件限制查询范围"
			return result
		}

		// 检查是否使用 SELECT *
		if strings.Contains(sql, "SELECT *") {
			result.Risk = model.RiskMedium
			result.Description = "使用 SELECT * 可能影响性能"
			result.Suggestion = "建议明确指定需要的字段"
			return result
		}

		// 检查是否有排序
		if strings.Contains(sql, "ORDER BY") {
			result.Risk = model.RiskLow
			result.Description = "包含排序操作"
			result.Suggestion = "建议检查排序字段是否建立索引"
		}

		// 检查是否有分组
		if strings.Contains(sql, "GROUP BY") {
			result.Risk = model.RiskLow
			result.Description = "包含分组操作"
			result.Suggestion = "建议检查分组字段是否建立索引"
		}
	}

	return result
}

// ValidateSQL 验证SQL语句安全性
func (s *SQLSecurityService) ValidateSQL(sql string) error {
	result := s.AnalyzeSQL(sql)
	if result.Risk == model.RiskHigh {
		return fmt.Errorf("SQL安全风险: %s, %s", result.Description, result.Suggestion)
	}
	return nil
}

// SanitizeSQL 清理SQL语句
func (s *SQLSecurityService) SanitizeSQL(sql string) string {
	// 移除注释
	sql = regexp.MustCompile(`--.*$`).ReplaceAllString(sql, "")
	sql = regexp.MustCompile(`/\*.*?\*/`).ReplaceAllString(sql, "")
	sql = regexp.MustCompile(`#.*$`).ReplaceAllString(sql, "")

	// 移除多余的空白字符
	sql = regexp.MustCompile(`\s+`).ReplaceAllString(sql, " ")
	
	return strings.TrimSpace(sql)
}
