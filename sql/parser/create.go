package parser

import (
	"github.com/asdine/genji/database"
	"github.com/asdine/genji/sql/query"
	"github.com/asdine/genji/sql/scanner"
)

// parseCreateStatement parses a create string and returns a Statement AST object.
// This function assumes the CREATE token has already been consumed.
func (p *Parser) parseCreateStatement() (query.Statement, error) {
	tok, pos, lit := p.ScanIgnoreWhitespace()
	switch tok {
	case scanner.TABLE:
		return p.parseCreateTableStatement()
	case scanner.UNIQUE:
		if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != scanner.INDEX {
			return nil, newParseError(scanner.Tokstr(tok, lit), []string{"INDEX"}, pos)
		}

		return p.parseCreateIndexStatement(true)
	case scanner.INDEX:
		return p.parseCreateIndexStatement(false)
	}

	return nil, newParseError(scanner.Tokstr(tok, lit), []string{"TABLE", "INDEX"}, pos)
}

// parseCreateTableStatement parses a create table string and returns a Statement AST object.
// This function assumes the CREATE TABLE tokens have already been consumed.
func (p *Parser) parseCreateTableStatement() (query.CreateTableStmt, error) {
	var stmt query.CreateTableStmt
	var err error

	// Parse IF NOT EXISTS
	stmt.IfNotExists, err = p.parseIfNotExists()
	if err != nil {
		return stmt, err
	}

	// Parse table name
	stmt.TableName, err = p.parseIdent()
	if err != nil {
		return stmt, err
	}

	// parse table config
	err = p.parseTableConfig(&stmt.Config)
	if err != nil {
		return stmt, err
	}

	return stmt, nil
}

func (p *Parser) parseIfNotExists() (bool, error) {
	// Parse "IF"
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok != scanner.IF {
		p.Unscan()
		return false, nil
	}

	// Parse "NOT"
	if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != scanner.NOT {
		return false, newParseError(scanner.Tokstr(tok, lit), []string{"NOT", "EXISTS"}, pos)
	}

	// Parse "EXISTS"
	if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != scanner.EXISTS {
		return false, newParseError(scanner.Tokstr(tok, lit), []string{"EXISTS"}, pos)
	}

	return true, nil
}

func (p *Parser) parseTableConfig(cfg *database.TableConfig) error {
	// Parse ( token.
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok != scanner.LPAREN {
		p.Unscan()
		return nil
	}

	var err error

	// Parse constraints.
	for {
		var fc database.FieldConstraint

		fc.Path, err = p.parseFieldRef()
		if err != nil {
			p.Unscan()
			break
		}

		fc.Type, err = p.parseType()
		if err != nil {
			return err
		}

		// Parse "PRIMARY"
		if tok, _, _ := p.ScanIgnoreWhitespace(); tok == scanner.PRIMARY {
			// Parse "KEY"
			if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != scanner.KEY {
				return newParseError(scanner.Tokstr(tok, lit), []string{"KEY"}, pos)
			}
			if len(cfg.PrimaryKey.Path) != 0 {
				return &ParseError{Message: "only one primary key is allowed"}
			}
			cfg.PrimaryKey = fc
		} else {
			p.Unscan()
			cfg.FieldConstraints = append(cfg.FieldConstraints, fc)
		}

		if tok, _, _ := p.ScanIgnoreWhitespace(); tok != scanner.COMMA {
			p.Unscan()
			break
		}
	}

	// Parse required ) token.
	if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != scanner.RPAREN {
		return newParseError(scanner.Tokstr(tok, lit), []string{")"}, pos)
	}

	return nil
}

// parseCreateIndexStatement parses a create index string and returns a Statement AST object.
// This function assumes the CREATE INDEX or CREATE UNIQUE INDEX tokens have already been consumed.
func (p *Parser) parseCreateIndexStatement(unique bool) (query.CreateIndexStmt, error) {
	var err error
	stmt := query.CreateIndexStmt{
		Unique: unique,
	}

	// Parse "IF"
	if tok, _, _ := p.ScanIgnoreWhitespace(); tok == scanner.IF {
		// Parse "NOT"
		if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != scanner.NOT {
			return stmt, newParseError(scanner.Tokstr(tok, lit), []string{"NOT", "EXISTS"}, pos)
		}

		// Parse "EXISTS"
		if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != scanner.EXISTS {
			return stmt, newParseError(scanner.Tokstr(tok, lit), []string{"EXISTS"}, pos)
		}

		stmt.IfNotExists = true
	} else {
		p.Unscan()
	}

	// Parse index name
	stmt.IndexName, err = p.parseIdent()
	if err != nil {
		return stmt, err
	}

	// Parse "ON"
	if tok, pos, lit := p.ScanIgnoreWhitespace(); tok != scanner.ON {
		return stmt, newParseError(scanner.Tokstr(tok, lit), []string{"ON"}, pos)
	}

	// Parse table name
	stmt.TableName, err = p.parseIdent()
	if err != nil {
		return stmt, err
	}

	paths, err := p.parsePathList()
	if err != nil {
		return stmt, err
	}
	if len(paths) == 0 {
		tok, pos, lit := p.ScanIgnoreWhitespace()
		return stmt, newParseError(scanner.Tokstr(tok, lit), []string{"("}, pos)
	}

	if len(paths) != 1 {
		return stmt, &ParseError{Message: "indexes on more than one field are not supported"}
	}

	stmt.Path = paths[0]

	return stmt, nil
}
