package types

type TaskMessage struct {
	TaskHash  		string    				`json:"taskHash"`
	Owner       	*TaskParticipant   		`json:"owner"`
	Partners 		[]*TaskParticipant  	`json:"partners"`
	Receivers       []*TaskResultReceiver  	`json:"receivers"`
	ContractCode 	string 					`json:"contractCode"`
	Extra 			string  				`json:"extra"`
}

type TaskParticipant struct {
	NodeId    		string  				`json:"nodeId"`
	IdentityId 		string 					`json:"identityId"`
	Alias 			string 					`json:"alias"`
	MetaData 		*ParticipantMetaData 	`json:"metaData"`
}

type ParticipantMetaData struct {
	MetaId  			string  		`json:"metaId"`
	FilePath 			string 			`json:"filePath"`
	CindexList 			[]uint32		`json:"cindexList"`
}

type TaskResultReceiver struct {
	NodeId   		string   				`json:"nodeId"`
	IdentityId  	string 					`json:"identityId"`
	Providers       []*TaskResultProvider 	`json:"providers"`
}

type TaskResultProvider struct {
	NodeId    		string  				`json:"nodeId"`
	IdentityId 		string 					`json:"identityId"`
	Alias 			string 					`json:"alias"`
}