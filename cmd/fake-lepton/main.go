// fake-lepton - read a cptv file and send it has raw frames to thermal-recorder
//  Copyright (C) 2020, The Cacophony Project
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "log"
    "net"
    "os"
    "time"

    arg "github.com/alexflint/go-arg"
    "gopkg.in/yaml.v1"

    "github.com/TheCacophonyProject/go-cptv"
    cptvframe "github.com/TheCacophonyProject/go-cptv/cptvframe"
    lepton3 "github.com/TheCacophonyProject/lepton3"
    "github.com/TheCacophonyProject/thermal-recorder/headers"
)

const (
    SEND_SOCKET = "/var/run/lepton-frames"
    framesHz    = 9
)

type Args struct {
    CPTV string `arg:"--cptv" help:"cptv file to stream"`
}

func procArgs() Args {
    var args Args
    args.CPTV = "test.cptv"
    arg.MustParse(&args)
    return args
}

func main() {
    err := runMain()
    if err != nil {
        log.Fatal(err)
    }
}

func runMain() error {
    args := procArgs()
    if _, err := os.Stat(args.CPTV); err != nil {
        fmt.Printf("%v does not exist\n", args.CPTV)
        return err
    }

    log.Printf("dialing frame output socket %s\n", SEND_SOCKET)
    conn, err := net.DialUnix("unix", nil, &net.UnixAddr{
        Net:  "unix",
        Name: SEND_SOCKET,
    })
    if err != nil {
        fmt.Printf("error %v\n", err)
        return errors.New("error: connecting to frame output socket failed")
    }
    defer conn.Close()

    conn.SetWriteBuffer(lepton3.FrameCols * lepton3.FrameCols * 2 * 20)

    camera_specs := map[string]interface{}{
        headers.YResolution: lepton3.FrameRows,
        headers.XResolution: lepton3.FrameCols,
        headers.FrameSize:   lepton3.BytesPerFrame,
        headers.Model:       lepton3.Model,
        headers.Brand:       lepton3.Brand,
        headers.FPS:         framesHz,
    }

    cameraYAML, err := yaml.Marshal(camera_specs)
    if _, err := conn.Write(cameraYAML); err != nil {
        return err
    }

    conn.Write([]byte("\n"))

    log.Print("reading frames")
    r, err := cptv.NewFileReader(args.CPTV)
    defer r.Close()
    frame := r.Reader.EmptyFrame()
    count := 0
    // Telemetry size of 640 -64(size of telemetry words)
    var reaminingBytes [576]byte
    for {
        err := r.ReadFrame(frame)
        if err == io.EOF {
            break
        }
        buf := rawTelemetryBytes(frame.Status)
        _ = binary.Write(buf, binary.BigEndian, reaminingBytes)
        for _, row := range frame.Pix {
            for x, _ := range row {
                _ = binary.Write(buf, binary.BigEndian, row[x])
            }
        }
        count++
        if _, err := conn.Write(buf.Bytes()); err != nil {
            return err
        }
    }
    return nil
}

func rawTelemetryBytes(t cptvframe.Telemetry) *bytes.Buffer {
    var tw telemetryWords
    tw.TimeOn = ToMS(t.TimeOn)
    tw.StatusBits = ffcStateToStatus(t.FFCState)
    tw.FrameCounter = uint32(t.FrameCount)
    tw.FrameMean = t.FrameMean
    tw.FPATemp = ToK(t.TempC)
    tw.FPATempLastFFC = ToK(t.LastFFCTempC)
    tw.TimeCounterLastFFC = ToMS(t.LastFFCTime)
    buf := new(bytes.Buffer)
    binary.Write(buf, lepton3.Big16, tw)
    return buf
}

const statusFFCStateShift uint32 = 4

func ffcStateToStatus(status string) uint32 {
    var state uint32 = 3
    switch status {
    case lepton3.FFCNever:
        state = 0
    case lepton3.FFCImminent:
        state = 1
    case lepton3.FFCRunning:
        state = 2
    }
    state = state << statusFFCStateShift
    return state
}

type durationMS uint32
type centiK uint16

func ToK(c float64) centiK {
    return centiK(c*100 + 27315)
}

func ToMS(d time.Duration) durationMS {
    return durationMS(d / time.Millisecond)
}

type telemetryWords struct {
    TelemetryRevision  uint16     // 0  *
    TimeOn             durationMS // 1  *
    StatusBits         uint32     // 3  * Bit field
    Reserved5          [8]uint16  // 5  *
    SoftwareRevision   uint64     // 13 - Junk.
    Reserved17         [3]uint16  // 17 *
    FrameCounter       uint32     // 20 *
    FrameMean          uint16     // 22 * The average value from the whole frame
    FPATempCounts      uint16     // 23
    FPATemp            centiK     // 24 *
    Reserved25         [4]uint16  // 25
    FPATempLastFFC     centiK     // 29
    TimeCounterLastFFC durationMS // 30 *
}
