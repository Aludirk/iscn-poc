package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ipfs/go-cid"
	"github.com/tidwall/pretty"

	icore "github.com/ipfs/interface-go-ipfs-core"
	iscn "github.com/likecoin/iscn-ipld/plugin/block"
)

func testRights(
	ctx context.Context,
	ipfs icore.CoreAPI,
	entities []iscn.IscnObject,
) iscn.IscnObject {
	log.Printf("Generating rights block ...")

	termCid, err := cid.Decode("Qmacpqc7EWQBU9q8cctAj1hdoVXdyMH7Geq7FcpZ8XA5M8")
	if err != nil {
		log.Panicf("Cannot create a CID for ISCN kernel: %s", err)
	}

	data := map[string]interface{}{
		"rights": []map[string]interface{}{
			{
				"holder": entities[0].Cid(),
				"type":   "license",
				"terms":  termCid,
			},
			{
				"holder": entities[0].Cid(),
				"type":   "license",
				"terms":  termCid,
				"period": map[string]interface{}{
					"from": "2020-01-01T12:34:56Z",
					"to":   "2046-01-01T12:34:56+08:00",
				},
				"territory": "Mars",
			},
			{
				"holder": entities[1].Cid(),
				"type":   "license",
				"terms":  termCid,
				"period": map[string]interface{}{
					"from": "2020-01-01T12:34:56Z",
				},
			},
			{
				"holder": entities[1].Cid(),
				"type":   "license",
				"terms":  termCid,
				"period": map[string]interface{}{
					"to": "2046-01-01T12:34:56+08:00",
				},
			},
			{
				"holder":    entities[2].Cid(),
				"type":      "license",
				"terms":     termCid,
				"territory": "Jupiter",
			},
		},
	}

	b, err := iscn.Encode(iscn.CodecRights, 1, data)
	if err != nil {
		log.Panicf("Cannot create rights block: %s", err)
	}
	log.Printf("New rights block %s", b.RawData())

	// --------------------------------------------------
	log.Printf("Pinning rights blocks ...")

	if err := ipfs.Dag().Pinning().Add(ctx, b); err != nil {
		log.Panicf("Cannot pin IPLD: %s", err)
	}

	// --------------------------------------------------
	log.Printf("Getting rights block ...")

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
	log.Println("Rights report")
	log.Println("********************************************************************************")

	c, err := b.Cid().StringOfBase('z')
	if err != nil {
		log.Panicf("Cannot retrieve CID from block: %s", err)
	}
	log.Printf("  CID: %s", c)

	log.Printf("  Raw data: %s", b.RawData())

	log.Printf("  Type: %s", obj.GetName())
	log.Printf("  Schema version: %d", obj.GetVersion())

	if rights, err := obj.GetArray("rights"); err == nil {
		for i, right := range rights {
			r, ok := right.(iscn.IscnObject)
			if !ok {
				log.Panicf("(Index %d) Right is not an \"IscnObject\"", i)
			}

			log.Printf("  Right %d -", i+1)
			log.Printf("    Type: %s", r.GetName())
			log.Printf("    Schema version: %d", obj.GetVersion())

			if val, err := r.GetCid("holder"); err == nil {
				c, err := val.StringOfBase('z')
				if err != nil {
					log.Panicf("Cannot retrieve CID from block: %s", err)
				}
				log.Printf("    Holder: %s (0x%x)", c, val.Type())
			} else {
				log.Panicf("%s", err)
			}

			if val, err := r.GetString("type"); err == nil {
				log.Printf("    Type: %q", val)
			} else {
				log.Panicf("%s", err)
			}

			if val, err := r.GetCid("terms"); err == nil {
				c, err := val.StringOfBase('z')
				if err != nil {
					log.Panicf("Cannot retrieve CID from stakeholder: %s", err)
				}
				log.Printf("    Terms: %s (0x%x)", c, val.Type())
			} else {
				log.Panicf("%s", err)
			}

			if val, err := r.GetObject("period"); err == nil {
				period, ok := val.(iscn.IscnObject)
				if !ok {
					log.Panic("Time period should be an \"IscnObject\"")
				}

				log.Println("    Period -")

				if val, err := period.GetString("from"); err == nil {
					log.Printf("      From: %q", val)
				} else {
					if err.Error() != fmt.Sprintf("%q is not found", "from") {
						log.Panicf("%s", err)
					}
				}

				if val, err := period.GetString("to"); err == nil {
					log.Printf("      To: %q", val)
				} else {
					if err.Error() != fmt.Sprintf("%q is not found", "to") {
						log.Panicf("%s", err)
					}
				}
			} else {
				if err.Error() != fmt.Sprintf("%q is not found", "period") {
					log.Panicf("%s", err)
				}
			}

			if val, err := r.GetString("territory"); err == nil {
				log.Printf("    Territory: %q", val)
			} else {
				if err.Error() != fmt.Sprintf("%q is not found", "territory") {
					log.Panicf("%s", err)
				}
			}
		}
	} else {
		log.Panicf("%s", err)
	}

	// --------------------------------------------------
	// JSON

	log.Println()
	log.Println("********************************************************************************")
	log.Println("Rights JSON")
	log.Println("********************************************************************************")

	json, err := obj.MarshalJSON()
	if err != nil {
		log.Panicf("Cannot marshal JSON: %s", err)
	}
	log.Println(string(json))
	log.Println(string(pretty.Pretty([]byte(json))))

	return b
}
