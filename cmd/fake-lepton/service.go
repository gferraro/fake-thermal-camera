package main

import (
    "errors"
    "fmt"
    "net"

    "github.com/godbus/dbus"
    "github.com/godbus/dbus/introspect"
)

const (
    dbusName = "org.cacophony.FakeLepton"
    dbusPath = "/org/cacophony/FakeLepton"
)

type service struct {
    conn *net.UnixConn
}

func startService(unixConn *net.UnixConn) error {
    conn, err := dbus.SystemBus()
    if err != nil {
        return err
    }
    reply, err := conn.RequestName(dbusName, dbus.NameFlagDoNotQueue)
    if err != nil {
        return err
    }
    if reply != dbus.RequestNameReplyPrimaryOwner {
        return errors.New("name already taken")
    }

    s := &service{conn: unixConn}
    conn.Export(s, dbusPath, dbusName)
    conn.Export(genIntrospectable(s), dbusPath, "org.freedesktop.DBus.Introspectable")
    return nil
}

func genIntrospectable(v interface{}) introspect.Introspectable {
    node := &introspect.Node{
        Interfaces: []introspect.Interface{{
            Name:    dbusName,
            Methods: introspect.Methods(v),
        }},
    }
    return introspect.NewIntrospectable(node)
}

// SendCPTV will send the raw frames of a cptv, to thermal-recorder
func (s *service) SendCPTV(filename string) *dbus.Error {
    fmt.Printf("Recieved cptv %v\n", filename)
    err := sendCPTV(s.conn, filename)
    if err != nil {
        return &dbus.Error{
            Name: dbusName + ".StayOnForError",
            Body: []interface{}{err.Error()},
        }
    }
    return nil
}
