package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	khdb "AGES/pkg/gee/keyhole_dbroot"
	"AGES/pkg/net"

	"github.com/golang/protobuf/proto"
)

//DBRootHandler returns a dbRoot object
func DBRootHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("config/dbRoot.raw"); os.IsNotExist(err) {
		net.DownloadFile("config/dbRoot.raw", "http://www.earthenterprise.org/3d/dbRoot.v5")
	}
	if _, err := os.Stat("config/dbRoot.js"); os.IsNotExist(err) {
		b := readFile("config/dbRoot.raw")
		edrp := khdb.EncryptedDbRootProto{}
		unProto(b, &edrp)
		//unobfuscate
		XOR(edrp.DbrootData, edrp.EncryptionData, true)
		//uncompress
		//fmt.Printf("DL uncompressed %d\n", len(edrp.DbrootData))
		dbRoot, _ := uncompressPacket(edrp.DbrootData)

		//fmt.Printf("DL compressed %d\n", len(dbRoot))
		//read the protocol buffer
		drp := khdb.DbRootProto{}
		unProto(dbRoot, &drp)
		//show the content
		b, err := json.MarshalIndent(drp, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		writeFile("config/dbRoot.js", b)
		edrp.DbrootData = nil
		e, err := json.MarshalIndent(edrp, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		writeFile("config/encDbRoot.js", e)
	}

	//get DbRoot json data
	drp := &khdb.DbRootProto{}
	err := unMarshalJSONFile("config/dbRoot.js", drp)
	if err != nil {
		fmt.Fprintln(w, "drp json")
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	//convert to protobuf
	drpBytes, err := proto.Marshal(drp)
	if err != nil {
		fmt.Fprintln(w, "drp proto")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//compress
	//fmt.Printf("UL uncompressed %d\n", len(drpBytes))
	cDrp, err := compressPacket(drpBytes)
	drpBytes, _ = uncompressPacket(cDrp)
	cDrp, err = compressPacket(drpBytes)

	//fmt.Printf("UL compressed %d\n", len(cDrp))
	if err != nil {
		fmt.Fprintln(w, "compress")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// encrypt / obfuscate
	XOR(cDrp, []byte(defaultKey), false)
	//get EncryptedDbRoot json data
	edrp2 := &khdb.EncryptedDbRootProto{}
	err = unMarshalJSONFile("config/encDbRoot.js", edrp2)
	if err != nil {
		fmt.Fprintln(w, "edrp json")
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	//embed DbRoot payload in EncryptedDbRoot
	edrp2.DbrootData = cDrp
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
