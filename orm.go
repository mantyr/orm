// Very light ORM, not depend of others packages
package orm

import (
    "fmt"
    "strings"
)


type ORM struct {
    Table string
    Ignore_ bool
    Type string
    Values map[string]interface{}
    ValuesUp map[string]interface{}
    Wheres []string
}

func New(params ...string) *ORM {
    o := new(ORM)
    o.Values = make(map[string]interface{})
    o.ValuesUp = make(map[string]interface{})
    if len(params) > 0 {
        o.Table = params[0]
    }
    return o
}
func (o *ORM) SetType(Type string) *ORM {
    if Type == "select" || Type == "insert" || Type == "update" || Type == "insert_update" || Type == "delete" {
        o.Type = Type
    }
    return o
}
func (o *ORM) Select(table_name string) *ORM {
    if table_name != "" {
        o.Table = table_name
    }
    o.Type = "select"
    return o
}
func (o *ORM) Insert(table_name string) *ORM {
    if table_name != "" {
        o.Table = table_name
    }
    o.Type = "insert"
    return o
}
func (o *ORM) Update(table_name string) *ORM {
    if table_name != "" {
        o.Table = table_name
    }
    o.Type = "update"
    return o
}
func (o *ORM) InsertUpdate(table_name string) *ORM {
    if table_name != "" {
        o.Table = table_name
    }
    o.Type = "insert_update"
    return o
}
func (o *ORM) Delete(table_name string) *ORM {
    if table_name != "" {
        o.Table = table_name
    }
    o.Type = "delete"
    return o
}
func (o *ORM) SetArr(values map[string]interface{}) *ORM {
    o.Values = values
    return o
}
func (o *ORM) Set(name string, params ...interface{}) *ORM {
    if len(params) > 0 {
        o.Values[name] = params[0]
    } else {
        o.Values[name] = nil
    }
    return o
}

func (o *ORM) UpArr(values map[string]interface{}) *ORM {
    o.ValuesUp = values
    return o
}
func (o *ORM) Up(name string, params ...interface{}) *ORM {
    if len(params) > 0 {
        o.ValuesUp[name] = params[0]
    } else {
        o.ValuesUp[name] = nil
    }
    return o
}

func (o *ORM) SetUp(name string, params ...interface{}) *ORM {
    o.Set(name, params...)
    o.Up(name, params...)
    return o
}

func (o *ORM) Where(st string, params ...interface{}) *ORM {
    for key, param := range params {
        params[key] = Escape_string(Convert_string(param))
    }
    st = strings.Replace(st, "?", "%v", len(params)) // заменить на регулярку что бы \\? заменять на просто знак ?
    st = fmt.Sprintf(st, params...)
    o.Wheres = append(o.Wheres, st)
    return o
}

func (o *ORM) Ignore(val bool) *ORM {
    o.Ignore_ = val
    return o
}

func (o *ORM) String() string {
    if o.Type == "insert_update" && len(o.ValuesUp) == 0 {
        o.Type = "insert"
    }

    var query string
    if o.Type == "insert" {
        if (o.Ignore_) {
            query = "INSERT IGNORE INTO `"+o.Table+"`\r\nSET\r\n"
        } else {
            query = "INSERT INTO `"+o.Table+"`\r\nSET\r\n"
        }
    } else if o.Type == "update" {
        query = "UPDATE `"+o.Table+"`\r\nSET\r\n"
    } else if o.Type == "delete" {
        query = "DELETE FROM `"+o.Table+"`\r\n"
    } else if o.Type == "insert_update" {
        query = "INSERT INTO `"+o.Table+"`\r\nSET\r\n"
    }

    first := true

    // SET section
    if o.Type == "insert" || o.Type == "update" || o.Type == "insert_update"{
        for key, value := range o.Values {
            if first {
                query += "	"
            } else {
                query += "	,"
            }
            if value == nil {
                query += Escape_string(key)+"\r\n"
            } else {
                query += "`"+Escape_string(key)+"` = \""+Escape_string(Convert_string(value))+"\"\r\n"
            }
            first = false
        }
    }
    if o.Type == "insert_update" {
        query += "ON DUPLICATE KEY UPDATE\r\n"
        first = true
        for key, value := range o.ValuesUp {
            if first {
                query += "	"
            } else {
                query += "	,"
            }
            if value == nil {
                query += Escape_string(key)+"\r\n"
            } else {
                query += "`"+Escape_string(key)+"` = \""+Escape_string(Convert_string(value))+"\"\r\n"
            }
            first = false
        }
    }
    // WHERE section
    if o.Type == "select" || o.Type == "update" || o.Type == "delete" {
        if len(o.Wheres) > 0 {
            query += "WHERE\r\n"
        }
        first = true
        for _, value := range o.Wheres {
            if first {
                query += "	"
            } else {
                query += "	AND	"
            }
            query += value
            first = false
        }
    }
    return query
}

func Convert_string(i interface{}) string {
    return fmt.Sprintf("%v", i)
}

func Escape_string(text string) string {
    text = strings.Replace(text, "'", "\\'", -1)
    text = strings.Replace(text, "\"", "\\\"", -1)
    return text
}