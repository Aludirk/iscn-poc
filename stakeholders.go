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

func testStakeholders(
	ctx context.Context,
	ipfs icore.CoreAPI,
	entities []iscn.IscnObject,
) iscn.IscnObject {
	log.Printf("Generating stakeholders block ...")

	kernelCid, err := cid.Decode("z4gAY85gBq5PF1xydzdg6wgW9Q88B7B5bu1LYD7AAmRxWnpjFGQ")
	if err != nil {
		log.Panicf("Cannot create a CID for ISCN kernel: %s", err)
	}

	data := map[string]interface{}{
		"stakeholders": []map[string]interface{}{
			{
				"type":        "Creator",
				"stakeholder": entities[0].Cid(),
				"sharing":     8,
			},
			{
				"type":        "FootprintStakeholder",
				"stakeholder": entities[1].Cid(),
				"sharing":     1,
				"footprint":   kernelCid,
			},
			{
				"type":        "FootprintStakeholder",
				"stakeholder": entities[2].Cid(),
				"sharing":     1,
				"footprint":   "https://example.com/footprint.html",
			},
		},
	}

	b, err := iscn.Encode(iscn.CodecStakeholders, 1, data)
	if err != nil {
		log.Panicf("Cannot create stakeholders block: %s", err)
	}
	log.Printf("New stakeholders block %s", b.RawData())

	// --------------------------------------------------
	log.Printf("Pinning stakeholders blocks ...")

	if err := ipfs.Dag().Pinning().Add(ctx, b); err != nil {
		log.Panicf("Cannot pin IPLD: %s", err)
	}

	// --------------------------------------------------
	log.Printf("Getting stakeholders block ...")

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
	log.Println("Stakeholders report")
	log.Println("********************************************************************************")

	c, err := b.Cid().StringOfBase('z')
	if err != nil {
		log.Panicf("Cannot retrieve CID from block: %s", err)
	}
	log.Printf("  CID: %s", c)

	log.Printf("  Raw data: %s", b.RawData())

	log.Printf("  Type: %s", obj.GetName())
	log.Printf("  Schema version: %d", obj.GetVersion())

	if stakeholders, err := obj.GetArray("stakeholders"); err == nil {
		for i, stakeholder := range stakeholders {
			s, ok := stakeholder.(iscn.IscnObject)
			if !ok {
				log.Panicf("(Index %d) Stakeholder is not an \"IscnObject\"", i)
			}

			log.Printf("  Stakeholder %d -", i+1)
			log.Printf("    Type: %s", s.GetName())
			log.Printf("    Schema version: %d", obj.GetVersion())

			if val, err := s.GetString("type"); err == nil {
				log.Printf("    Type: %q", val)
			} else {
				log.Panicf("%s", err)
			}

			if val, err := s.GetCid("stakeholder"); err == nil {
				c, err := val.StringOfBase('z')
				if err != nil {
					log.Panicf("Cannot retrieve CID from stakeholder: %s", err)
				}
				log.Printf("    Stakeholder: %s (0x%x)", c, val.Type())
			} else {
				log.Panicf("%s", err)
			}

			if val, err := s.GetUint32("sharing"); err == nil {
				log.Printf("    Sharing: %d", val)
			} else {
				log.Panicf("%s", err)
			}

			if cidObj, url, err := s.GetLink("footprint"); err == nil {
				if cidObj.Defined() {
					c, err := cidObj.StringOfBase('z')
					if err != nil {
						log.Panicf("Cannot retrieve CID from footprint: %s", err)
					}
					log.Printf("    Footprint: %s (0x%x)", c, cidObj.Type())
				} else if len(url) > 0 {
					log.Printf("    Footprint: %q", url)
				} else {
					log.Panic("Unknown error from `GetLink`: not a Cid nor URL")
				}
			} else {
				if err.Error() != fmt.Sprintf("%q is not found", "footprint") {
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
	log.Println("Stakeholders JSON")
	log.Println("********************************************************************************")

	json, err := obj.MarshalJSON()
	if err != nil {
		log.Panicf("Cannot marshal JSON: %s", err)
	}
	log.Println(string(json))
	log.Println(string(pretty.Pretty([]byte(json))))

	return b
}
