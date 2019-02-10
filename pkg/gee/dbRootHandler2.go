package gee

import (
	"encoding/json"
	"fmt"
	"net/http"

	khdb "AGES/pkg/gee/keyhole_dbroot"

	"github.com/golang/protobuf/proto"
)

const minConfig = `
{
	"imagery_present": true,
	"proto_imagery": true,
	"provider_info": [
	  {
		"provider_id": 1,
		"copyright_string": {
		  "value": "Imagery © 2005 USGS"
		},
		"vertical_pixel_offset": -1
	  }
	],
	"end_snippet": {
	  "model": {
		"radius": 6371.01,
		"negative_altitude_exponent_bias": 32,
		"compressed_negative_altitude_threshold": 1e-12
	  },
	  "disable_authentication": true
	},
	"database_version": {
	  "quadtree_version": 5
	}
  }
`

//DBRootHandler2 returns a dbRoot object from scratch
func DBRootHandler2(w http.ResponseWriter, r *http.Request) {
	//get DbRoot json data
	drp := &khdb.DbRootProto{}
	json.Unmarshal([]byte(minConfig), drp)
	//convert to protobuf
	drpBytes, err := proto.Marshal(drp)
	if err != nil {
		fmt.Fprintln(w, "drp proto")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//compress
	fmt.Printf("UL uncompressed %d\n", len(drpBytes))
	cDrp, err := compressPacket(drpBytes)
	drpBytes, _ = uncompressPacket(cDrp)
	cDrp, err = compressPacket(drpBytes)

	fmt.Printf("UL compressed %d\n", len(cDrp))
	if err != nil {
		fmt.Fprintln(w, "compress")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// encrypt / obfuscate
	XOR(cDrp, []byte(defaultKey), false)
	//embed DbRoot payload in EncryptedDbRoot
	ec := khdb.EncryptedDbRootProto_ENCRYPTION_XOR
	edrp2 := &khdb.EncryptedDbRootProto{
		EncryptionType: &ec,
		EncryptionData: []byte(defaultKey),
		DbrootData:     cDrp,
	}
	//convert to protobuf
	edrpBytes, err := proto.Marshal(edrp2)
	if err != nil {
		fmt.Fprintf(w, "edrp proto\n%+v\n%v", edrp2, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//send bytes
	w.Write(edrpBytes)
}
