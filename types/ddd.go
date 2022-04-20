package types

import "encoding/json"

type AParams struct {  // type 1

}


type BParams struct {  // type 2

}




func A (dyparams string) string {

	var dym map[string]interface{}

	json.Unmarshal([]byte(dyparams), &dym)

	m := make(map[string]interface{}, 0)

	m["self_params"] = AParams{} // m["self_params"] = BParams{}
	m["common_params"] = dym

	b, _ := json.Marshal(m)
	return string(b)
}