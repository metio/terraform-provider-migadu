/*
 * SPDX-FileCopyrightText: The terraform-provider-migadu Authors
 * SPDX-License-Identifier: 0BSD
 */

package simulator

import (
	"encoding/json"
	"net/http"
	"testing"
)

func writeJsonResponse(t *testing.T, w http.ResponseWriter, value any) {
	bytes, err := json.Marshal(value)
	if err != nil {
		t.Errorf("Could not marshall data")
	}
	_, err = w.Write(bytes)
	if err != nil {
		t.Errorf("Could not write data")
	}
}
