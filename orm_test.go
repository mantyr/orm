package orm

import (
    "testing"
)

func TestEscape_string(t *testing.T) {
    text := "123'; select"

    val := Escape_string(text)
    if val != `123\'; select` {
        t.Errorf("Error Zero return in ParseDateString, date %q", val)
    }
}

func TestInsertUpdate(t *testing.T) {
    o := New()
    o.InsertUpdate("table").Set("status", "no")
    o.Up("status", "ok")

    val := o.String()
    text := "INSERT INTO `table`\r\nSET\r\n	`status` = \"no\"\r\nON DUPLICATE KEY UPDATE\r\n\t`status` = \"ok\"\r\n"
    if val != text {
        t.Errorf("Error INSERT ON DUPLICATE KEY UPDATE %q", val)
    }
}

func TestInsertUpdatePlus(t *testing.T) {
    o := New()
    o.InsertUpdate("table").Set("cnt", 1)
    o.Up("cnt = cnt + 1")

    val := o.String()
    text := "INSERT INTO `table`\r\nSET\r\n	`cnt` = \"1\"\r\nON DUPLICATE KEY UPDATE\r\n	cnt = cnt + 1\r\n"
    if val != text {
        t.Errorf("Error INSERT ON DUPLICATE KEY UPDATE %q", val)
    }
}

// It does not work - the location of the fields is mixed, it's ok
func TestInsertUpdatePlus_position(t *testing.T) {
    o := New()
    o.InsertUpdate("table").Set("param", "value").Set("cnt", 1)
    o.Up("cnt = cnt + 1").Up("param", "value2")

    val := o.String()
    text := "INSERT INTO `table`\r\nSET\r\n	`param` = \"value\"\r\n\t,`cnt` = \"1\"\r\nON DUPLICATE KEY UPDATE\r\n	cnt = cnt + 1\r\n	,`param` = \"value2\"\r\n"
    if val != text {
        t.Errorf("Error - the location of the fields is mixed %q %q", val, text)
    }
}