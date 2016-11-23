package main

import (
	"fmt"
	"log"
	"time"
)

func watcher(sub Subscription) {
	log.Println("Start watcher for", sub.Name)

	var lastUpdated = time.Time{}
	var latestPublished = time.Time{}

	refresh := func(t time.Time) {
		log.Printf("Refreshing %s at %s\n", sub.Name, t)

		feed, err := fetchPttFeed(sub.FeedUrl)
		if err != nil {
			log.Println("Failed to fetch feed")
			return
		}

		feedUpdated, err := parsePttTime(feed.Updated)
		if err != nil {
			log.Println("Failed to parse feed's update time")
			return
		}

		if feedUpdated.Equal(lastUpdated) {
			// The feed XML has not changed
			return
		}

		lastUpdated = feedUpdated
		log.Printf("%s updated at %s", sub.Name, feedUpdated.Local())

		var notification NotificationMessage
		size := len(feed.EntryList)
		for i := size - 1; i >= 0; i-- {
			var entry = feed.EntryList[i]
			// Try to parse the publish time of entry
			published, err := parsePttTime(entry.Published)
			if err != nil {
				log.Fatal("Error while parsing entry's publish time")
				return
			}

			// This entry has been traversed
			if !published.After(latestPublished) {
				continue
			}

			latestPublished = published

			// Filtering
			if filteredAny(entry.Title, sub.Filters) {
				// Add this entry to notification
				notification.Body += fmt.Sprintf("><%s|%s>\n", entry.Link.Href, entry.Title)
				log.Println("New interesting entry:", entry.Title)
				continue
			}
		}

		// Send Slack notification if any interesting post was found
		if len(notification.Body) > 0 {
			banner := fmt.Sprintf("New interesting entries in *%s*\n", sub.Name)
			notification.Body = banner + notification.Body
			if contains(sub.NotifyMethods, "slack") {
				nChan <- notification
			}
		}
	}

	// Refresh when the watcher started, and then every ticks
	refresh(time.Now())
	refreshTime := time.Duration(sub.RefreshTime)
	ticker := time.NewTicker(refreshTime * time.Second)
	go func() {
		for t := range ticker.C {
			refresh(t)
		}
	}()
}
