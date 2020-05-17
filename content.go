package main

import (
	"context"
	"log"

	"github.com/tidwall/pretty"

	icore "github.com/ipfs/interface-go-ipfs-core"
	iscn "github.com/likecoin/iscn-ipld/plugin/block"
)

func testContent(ctx context.Context, ipfs icore.CoreAPI) iscn.IscnObject {
	// --------------------------------------------------
	log.Printf("Generating content v1 block ...")
	data := map[string]interface{}{
		"type":        "article",
		"version":     1,
		"parent":      nil,
		"source":      "https://example.com/index.html",
		"edition":     "v.0.1",
		"fingerprint": "hash://sha256/9f86d081884c7d659a2feaa0",
		"title":       "Hello World!!!",
		"description": "Just to say hello to world.",
		"tags":        []string{"hello", "world", "blog"},
	}

	b1, err := iscn.Encode(iscn.CodecContent, 1, data)
	if err != nil {
		log.Panicf("Cannot create content v1 block: %s", err)
	}
	log.Printf("New content v1 block %s", b1.RawData())

	log.Printf("Generating content v2 block ...")
	data = map[string]interface{}{
		"type":        "article",
		"version":     2,
		"parent":      b1.Cid(),
		"fingerprint": "hash://sha256/9f86d081884c7d659a2feaa0",
		"title":       "Hello World!!!",
	}

	b2, err := iscn.Encode(iscn.CodecContent, 1, data)
	if err != nil {
		log.Panicf("Cannot create content v2 block: %s", err)
	}
	log.Printf("New content v2 block %s", b2.RawData())

	// --------------------------------------------------
	log.Printf("Pinning content blocks ...")

	if err := ipfs.Dag().Pinning().Add(ctx, b1); err != nil {
		log.Panicf("Cannot pin IPLD: %s", err)
	}
	if err := ipfs.Dag().Pinning().Add(ctx, b2); err != nil {
		log.Panicf("Cannot pin IPLD: %s", err)
	}

	// --------------------------------------------------
	log.Printf("Getting content blocks ...")

	ret1, err := ipfs.Dag().Get(ctx, b1.Cid())
	if err != nil {
		log.Panicf("Cannot fetch IPLD: %s", err)
		return nil
	}

	obj1, err := iscn.Decode(ret1.RawData(), b1.Cid())
	if err != nil {
		log.Panicf("Cannot decode IPLD raw data: %s", err)
	}

	ret2, err := ipfs.Dag().Get(ctx, b2.Cid())
	if err != nil {
		log.Panicf("Cannot fetch IPLD: %s", err)
		return nil
	}

	obj2, err := iscn.Decode(ret2.RawData(), b2.Cid())
	if err != nil {
		log.Panicf("Cannot decode IPLD raw data: %s", err)
	}

	// --------------------------------------------------
	// Report
	log.Println("********************************************************************************")
	log.Println("Content report")
	log.Println("********************************************************************************")

	log.Printf("Content version 1")
	c1, err := b1.Cid().StringOfBase('z')
	if err != nil {
		log.Panicf("Cannot retrieve CID from block: %s", err)
	}
	log.Printf("  CID: %s", c1)

	log.Printf("  Raw data: %s", b1.RawData())

	log.Printf("  Type: %s", obj1.GetName())
	log.Printf("  Schema version: %d", obj1.GetVersion())

	if val, err := obj1.GetString("type"); err == nil {
		log.Printf("  Content type: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj1.GetUint64("version"); err == nil {
		log.Printf("  Version: %d", val)
	} else {
		log.Panicf("%s", err)
	}

	if _, err := obj1.GetCid("parent"); err == nil {
		log.Panic("Should not have property \"parent\"")
	}

	if val, err := obj1.GetString("source"); err == nil {
		log.Printf("  Source: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj1.GetString("edition"); err == nil {
		log.Printf("  Edition: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj1.GetString("fingerprint"); err == nil {
		log.Printf("  Fingerprint: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj1.GetString("title"); err == nil {
		log.Printf("  Title: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj1.GetString("description"); err == nil {
		log.Printf("  Description: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj1.GetArray("tags"); err == nil {
		log.Printf("  Tags: %v", val)
	} else {
		log.Panicf("%s", err)
	}

	log.Printf("Content version 2")
	c2, err := b2.Cid().StringOfBase('z')
	if err != nil {
		log.Panicf("Cannot retrieve CID from block: %s", err)
	}
	log.Printf("  CID: %s", c2)

	log.Printf("  Raw data: %s", b2.RawData())

	log.Printf("  Type: %s", obj2.GetName())
	log.Printf("  Schema version: %d", obj2.GetVersion())

	if val, err := obj2.GetString("type"); err == nil {
		log.Printf("  Content type: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj2.GetUint64("version"); err == nil {
		log.Printf("  Version: %d", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj2.GetCid("parent"); err == nil {
		c, err := val.StringOfBase('z')
		if err != nil {
			log.Panicf("Cannot retrieve CID for parent block: %s", err)
		}
		log.Printf("  Parent: %s", c)
	} else {
		log.Panicf("%s", err)
	}

	if _, err := obj2.GetString("source"); err == nil {
		log.Panic("Should not have property \"source\"")
	}

	if _, err := obj2.GetString("edition"); err == nil {
		log.Panic("Should not have property \"edition\"")
	}

	if val, err := obj2.GetString("fingerprint"); err == nil {
		log.Printf("  Fingerprint: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj2.GetString("title"); err == nil {
		log.Printf("  Title: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if _, err := obj2.GetString("description"); err == nil {
		log.Panic("Should not have property \"description\"")
	}

	if _, err := obj2.GetArray("tags"); err == nil {
		log.Panic("Should not have property \"tags\"")
	}

	// --------------------------------------------------
	// JSON

	log.Println()
	log.Println("********************************************************************************")
	log.Println("Content JSON")
	log.Println("********************************************************************************")

	log.Printf("Content version 1")
	json1, err := obj1.MarshalJSON()
	if err != nil {
		log.Panicf("Cannot marshal JSON: %s", err)
	}
	log.Println(string(json1))
	log.Println(string(pretty.Pretty([]byte(json1))))

	log.Printf("Content version 2")
	json2, err := obj2.MarshalJSON()
	if err != nil {
		log.Panicf("Cannot marshal JSON: %s", err)
	}
	log.Println(string(json2))
	log.Println(string(pretty.Pretty([]byte(json2))))

	return b2
}
