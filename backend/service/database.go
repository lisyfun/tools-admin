package service

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
	"tools-admin/backend/model"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	_ "github.com/lib/pq"
)

// DatabaseService 数据库服务
type DatabaseService struct {
	db            *gorm.DB
	connPool      sync.Map // map[uint]*sql.DB
	sqlSecurity   *SQLSecurityService
}

// NewDatabaseService 创建数据库服务实例
func NewDatabaseService(db *gorm.DB) *DatabaseService {
	return &DatabaseService{
		db:          db,
		sqlSecurity: NewSQLSecurityService(),
	}
}

// List 获取数据库连接列表
func (s *DatabaseService) List(req *model.DatabaseListReq) (*model.DatabaseListResp, error) {
	var total int64
	var list []model.Database

	query := s.db.Model(&model.Database{})
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	err = query.Offset(offset).Limit(req.PageSize).Find(&list).Error
	if err != nil {
		return nil, err
	}

	return &model.DatabaseListResp{
		Total: total,
		List:  list,
	}, nil
}

// Create 创建数据库连接
func (s *DatabaseService) Create(req *model.DatabaseCreateReq) error {
	db := &model.Database{
		Name:     req.Name,
		Type:     req.Type,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Database: req.Database,
	}

	return s.db.Create(db).Error
}

// TestConnection 测试数据库连接
func (s *DatabaseService) TestConnection(req *model.DatabaseTestReq) error {
	var db *sql.DB
	var err error

	switch req.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			req.Username, req.Password, req.Host, req.Port, req.Database)
		db, err = sql.Open("mysql", dsn)
	case "postgresql":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			req.Host, req.Port, req.Username, req.Password, req.Database)
		db, err = sql.Open("postgres", dsn)
	default:
		return fmt.Errorf("unsupported database type: %s", req.Type)
	}

	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	return nil
}

// getConnection 获取数据库连接
func (s *DatabaseService) getConnection(dbID uint) (*sql.DB, error) {
	// 从连接池中获取连接
	if conn, ok := s.connPool.Load(dbID); ok {
		return conn.(*sql.DB), nil
	}

	// 获取数据库配置
	var dbConfig model.Database
	err := s.db.First(&dbConfig, dbID).Error
	if err != nil {
		return nil, fmt.Errorf("database not found: %v", err)
	}

	// 创建新连接
	var db *sql.DB
	switch dbConfig.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
		db, err = sql.Open("mysql", dsn)
	case "postgresql":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Database)
		db, err = sql.Open("postgres", dsn)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbConfig.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// 存入连接池
	s.connPool.Store(dbID, db)

	return db, nil
}

// ExecuteQuery 执行SQL查询
func (s *DatabaseService) ExecuteQuery(req *model.QueryExecuteReq, c *gin.Context) (*model.QueryExecuteResp, error) {
	// 清理SQL语句
	req.SQL = s.sqlSecurity.SanitizeSQL(req.SQL)

	// SQL安全检查
	if err := s.sqlSecurity.ValidateSQL(req.SQL); err != nil {
		return nil, err
	}

	// SQL分析
	analysis := s.sqlSecurity.AnalyzeSQL(req.SQL)
	if analysis.Risk == model.RiskHigh {
		return nil, fmt.Errorf("SQL存在高风险: %s, %s", analysis.Description, analysis.Suggestion)
	}

	// 获取数据库连接
	db, err := s.getConnection(req.DatabaseID)
	if err != nil {
		return nil, err
	}

	// 记录开始时间
	startTime := time.Now()

	// 执行查询
	rows, err := db.Query(req.SQL)
	
	// 计算执行时长
	duration := time.Since(startTime).Milliseconds()

	// 创建审计日志
	audit := &model.SQLAudit{
		DatabaseID: req.DatabaseID,
		UserID:     c.GetUint("user_id"),
		Username:   c.GetString("username"),
		SQL:        req.SQL,
		Duration:   duration,
		ClientIP:   c.ClientIP(),
	}

	if err != nil {
		// 记录失败日志
		audit.Status = "failed"
		audit.Error = err.Error()
		s.db.Create(audit)
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		audit.Status = "failed"
		audit.Error = err.Error()
		s.db.Create(audit)
		return nil, fmt.Errorf("failed to get columns: %v", err)
	}

	// 准备数据容器
	resp := &model.QueryExecuteResp{
		Columns: columns,
		Rows:    make([][]interface{}, 0),
	}

	// 读取数据
	rowCount := int64(0)
	for rows.Next() {
		rowCount++
		// 创建数据容器
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// 扫描数据
		err = rows.Scan(scanArgs...)
		if err != nil {
			audit.Status = "failed"
			audit.Error = err.Error()
			s.db.Create(audit)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// 转换数据类型
		row := make([]interface{}, len(columns))
		for i, v := range values {
			if v == nil {
				row[i] = nil
				continue
			}

			switch v.(type) {
			case []byte:
				row[i] = string(v.([]byte))
			default:
				row[i] = v
			}
		}

		resp.Rows = append(resp.Rows, row)
	}

	// 记录成功日志
	audit.Status = "success"
	audit.AffectedRows = rowCount
	s.db.Create(audit)

	return resp, nil
}

// GetTables 获取表列表
func (s *DatabaseService) GetTables(req *model.TableListReq) ([]model.TableInfo, error) {
	db, err := s.getConnection(req.DatabaseID)
	if err != nil {
		return nil, err
	}

	var dbConfig model.Database
	err = s.db.First(&dbConfig, req.DatabaseID).Error
	if err != nil {
		return nil, fmt.Errorf("database not found: %v", err)
	}

	var tables []model.TableInfo
	switch dbConfig.Type {
	case "mysql":
		query := `
			SELECT 
				TABLE_NAME as name,
				TABLE_COMMENT as comment
			FROM 
				information_schema.TABLES 
			WHERE 
				TABLE_SCHEMA = ?
		`
		rows, err := db.Query(query, dbConfig.Database)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var table model.TableInfo
			err := rows.Scan(&table.Name, &table.Comment)
			if err != nil {
				return nil, err
			}
			tables = append(tables, table)
		}

	case "postgresql":
		query := `
			SELECT 
				tablename as name,
				obj_description((quote_ident(schemaname) || '.' || quote_ident(tablename))::regclass, 'pg_class') as comment
			FROM 
				pg_catalog.pg_tables
			WHERE 
				schemaname = 'public'
		`
		rows, err := db.Query(query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var table model.TableInfo
			err := rows.Scan(&table.Name, &table.Comment)
			if err != nil {
				return nil, err
			}
			tables = append(tables, table)
		}
	}

	return tables, nil
}

// GetTableSchema 获取表结构
func (s *DatabaseService) GetTableSchema(req *model.TableSchemaReq) ([]model.ColumnInfo, error) {
	db, err := s.getConnection(req.DatabaseID)
	if err != nil {
		return nil, err
	}

	var dbConfig model.Database
	err = s.db.First(&dbConfig, req.DatabaseID).Error
	if err != nil {
		return nil, fmt.Errorf("database not found: %v", err)
	}

	var columns []model.ColumnInfo
	switch dbConfig.Type {
	case "mysql":
		query := `
			SELECT 
				COLUMN_NAME as name,
				DATA_TYPE as type,
				IFNULL(CHARACTER_MAXIMUM_LENGTH, NUMERIC_PRECISION) as length,
				IS_NULLABLE = 'YES' as nullable,
				COLUMN_KEY = 'PRI' as is_primary_key,
				EXTRA = 'auto_increment' as is_auto_increment,
				COLUMN_DEFAULT as default_value,
				COLUMN_COMMENT as comment
			FROM 
				information_schema.COLUMNS 
			WHERE 
				TABLE_SCHEMA = ? AND TABLE_NAME = ?
			ORDER BY 
				ORDINAL_POSITION
		`
		rows, err := db.Query(query, dbConfig.Database, req.TableName)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var col model.ColumnInfo
			var defaultValue sql.NullString
			err := rows.Scan(
				&col.Name, &col.Type, &col.Length, &col.Nullable,
				&col.IsPrimaryKey, &col.IsAutoIncrement, &defaultValue, &col.Comment,
			)
			if err != nil {
				return nil, err
			}
			if defaultValue.Valid {
				col.DefaultValue = defaultValue.String
			}
			columns = append(columns, col)
		}

	case "postgresql":
		query := `
			SELECT 
				a.attname as name,
				format_type(a.atttypid, a.atttypmod) as type,
				CASE 
					WHEN a.atttypmod > 0 THEN a.atttypmod - 4
					ELSE a.attlen
				END as length,
				NOT a.attnotnull as nullable,
				CASE WHEN pk.contype = 'p' THEN true ELSE false END as is_primary_key,
				CASE WHEN a.attidentity != '' THEN true ELSE false END as is_auto_increment,
				pg_get_expr(ad.adbin, ad.adrelid) as default_value,
				col_description(a.attrelid, a.attnum) as comment
			FROM 
				pg_attribute a
				LEFT JOIN pg_attrdef ad ON a.attrelid = ad.adrelid AND a.attnum = ad.adnum
				LEFT JOIN pg_constraint pk ON pk.conrelid = a.attrelid AND pk.conkey[1] = a.attnum AND pk.contype = 'p'
			WHERE 
				a.attrelid = $1::regclass
				AND a.attnum > 0 
				AND NOT a.attisdropped
			ORDER BY 
				a.attnum
		`
		rows, err := db.Query(query, req.TableName)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var col model.ColumnInfo
			var defaultValue sql.NullString
			err := rows.Scan(
				&col.Name, &col.Type, &col.Length, &col.Nullable,
				&col.IsPrimaryKey, &col.IsAutoIncrement, &defaultValue, &col.Comment,
			)
			if err != nil {
				return nil, err
			}
			if defaultValue.Valid {
				col.DefaultValue = defaultValue.String
			}
			columns = append(columns, col)
		}
	}

	return columns, nil
}

// Update 更新数据库连接
func (s *DatabaseService) Update(req *model.DatabaseUpdateReq) error {
	id := req.ID
	db := model.Database{
		Name:     req.Name,
		Type:     req.Type,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Database: req.Database,
	}
	return s.db.Where("id = ?", id).Updates(&db).Error
}

// Delete 删除数据库连接
func (s *DatabaseService) Delete(id string) error {
	return s.db.Delete(&model.Database{}, id).Error
}

// TestConnectionByID 根据ID测试数据库连接
func (s *DatabaseService) TestConnectionByID(id string) error {
	var db model.Database
	if err := s.db.First(&db, id).Error; err != nil {
		return err
	}

	req := &model.DatabaseTestReq{
		Type:     db.Type,
		Host:     db.Host,
		Port:     db.Port,
		Username: db.Username,
		Password: db.Password,
		Database: db.Database,
	}
	return s.TestConnection(req)
}
