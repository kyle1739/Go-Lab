package command

import (
	"context"
	"log"

	"github.com/olivere/elastic"
)

func EsInit(elasticclient *elastic.Client, ctx context.Context) {
	deleteIndex, err := elasticclient.DeleteIndex("twitter2").Do(ctx)
	if err != nil {
		println(err.Error())
		return
	}
	log.Printf("delete result %v\n", deleteIndex.Acknowledged)

	var IndexName = "twitter2"
	// Create a new index.
	// Create index
	createIndex, err := elasticclient.CreateIndex(IndexName).Do(ctx)
	if err != nil {
		log.Printf("expected CreateIndex to succeed; got: %v", err)
	}
	if createIndex == nil {
		log.Printf("expected result to be != nil; got: %v", createIndex)
	}

	log.Printf("createIndex response; got: %v", createIndex)

	mapping := ` {
					"properties" : {
						"historyUid":{
							"type":"long"
						},
						"roomName":{
							"type":"text"
						},
						"userId":{
							"type":"text"
						},
						"nickName":{
							"type":"text"
						},
						"icon":{
							"type":"text"
						},
						"stamp":{
							"type":"text"
						},
						"message":{
							"type":"text"
						}
					}
				}`

	putresp, err := elasticclient.PutMapping().Index(IndexName).BodyString(mapping).Do(context.TODO())
	if err != nil {
		log.Printf("expected put mapping to succeed; got: %v", err)
	}
	if putresp == nil {
		log.Printf("expected put mapping response; got: %v", putresp)
	}
	if !putresp.Acknowledged {
		log.Printf("expected put mapping ack; got: %v", putresp.Acknowledged)
	}

	log.Printf("putresp response; got: %v", putresp)

	getresp, err := elasticclient.GetMapping().Index(IndexName).Do(context.TODO())
	if err != nil {
		log.Printf("expected get mapping to succeed; got: %v", err)
	}
	if getresp == nil {
		log.Printf("expected get mapping response; got: %v", getresp)
	}

	log.Printf("get mapping response; got: %v", getresp)

	props, ok := getresp[IndexName]
	if !ok {
		log.Printf("expected JSON root to be of type map[string]interface{}; got: %#v", props)
	}

	log.Printf("props response; got: %v", props)

}
