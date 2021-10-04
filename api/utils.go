package api

import (
	"encoding/xml"
	"log"
	"strings"
	"youtube-stream-notifier-linebot/controller"

	_ "github.com/lib/pq"
)

type Query struct {
	XMLName xml.Name `xml:"feed"`
	Entry   struct {
		XMLName xml.Name `xml:"entry"`
		Title   string   `xml:"title"`
		VideoId string   `xml:"id"`
	} `xml:"entry"`
}

func ParseXML(data []byte) (string, string) {
	var q Query
	xml.Unmarshal(data, &q)
	title := q.Entry.Title
	video_link := strings.Replace(q.Entry.VideoId, "yt:video:", "http://www.youtube.com/watch?v=", 1)
	return title, video_link
}

type Group struct {
	id         int
	group_id   string
	group_name string
}

func ListGroup() []Group {
	var group Group
	rows, err := controller.DB.Query(`SELECT * FROM groups;`)
	if err != nil {
		log.Fatal(err)
		return []Group{}
	}
	defer rows.Close()
	var groups []Group
	for rows.Next() {
		err := rows.Scan(&group.id, &group.group_id, &group.group_name)
		if err != nil {
			log.Fatal(err)
		}
		groups = append(groups, group)
	}
	log.Println("The number of chat groups retrieved:", len(groups))
	return groups

}

func DeleteGroup(GroupId string) {
	sqlDelete := `DELETE FROM groups WHERE group_id=$1`
	_, err := controller.DB.Exec(sqlDelete, GroupId)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Delete:", GroupId)
	}
}

func InsertGroup(GroupId string, GroupName string) {
	sqlInsert := `
	INSERT INTO groups(group_id,group_name) 
	VALUES ($1, $2)`
	_, err := controller.DB.Exec(sqlInsert, GroupId, GroupName)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Insert:", GroupId, GroupName)
	}

}
