package handler

import (
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/calendar/v3"
)

type Event struct {
	ID          string         `json:"id" omitempty:"true"`
	Summary     string         `json:"summary"`
	Description string         `json:"description"`
	Conference  ConferenceData `json:"conference"`
	Date        string         `json:"date"`
	Start       DateTime       `json:"start"`
	End         DateTime       `json:"end"`
}

type ConferenceData struct {
	Label string `json:"label"`
	Uri   string `json:"uri"`
}

type Response struct {
	Events []Event
}

type DateTime struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timeZone"`
}

func (h *ClientHandler) GetEvents(c *gin.Context) {
	period := c.Query("period")

	var startTime, endTime time.Time
	switch period {
	case "day":
		dayNum, _ := strconv.Atoi(c.Query("day"))
		startTime = time.Now().AddDate(0, 0, dayNum)
		endTime = startTime.Add(24 * time.Hour)
	case "week":
		startTime = time.Now()
		endTime = startTime.AddDate(0, 0, 7)
	case "month":
		startTime = time.Now()
		endTime = startTime.AddDate(0, 1, 0)
	default:
		startTime = time.Now()
		endTime = startTime.Add(24 * time.Hour)
	}

	eType := c.Query("type")
	switch eType {
	case "conference":
	}

	startTimeStr := startTime.Format(time.RFC3339)
	endTimeStr := endTime.Format(time.RFC3339)

	eventsData, err := h.service.Events.List("primary").
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(startTimeStr).
		TimeMax(endTimeStr).
		MaxResults(10).
		OrderBy("startTime").
		Do()
	if err != nil {
		h.log.Error("failed to get events", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events"})
		return
	}

	var events []Event
	for _, item := range eventsData.Items {
		date := item.Start.DateTime
		if date == "" {
			date = item.Start.Date
		}

		events = append(events, Event{
			ID:          item.Id,
			Summary:     item.Summary,
			Description: item.Description,
			Date:        date,
			Start:       DateTime{DateTime: item.Start.DateTime},
			End:         DateTime{DateTime: item.End.DateTime},
		})
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (h *ClientHandler) CreateEvent(c *gin.Context) {
	var eventForm Event
	if err := c.ShouldBindJSON(&eventForm); err != nil {
		h.log.Error("failed to bind json", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	event := &calendar.Event{
		Summary:     eventForm.Summary,
		Description: eventForm.Description,
		ConferenceData: &calendar.ConferenceData{EntryPoints: []*calendar.EntryPoint{
			{Uri: eventForm.Conference.Uri, Label: eventForm.Conference.Label},
		}},
		Start: &calendar.EventDateTime{
			DateTime: eventForm.Start.DateTime,
			TimeZone: eventForm.Start.TimeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: eventForm.End.DateTime,
			TimeZone: eventForm.End.TimeZone,
		},
	}

	_, err := h.service.Events.Insert("primary", event).Do()
	if err != nil {
		h.log.Error("failed to create event", err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "Event created")
}

func (h *ClientHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")

	err := h.service.Events.Delete("primary", eventID).Do()
	if err != nil {
		h.log.Error("Failed to delete event", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

func (h *ClientHandler) GetEvent(c *gin.Context) {
	eventID := c.Param("id")

	event, err := h.service.Events.Get("primary", eventID).Do()
	if err != nil {
		h.log.Error("Failed to get event", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get event"})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *ClientHandler) GetFirstZoomMeeting(c *gin.Context) {
	start := time.Now().Format(time.RFC3339)
	end := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
	eventsData, err := h.service.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(start).TimeMax(end).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		h.log.Error("failed to get events", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events from Google Calendar"})
		return
	}

	for _, item := range eventsData.Items {
		if strings.Contains(item.Description, "zoom.us") {
			re := regexp.MustCompile(`https://[^"]+`)
			link := re.FindString(item.Description)

			err := exec.Command("cmd", "/c", "start", link).Run()
			if err != nil {
				h.log.Error("failed to start zoom", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events from Google Calendar"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"zoom_meeting_link": link})

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "there is no zoom meeting today"})
}
