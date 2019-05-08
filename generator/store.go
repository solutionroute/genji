package generator

const storeTmpl = `
{{ define "store" }}

{{ template "store-Struct" . }}
{{ template "store-New" . }}
{{ template "store-NewWithTx" . }}
{{ template "store-Insert" . }}
{{ template "store-Get" . }}
{{ template "store-Delete" . }}
{{ template "store-List" . }}
{{ template "store-Replace" . }}
{{ end }}
`

const storeStructTmpl = `
{{ define "store-Struct" }}
// {{$structName}}Store manages the table. It provides several typed helpers
// that simplify common operations.
type {{$structName}}Store struct {
	*genji.Store
}
{{ end }}
`

const storeNewTmpl = `
{{ define "store-New" }}
// {{.NameWithPrefix "New"}}Store creates a {{$structName}}Store.
func {{.NameWithPrefix "New"}}Store(db *genji.DB) *{{$structName}}Store {
	schema := record.Schema{
		Fields: []field.Field{
		{{- range .Fields}}
			{{- if eq .Type "string"}}
			{Name: "{{.Name}}", Type: field.String},
			{{- else if eq .Type "int64"}}
			{Name: "{{.Name}}", Type: field.Int64},
			{{- end}}
		{{- end}}
		},
	}

	return &{{$structName}}Store{Store: genji.NewStore(db, "{{$structName}}", &schema)}
}
{{ end }}
`

const storeNewWithTxTmpl = `
{{ define "store-NewWithTx" }}
// {{.NameWithPrefix "New"}}StoreWithTx creates a {{$structName}}Store valid for the lifetime of the given transaction.
func {{.NameWithPrefix "New"}}StoreWithTx(tx *genji.Tx) *{{$structName}}Store {
	schema := record.Schema{
		Fields: []field.Field{
		{{- range .Fields}}
			{{- if eq .Type "string"}}
			{Name: "{{.Name}}", Type: field.String},
			{{- else if eq .Type "int64"}}
			{Name: "{{.Name}}", Type: field.Int64},
			{{- end}}
		{{- end}}
		},
	}

	return &{{$structName}}Store{Store: genji.NewStoreWithTx(tx, "{{$structName}}", &schema)}
}
{{ end }}
`

const storeInsertTmpl = `
{{ define "store-Insert" }}
// Insert a record in the table and return the primary key.
{{- if eq .Pk.Name ""}}
func ({{$fl}} *{{$structName}}Store) Insert(record *{{$structName}}) (rowid []byte, err error) {
	return {{$fl}}.Store.Insert(record)
}
{{- else }}
func ({{$fl}} *{{$structName}}Store) Insert(record *{{$structName}}) (err error) {
	_, err = {{$fl}}.Store.Insert(record)
	return err
}
{{- end}}
{{ end }}
`

const storeGetTmpl = `
{{ define "store-Get" }}
// Get a record using its primary key.
{{- if eq .Pk.Name ""}}
func ({{$fl}} *{{$structName}}Store) Get(rowid []byte) (*{{$structName}}, error) {
{{- else}}
	{{- if eq .Pk.Type "string"}}
func ({{$fl}} *{{$structName}}Store) Get(pk string) (*{{$structName}}, error) {
	{{- else if eq .Pk.Type "int64"}}
func ({{$fl}} *{{$structName}}Store) Get(pk int64) (*{{$structName}}, error) {
	{{- end}}
{{- end}}
	var record {{$structName}}

	{{- if ne .Pk.Name ""}}
		{{- if eq .Pk.Type "string"}}
			rowid := []byte(pk)
		{{- else if eq .Pk.Type "int64"}}
			rowid := field.EncodeInt64(pk)
		{{end}}
	{{- end}}

	return &record, {{$fl}}.Store.Get(rowid, &record)
}
{{ end }}
`

const storeDeleteTmpl = `
{{ define "store-Delete" }}
{{- if ne .Pk.Name ""}}
// Delete a record using its primary key.
	{{- if eq .Pk.Type "string"}}
func ({{$fl}} *{{$structName}}Store) Delete(pk string) error {
	rowid := []byte(pk)
	{{- else if eq .Pk.Type "int64"}}
func ({{$fl}} *{{$structName}}Store) Delete(pk int64) error {
	rowid := field.EncodeInt64(pk)
	{{- end}}
	return {{$fl}}.Store.Delete(rowid)
}
{{- end}}
{{ end }}
`

const storeListTmpl = `
{{ define "store-List" }}
// List records from the specified offset. If the limit is equal to -1, it returns all records after the selected offset.
func ({{$fl}} *{{$structName}}Store) List(offset, limit int) ([]{{$structName}}, error) {
	size := limit
	if size == -1 {
		size = 0
	}
	list := make([]{{$structName}}, 0, size)
	err := {{$fl}}.Store.List(offset, limit, func(rowid []byte, r record.Record) error {
		var record {{$structName}}
		err := record.ScanRecord(r)
		if err != nil {
			return err
		}
		list = append(list, record)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}
{{ end }}
`

const storeReplaceTmpl = `
{{ define "store-Replace" }}
{{ if eq .Pk.Name ""}}
func ({{$fl}} *{{$structName}}Store) Replace(rowid []byte, record *{{$structName}}) error {
{{- else}}
	{{- if eq .Pk.Type "string"}}
func ({{$fl}} *{{$structName}}Store) Replace(pk string, record *{{$structName}}) error {
	rowid := []byte(pk)
	if record.{{ .Pk.Name }} == "" && record.{{ .Pk.Name }} != pk {
		record.{{ .Pk.Name }} = pk
	}

	{{- else if eq .Pk.Type "int64"}}
func ({{$fl}} *{{$structName}}Store) Replace(pk int64, record *{{$structName}}) error {
	rowid := field.EncodeInt64(pk)
	if record.{{ .Pk.Name }} == 0 && record.{{ .Pk.Name }} != pk {
		record.{{ .Pk.Name }} = pk
	}
	{{- end}}
{{- end}}
	return {{$fl}}.Store.Replace(rowid, record)
}
{{ end }}
`