package main

import (
	"bytes"
	"context"
	"log"
	"math/rand"

	"github.com/btcsuite/btcutil/base58"
	"github.com/tidwall/pretty"

	icore "github.com/ipfs/interface-go-ipfs-core"
	iscn "github.com/likecoin/iscn-ipld/plugin/block"
)

func testIscnKernel(
	ctx context.Context,
	ipfs icore.CoreAPI,
	stakeholders iscn.IscnObject,
	content iscn.IscnObject,
) {
	// --------------------------------------------------
	log.Printf("Generating ISCN kernel block ...")

	id := make([]byte, 32)
	rand.Read(id)

	data := map[string]interface{}{
		"id":           id,
		"timestamp":    "2020-01-01T12:34:56Z",
		"version":      1,
		"stakeholders": stakeholders.Cid(),
		"content":      content.Cid(),
		"zzz":          -987654321,
		"yyy":          []string{"abc", "def", "ghi"},
		"xxx":          []byte{'x', 'y', 'z'},
		"p": map[string]interface{}{
			"a": 10,
			"b": map[string]interface{}{
				"ba": "abc",
				"bb": 123,
			},
		},
	}

	b, err := iscn.Encode(iscn.CodecISCN, 1, data)
	if err != nil {
		log.Panicf("Cannot create ISCN kernel block: %s", err)
	}

	// --------------------------------------------------
	log.Printf("Pinning ISCN kernel block ...")

	if err := ipfs.Dag().Pinning().Add(ctx, b); err != nil {
		log.Panicf("Cannot pin IPLD: %s", err)
	}

	// --------------------------------------------------
	log.Printf("Getting ISCN kernel block ...")

	ret, err := ipfs.Dag().Get(ctx, b.Cid())
	if err != nil {
		log.Panicf("Cannot fetch IPLD: %s", err)
	}

	obj, err := iscn.Decode(ret.RawData(), b.Cid())
	if err != nil {
		log.Panicf("Cannot decode IPLD raw data: %s", err)
	}

	// --------------------------------------------------
	// Report
	log.Println("********************************************************************************")
	log.Println("ISCN kernel report")
	log.Println("********************************************************************************")

	c, err := b.Cid().StringOfBase('z')
	if err != nil {
		log.Panicf("Cannot retrieve CID from block: %s", err)
	}
	log.Printf("  CID: %s", c)

	log.Printf("  Raw data: %s", b.RawData())

	log.Printf("  Type: %s", obj.GetName())
	log.Printf("  Schema version: %d", obj.GetVersion())

	if val, err := obj.GetBytes("id"); err == nil {
		if !bytes.Equal(id, val) {
			log.Panic("ID is not matched")
		}

		log.Printf("  ID (original): %s", base58.Encode(id))
		log.Printf("  ID           : %s", base58.Encode(val))
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj.GetString("timestamp"); err == nil {
		log.Printf("  Timestamp: %q", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj.GetUint64("version"); err == nil {
		log.Printf("  Version: %d", val)
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj.GetCid("stakeholders"); err == nil {
		c, err := val.StringOfBase('z')
		if err != nil {
			log.Panicf("Cannot retrieve CID from block: %s", err)
		}
		log.Printf("  Stakeholders: %s (0x%x)", c, val.Type())
	} else {
		log.Panicf("%s", err)
	}

	if val, err := obj.GetCid("content"); err == nil {
		c, err := val.StringOfBase('z')
		if err != nil {
			log.Panicf("Cannot retrieve CID from block: %s", err)
		}
		log.Printf("  Content: %s (0x%x)", c, val.Type())
	} else {
		log.Panicf("%s", err)
	}

	log.Println("  Custom properties:")
	for key, value := range obj.GetCustom() {
		log.Printf("    %q:", key)
		log.Printf("      %T -> %v", value, value)
	}

	// --------------------------------------------------
	// JSON

	log.Println()
	log.Println("********************************************************************************")
	log.Println("ISCN kernel JSON")
	log.Println("********************************************************************************")

	json, err := obj.MarshalJSON()
	if err != nil {
		log.Panicf("Cannot marshal JSON: %s", err)
	}
	log.Println(string(json))
	log.Println(string(pretty.Pretty([]byte(json))))
}
