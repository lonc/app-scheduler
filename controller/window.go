package controller

import (
	"fmt"
	"time"

	"github.com/ninjasphere/app-scheduler/model"
)

type window struct {
	from  Event
	until Event
}

func (w *window) init(m *model.Window) error {
	var err error
	if m == nil {
		return fmt.Errorf("illegal argument: m == nil")
	}
	w.from, err = newEvent(m.From, false)
	if err == nil {
		w.until, err = newEvent(m.Until, true)
	}
	return err
}

// Answer true if the window is permanently closed with respect to the specified time.
// That is: the close event is a non-recurring event has already occurred.
//
// This method is never true for recurring close events provided the open
// event has not yet occurred or itself recurring
func (w *window) isPermanentlyClosed(ref time.Time) bool {

	if w.until.hasTimestamp() && w.from.hasTimestamp() {
		openTimestamp := w.from.asTimestamp(ref)
		if w.from.hasEventOccurred(ref, ref) {
			if w.from.hasFinalEventOccurred(openTimestamp) {
				return w.until.hasEventOccurred(openTimestamp, ref)
			} else {
				return w.until.hasFinalEventOccurred(ref)
			}
		} else {
			return w.until.hasFinalEventOccurred(ref)
		}
	} else if !w.until.hasTimestamp() && !w.from.hasTimestamp() {
		return false
	} else if w.until.hasTimestamp() {
		return w.until.hasFinalEventOccurred(ref)
	} else { //if w.from.hasTimestamp() {
		return w.from.hasFinalEventOccurred(ref)
	}
}

// Answer true if the window is open with respect to the specified time.
func (w *window) isOpen(ref time.Time) bool {
	openWaitsForTime := w.from.hasTimestamp()
	closeWaitsForTime := w.until.hasTimestamp()

	if openWaitsForTime && closeWaitsForTime {

		// when both events are timestamp based, check
		// that the reference timestamp is within the boundaries
		// of those timestamp

		openTimestamp := w.from.asTimestamp(ref)
		closeTimestamp := w.until.asTimestamp(openTimestamp)

		return openTimestamp.Sub(ref) <= 0 &&
			ref.Sub(closeTimestamp) < 0 &&
			openTimestamp.Sub(closeTimestamp) < 0
	} else if !openWaitsForTime && !closeWaitsForTime {

		// when neither events are timestamp based, we have to
		// wait to wait for the open event to know we are open

		return false
	} else if closeWaitsForTime {

		// when only the close event is timestamp based we
		// the reference time is in the window, only if
		// it is less than the close event

		closeTimestamp := w.until.asTimestamp(ref)
		return ref.Sub(closeTimestamp) < 0
	} else { // if openWaitsForTime

		// when only the open event is timestamp basedd
		// we are in the window, only if reference
		// timestamp is greater than the open timestamp

		openTimestamp := w.until.asTimestamp(ref)
		return ref.Sub(openTimestamp) >= 0
	}
}

// Answer a channel that will receive an event when the next open event occurs.
func (w *window) whenOpen(ref time.Time) chan time.Time {
	return w.from.waiter(ref)
}

// Answer a channel that will receive an event when the next close event after the specified open event occurs.
func (w *window) whenClosed(opened time.Time) chan time.Time {
	return w.until.waiter(opened)
}

func (w *window) StringAt(ref time.Time) string {
	return fmt.Sprintf("window[from: %s, until:%s]", dump(w.from, ref), dump(w.until, ref))
}

func (w *window) String() string {
	return fmt.Sprintf("[%s, %s]", w.from.String(), w.until.String())
}
