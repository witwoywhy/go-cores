package apps

import "net/http"

func MaskHeader(header http.Header) {
	for k := range header {
		if _, ok := HeaderMaskingList[k]; ok {
			header.Set(k, "MASKING")
		}
	}
}
