package main

import (
	"context"
	"log"

	"github.com/tidwall/pretty"

	icore "github.com/ipfs/interface-go-ipfs-core"
	iscn "github.com/likecoin/iscn-ipld/plugin/block"
)

func testEntity(ctx context.Context, ipfs icore.CoreAPI) []iscn.IscnObject {
	// --------------------------------------------------
	log.Printf("Generating entity block 1 ...")
	data := map[string]interface{}{
		"id":          "lcc://id/comsos1xxxxxxxxxxxxxxxxxxxxxx",
		"name":        "Alice",
		"description": "I am the Alice.",
	}

	b1, err := iscn.Encode(iscn.CodecEntity, 1, data)
	if err != nil {
		log.Panicf("Cannot create entity block 1: %s", err)
	}
	log.Printf("New entity block 1 %s", b1.RawData())

	log.Printf("Generating entity block 2 ...")
	data = map[string]interface{}{
		"id": "lcc://id/comsos1yyyyyyyyyyyyyyyyyyyyyy",
	}

	b2, err := iscn.Encode(iscn.CodecEntity, 1, data)
	if err != nil {
		log.Panicf("Cannot create entity block 2: %s", err)
	}
	log.Printf("New entity block 2 %s", b2.RawData())

	log.Printf("Generating entity block 3 ...")
	data = map[string]interface{}{
		"id":   "lcc://id/comsos1zzzzzzzzzzzzzzzzzzzzzz",
		"name": "Calos",
	}

	b3, err := iscn.Encode(iscn.CodecEntity, 1, data)
	if err != nil {
		log.Panicf("Cannot create entity block 3: %s", err)
	}
	log.Printf("New entity block 3 %s", b3.RawData())

	// --------------------------------------------------
	log.Printf("Pinning entity blocks ...")

	if err := ipfs.Dag().Pinning().Add(ctx, b1); err != nil {
		log.Panicf("Cannot pin IPLD: %s", err)
	}
	if err := ipfs.Dag().Pinning().Add(ctx, b2); err != nil {
		log.Panicf("Cannot pin IPLD: %s", err)
	}
	if err := ipfs.Dag().Pinning().Add(ctx, b3); err != nil {
		log.Panicf("Cannot pin IPLD: %s", err)
	}

	// --------------------------------------------------
	log.Printf("Getting entity blocks ...")

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

	ret3, err := ipfs.Dag().Get(ctx, b3.Cid())
	if err != nil {
		log.Panicf("Cannot fetch IPLD: %s", err)
		return nil
	}

	obj3, err := iscn.Decode(ret3.RawData(), b3.Cid())
	if err != nil {
		log.Panicf("Cannot decode IPLD raw data: %s", err)
	}

	// --------------------------------------------------
	// Report
	log.Println("********************************************************************************")
	log.Println("Entity report")
	log.Println("********************************************************************************")

	log.Printf("Entity 1")
	c1, err := b1.Cid().StringOfBase('z')
	if err != nil {
		log.Panicf("Cannot retrieve CID from block: %s", err)
	}
	log.Printf("  CID: %s", c1)

	log.Printf("  Raw data: %s", b1.RawData())

	log.Printf("  Type: %s", obj1.GetName())
	log.Printf("  Schema version: %d", obj1.GetVersion())

	if val, err := obj1.GetString("id"); err == nil {
		log.Printf("  ID: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj1.GetString("name"); err == nil {
		log.Printf("  Name: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj1.GetString("description"); err == nil {
		log.Printf("  Description: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	log.Printf("Entity 2")
	c2, err := b2.Cid().StringOfBase('z')
	if err != nil {
		log.Panicf("Cannot retrieve CID from block: %s", err)
	}
	log.Printf("  CID: %s", c2)

	log.Printf("  Raw data: %s", b2.RawData())

	log.Printf("  Type: %s", obj2.GetName())
	log.Printf("  Schema version: %d", obj2.GetVersion())

	if val, err := obj2.GetString("id"); err == nil {
		log.Printf("  ID: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if _, err := obj2.GetString("name"); err == nil {
		log.Panic("Should not have property \"name\"")
	}

	if _, err := obj2.GetString("description"); err == nil {
		log.Panic("Should not have property \"description\"")
	}

	log.Printf("Entity 3")
	c3, err := b3.Cid().StringOfBase('z')
	if err != nil {
		log.Panicf("Cannot retrieve CID from block: %s", err)
	}
	log.Printf("  CID: %s", c3)

	log.Printf("  Raw data: %s", b3.RawData())

	log.Printf("  Type: %s", obj3.GetName())
	log.Printf("  Schema version: %d", obj3.GetVersion())

	if val, err := obj3.GetString("id"); err == nil {
		log.Printf("  ID: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj3.GetString("name"); err == nil {
		log.Printf("  Name: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if _, err := obj3.GetString("description"); err == nil {
		log.Panic("Should not have property \"description\"")
	}

	// --------------------------------------------------
	// JSON

	log.Println()
	log.Println("********************************************************************************")
	log.Println("Entity JSON")
	log.Println("********************************************************************************")

	log.Printf("Entity 1")
	json1, err := obj1.MarshalJSON()
	if err != nil {
		log.Panicf("Cannot marshal JSON: %s", err)
	}
	log.Println(string(json1))
	log.Println(string(pretty.Pretty([]byte(json1))))

	log.Printf("Entity 2")
	json2, err := obj2.MarshalJSON()
	if err != nil {
		log.Panicf("Cannot marshal JSON: %s", err)
	}
	log.Println(string(json2))
	log.Println(string(pretty.Pretty([]byte(json2))))

	log.Printf("Entity 3")
	json3, err := obj3.MarshalJSON()
	if err != nil {
		log.Panicf("Cannot marshal JSON: %s", err)
	}
	log.Println(string(json3))
	log.Println(string(pretty.Pretty([]byte(json3))))

	return []iscn.IscnObject{b1, b2, b3}
}
