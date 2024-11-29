package eventpolling

import "github.com/L4B0MB4/EVTSRC/pkg/models"

type EventHanlder interface {
	HandleEvent(event models.Event) error
}
