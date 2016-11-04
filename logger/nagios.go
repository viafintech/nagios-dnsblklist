package logger

import (
  "fmt"
  "strconv"

  "github.com/go-playground/ansi"
  "github.com/go-playground/log"
  "github.com/go-playground/log/handlers/console"
)

const (
  space   = byte(' ')
  colon   = byte(':')
  base10  = 10
  equals  = byte('=')
  newLine = byte('\n')
  v       = "%v"
  gopath  = "GOPATH"

  ok       = "OK"
  warning  = "WARNING"
  critical = "CRITICAL"
  unknown  = "UNKNOWN"
)

/**
  The NagiosFormatFunc is the function which takes care of formating the
  log output correctly for nagios checks.

  This formatfunction is for use with the github.com/go-playground/log library.

  For outputting the 4 output types of nagios do the following:

  Define new console.log and apply format function:

  consoleLog := console.New()
  consoleLog.SetFormatFunc(logger.NagiosFormatFunc)

  From now on you can use:
  Log.Info("Everythings ok.") OK: Everythings ok.
  Log.Warn("Something is not completely correct.") => WARNING: Something is not completely correct.
  Log.Error("This check isn't working out.") => CRITICAL: This check isn't working out.
  Log.Alert("No idea what happened.") => UNKNOWN: No idea what happened.
**/
func NagiosFormatFunc(c *console.Console) console.Formatter {

  var b []byte
  var lvl string

  if c.DisplayColor() {

    var color ansi.EscSeq

    return func(e *log.Entry) []byte {
      b = b[0:0]
      color = c.GetDisplayColor(e.Level)

      b = append(b, color...)

      switch e.Level {
      case log.InfoLevel:
        lvl = ok
      case log.WarnLevel:
        lvl = warning
      case log.ErrorLevel:
        lvl = critical
      case log.AlertLevel:
        lvl = unknown
      default:
        lvl = e.Level.String()
      }

      b = append(b, lvl...)
      b = append(b, ansi.Reset...)
      b = append(b, space)
      b = append(b, e.Message...)

      for _, f := range e.Fields {
        b = append(b, space)
        b = append(b, color...)
        b = append(b, f.Key...)
        b = append(b, ansi.Reset...)
        b = append(b, equals)

        switch f.Value.(type) {
        case string:
          b = append(b, f.Value.(string)...)
        case int:
          b = strconv.AppendInt(b, int64(f.Value.(int)), base10)
        case int8:
          b = strconv.AppendInt(b, int64(f.Value.(int8)), base10)
        case int16:
          b = strconv.AppendInt(b, int64(f.Value.(int16)), base10)
        case int32:
          b = strconv.AppendInt(b, int64(f.Value.(int32)), base10)
        case int64:
          b = strconv.AppendInt(b, f.Value.(int64), base10)
        case uint:
          b = strconv.AppendUint(b, uint64(f.Value.(uint)), base10)
        case uint8:
          b = strconv.AppendUint(b, uint64(f.Value.(uint8)), base10)
        case uint16:
          b = strconv.AppendUint(b, uint64(f.Value.(uint16)), base10)
        case uint32:
          b = strconv.AppendUint(b, uint64(f.Value.(uint32)), base10)
        case uint64:
          b = strconv.AppendUint(b, f.Value.(uint64), base10)
        case bool:
          b = strconv.AppendBool(b, f.Value.(bool))
        default:
          b = append(b, fmt.Sprintf(v, f.Value)...)
        }
      }

      b = append(b, newLine)

      return b
    }
  }

  return func(e *log.Entry) []byte {
    b = b[0:0]

    switch e.Level {
    case log.InfoLevel:
      lvl = ok
    case log.WarnLevel:
      lvl = warning
    case log.ErrorLevel:
      lvl = critical
    case log.AlertLevel:
      lvl = unknown
    default:
      lvl = e.Level.String()
    }

    b = append(b, lvl...)
    b = append(b, space)
    b = append(b, e.Message...)

    for _, f := range e.Fields {
      b = append(b, space)
      b = append(b, f.Key...)
      b = append(b, equals)

      switch f.Value.(type) {
      case string:
        b = append(b, f.Value.(string)...)
      case int:
        b = strconv.AppendInt(b, int64(f.Value.(int)), base10)
      case int8:
        b = strconv.AppendInt(b, int64(f.Value.(int8)), base10)
      case int16:
        b = strconv.AppendInt(b, int64(f.Value.(int16)), base10)
      case int32:
        b = strconv.AppendInt(b, int64(f.Value.(int32)), base10)
      case int64:
        b = strconv.AppendInt(b, f.Value.(int64), base10)
      case uint:
        b = strconv.AppendUint(b, uint64(f.Value.(uint)), base10)
      case uint8:
        b = strconv.AppendUint(b, uint64(f.Value.(uint8)), base10)
      case uint16:
        b = strconv.AppendUint(b, uint64(f.Value.(uint16)), base10)
      case uint32:
        b = strconv.AppendUint(b, uint64(f.Value.(uint32)), base10)
      case uint64:
        b = strconv.AppendUint(b, f.Value.(uint64), base10)
      case bool:
        b = strconv.AppendBool(b, f.Value.(bool))
      default:
        b = append(b, fmt.Sprintf(v, f.Value)...)
      }
    }

    b = append(b, newLine)

    return b
  }
}
