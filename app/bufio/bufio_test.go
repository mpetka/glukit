package bufio_test

import (
	. "github.com/alexandre-normand/glukit/app/bufio"
	"github.com/alexandre-normand/glukit/app/model"
	"log"
	"testing"
)

type statsCalibrationWriter struct {
	total      int
	batchCount int
	writeCount int
	flushCount int
}

func (w *statsCalibrationWriter) WriteCalibrationBatch(p []model.CalibrationRead) (n int, err error) {
	log.Printf("WriteCalibrationBatch with [%d] elements: %v", len(p), p)
	dayOfCalibrationReads := make([]model.DayOfCalibrationReads, 1)
	dayOfCalibrationReads[0] = model.DayOfCalibrationReads{p}

	return w.WriteCalibrationBatches(dayOfCalibrationReads)
}

func (w *statsCalibrationWriter) WriteCalibrationBatches(p []model.DayOfCalibrationReads) (n int, err error) {
	log.Printf("WriteCalibrationBatch with [%d] batches: %v", len(p), p)
	for _, dayOfData := range p {
		w.total += len(dayOfData.Reads)
	}
	log.Printf("WriteCalibrationBatch with total of %d", w.total)
	w.batchCount += len(p)
	w.writeCount++

	return len(p), nil
}

func (w *statsCalibrationWriter) Flush() error {
	w.flushCount++
	return nil
}

func TestSimpleWriteOfSingleBatch(t *testing.T) {
	statsWriter := new(statsCalibrationWriter)
	w := NewWriterSize(statsWriter, 10)
	batches := make([]model.DayOfCalibrationReads, 10)
	for i := 0; i < 10; i++ {
		calibrations := make([]model.CalibrationRead, 24)
		for j := 0; j < 24; j++ {
			calibrations[j] = model.CalibrationRead{model.Timestamp{"", 0}, 75}
		}
		batches[i] = model.DayOfCalibrationReads{calibrations}
	}
	w.WriteCalibrationBatches(batches)
	w.Flush()

	if statsWriter.total != 240 {
		t.Errorf("TestSimpleWriteOfSingleBatch failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestSimpleWriteOfSingleBatch failed: got a batchCount of %d but expected %d", statsWriter.total, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestSimpleWriteOfSingleBatch failed: got a writeCount of %d but expected %d", statsWriter.writeCount, 1)
	}

	if statsWriter.flushCount != 1 {
		t.Errorf("TestSimpleWriteOfSingleBatch failed: got a flushCount of %d but expected %d", statsWriter.flushCount, 1)
	}
}

func TestIndividualWrite(t *testing.T) {
	statsWriter := new(statsCalibrationWriter)
	w := NewWriterSize(statsWriter, 10)
	calibrations := make([]model.CalibrationRead, 24)
	for j := 0; j < 24; j++ {
		calibrations[j] = model.CalibrationRead{model.Timestamp{"", 0}, 75}
	}
	w.WriteCalibrationBatch(calibrations)
	w.Flush()

	if statsWriter.total != 24 {
		t.Errorf("TestIndividualWrite failed: got a total of %d but expected %d", statsWriter.total, 24)
	}

	if statsWriter.batchCount != 1 {
		t.Errorf("TestIndividualWrite failed: got a batchCount of %d but expected %d", statsWriter.total, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestIndividualWrite failed: got a writeCount of %d but expected %d", statsWriter.batchCount, 1)
	}

	if statsWriter.flushCount != 1 {
		t.Errorf("TestIndividualWrite failed: got a flushCount of %d but expected %d", statsWriter.flushCount, 1)
	}
}

func TestSimpleWriteLargerThanOneBatch(t *testing.T) {
	statsWriter := new(statsCalibrationWriter)
	w := NewWriterSize(statsWriter, 10)
	batches := make([]model.DayOfCalibrationReads, 11)
	for i := 0; i < 11; i++ {
		calibrations := make([]model.CalibrationRead, 24)
		for j := 0; j < 24; j++ {
			calibrations[j] = model.CalibrationRead{model.Timestamp{"", 0}, 75}
		}
		batches[i] = model.DayOfCalibrationReads{calibrations}
	}
	w.WriteCalibrationBatches(batches)

	if statsWriter.total != 240 {
		t.Errorf("TestSimpleWriteLargerThanOneBatch test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestSimpleWriteLargerThanOneBatch test: got a batchCount of %d but expected %d", statsWriter.batchCount, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestSimpleWriteLargerThanOneBatch test failed: got a writeCount of %d but expected %d", statsWriter.total, 1)
	}

	if statsWriter.flushCount != 0 {
		t.Errorf("TestSimpleWriteLargerThanOneBatch test failed: got a flushCount of %d but expected %d", statsWriter.flushCount, 0)
	}

	// Flushing should cause the extra calibration to be written
	w.Flush()

	if statsWriter.total != 264 {
		t.Errorf("TestSimpleWriteLargerThanOneBatch test failed: got a total of %d but expected %d", statsWriter.total, 264)
	}

	if statsWriter.batchCount != 11 {
		t.Errorf("TestSimpleWriteLargerThanOneBatch test: got a batchCount of %d but expected %d", statsWriter.batchCount, 11)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestSimpleWriteLargerThanOneBatch test failed: got a writeCount of %d but expected %d", statsWriter.total, 2)
	}

	if statsWriter.flushCount != 1 {
		t.Errorf("TestSimpleWriteLargerThanOneBatch test failed: got a flushCount of %d but expected %d", statsWriter.flushCount, 1)
	}
}

func TestWriteTwoFullBatches(t *testing.T) {
	statsWriter := new(statsCalibrationWriter)
	w := NewWriterSize(statsWriter, 10)
	batches := make([]model.DayOfCalibrationReads, 20)
	for i := 0; i < 20; i++ {
		calibrations := make([]model.CalibrationRead, 24)
		for j := 0; j < 24; j++ {
			calibrations[j] = model.CalibrationRead{model.Timestamp{"", 0}, 75}
		}
		batches[i] = model.DayOfCalibrationReads{calibrations}
	}
	w.WriteCalibrationBatches(batches)

	if statsWriter.total != 240 {
		t.Errorf("TestWriteTwoFullBatches test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 10 {
		t.Errorf("TestWriteTwoFullBatches test: got a batchCount of %d but expected %d", statsWriter.batchCount, 10)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteTwoFullBatches test failed: got a writeCount of %d but expected %d", statsWriter.total, 1)
	}

	if statsWriter.writeCount != 1 {
		t.Errorf("TestWriteTwoFullBatches test failed: got a writeCount of %d but expected %d", statsWriter.total, 1)
	}

	// Flushing should cause the extra batch to be written
	w.Flush()

	if statsWriter.total != 480 {
		t.Errorf("TestWriteTwoFullBatches test failed: got a total of %d but expected %d", statsWriter.total, 240)
	}

	if statsWriter.batchCount != 20 {
		t.Errorf("TestWriteTwoFullBatches test: got a batchCount of %d but expected %d", statsWriter.batchCount, 20)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteTwoFullBatches test failed: got a writeCount of %d but expected %d", statsWriter.total, 2)
	}

	if statsWriter.writeCount != 2 {
		t.Errorf("TestWriteTwoFullBatches test failed: got a writeCount of %d but expected %d", statsWriter.total, 2)
	}
}