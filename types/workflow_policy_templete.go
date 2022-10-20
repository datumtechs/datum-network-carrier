package types

//[
//  {
//    "origin": "t2",
//    "reference": [
//      {
//        "target": "t1",
//        "dependPartyId": [
//          "result1",
//          "result2"
//        ],
//        "merge": false,
//        "dependParams": [
//          "{\"inputType\": 1, \"keyColumnName\": \"\", \"selectedColumnNames\": []}",
//          "{\"inputType\": 1, \"keyColumnName\": \"\", \"selectedColumnNames\": []}"
//        ]
//      }
//    ]
//  },
//  {
//    "origin": "t3",
//    "reference": [
//      {
//        "target": "t2",
//        "dependPartyId": [
//          "result1",
//          "result2"
//        ],
//        "merge": false,
//        "dependParams": [
//          "{\"inputType\": 3}"
//        ]
//      }
//    ]
//  }
//]

type ReferTo struct {
	Target           string   `json:"target"`
	DependPartyId    []string `json:"dependPartyId"`
	Merge            bool     `json:"merge"`
	DependParamsType []uint32 `json:"dependParamsType"`
	DependParams     []string `json:"dependParams"`
}
type Relation struct {
	Origin    string    `json:"origin"`
	Reference []ReferTo `json:"reference"`
}
type WorkflowPolicy []Relation

func (s WorkflowPolicy) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s WorkflowPolicy) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s WorkflowPolicy) Less(i, j int) bool { return s[i].Origin < s[j].Origin }
