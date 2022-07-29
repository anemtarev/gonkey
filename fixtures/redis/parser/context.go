package parser

type context struct {
	keyRefs map[string]Keys
	mapRefs map[string]MapRecordValue
	setRefs map[string]SetRecordValue
}

func NewContext() *context {
	return &context{
		keyRefs: make(map[string]Keys),
		mapRefs: make(map[string]MapRecordValue),
		setRefs: make(map[string]SetRecordValue),
	}
}
