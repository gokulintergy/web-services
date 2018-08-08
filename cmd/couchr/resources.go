package main

import (
	"log"
	"fmt"
	"github.com/cardiacsociety/web-services/internal/generic"
	"github.com/cardiacsociety/web-services/internal/resource"
	"time"
	"strings"
)

const resourceIdPrefix = "resource"

type ResourceDoc struct {
	Type        string    `json:"type"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Date        string    `json:"date,omitempty"`
	Category    string    `json:"category,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Keywords    []string  `json:"keywords,omitempty"`
	URL         string    `json:"url"`
}

func syncResources() {

	xi, err := generic.GetIDs(ds, "ol_resource", "where id > 20000 limit 100")
	if err != nil {
		log.Fatalln("mysql err", err)
	}

	for _, id := range xi {

		rd := &ResourceDoc{}
		fmt.Println("Syncing resource id ", id)
		r, err := resource.ByID(ds, id)
		if err != nil {
			log.Fatalln("Could not get resource id ", id, "-", err)
		}

		rd.mapResource(*r)
		id := fmt.Sprintf("%v::%v", resourceIdPrefix, r.ID)
		_, err = cb.Upsert(id, rd, 0)
		if err != nil {
			log.Println("Upsert error", err)
		}
	}
}

func (rd *ResourceDoc) mapResource(r resource.Resource) {

	// get yyyy-mm-dd part of date
	date := strings.Fields(r.PubDate.Date.String())[0]

	rd.Type = "resource"
	rd.Created = r.CreatedAt
	rd.Updated = r.UpdatedAt
	rd.Date = date
	rd.Category = r.Type
	rd.Title = r.Name
	rd.Description = r.Description
	rd.Keywords = r.Keywords
	rd.URL = r.ResourceURL
}
