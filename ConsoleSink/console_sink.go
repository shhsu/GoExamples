package console_sink

import (
	"os/exec"
	"sync/atomic"
)

const (
	// LineTypeStdout is a value for SinkData.LineType that represents a single write to Stdout
	LineTypeStdout = 0
	// LineTypeStderr is a value for SinkData.LineType that represents a single write to Stderr
	LineTypeStderr = 1
)

// SinkData captures a single write to console output
type SinkData struct {
	LineType int    `json:"type"`
	Data     []byte `json:"data"`
}

// ConsoleSink collects all console output from its execution.Runner object
type ConsoleSink struct {
	out     outputWriter
	err     errorWriter
	lines   []SinkData
	channel chan SinkData
	closed  int32
}

// NewConsoleSink create a ConsoleSink and starts collection console outputs
func NewConsoleSink() ConsoleSink {
	sink := ConsoleSink{
		lines:   make([]SinkData, 16384),
		channel: make(chan SinkData),
	}
	sink.out = outputWriter{sink: &sink}
	sink.err = errorWriter{sink: &sink}
	go sink.collect()
	return sink
}

func (sink ConsoleSink) AttachCommand(cmd *exec.Cmd) {
	cmd.Stdout = sink.out
	cmd.Stderr = sink.err
}

// GetLines gets the console outputs starting from line startIndex, inclusive
func (sink ConsoleSink) GetLines(startIndex int) (int, []SinkData) {
	// TODO: possible optimization: merge lines of the same type to reduce traffic over network?
	// Note that if we do that numLines would no longer be the size of []SinkData returned
	lines := sink.lines[startIndex:]
	numLines := len(lines)
	return numLines, lines
}

// Close will shut down the ConsoleSink. No more input will be collected
func (sink ConsoleSink) Close() error {
	atomic.StoreInt32(&sink.closed, 1)
}

// IsActive returns whether the sink is active
func (sink ConsoleSink) IsActive() bool {
	return atomic.LoadInt32(&sink.closed) == 0
}

func (sink ConsoleSink) consume(data []byte, lineType int) {
	in := SinkData{
		Data:     data,
		LineType: lineType,
	}
	sink.channel <- in
}

// collect starts up a loop for collecting all console outputs. It should be called as goroutine
func (sink ConsoleSink) collect() {
	for sink.IsActive() {
		data := <-sink.channel
		sink.lines = append(sink.lines, data)
	}
}

type outputWriter struct {
	sink *ConsoleSink
}

func (writer outputWriter) Write(p []byte) (n int, err error) {
	writer.sink.consume(p, LineTypeStdout)
	return n, nil
}

type errorWriter struct {
	sink *ConsoleSink
}

func (writer errorWriter) Write(p []byte) (n int, err error) {
	writer.sink.consume(p, LineTypeStderr)
	return n, nil
}
