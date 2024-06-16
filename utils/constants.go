package utils

var PathApproach = map[string]interface{}{
	"Script":  nil,
	"Static":  nil,
	"Dynamic": nil,
	"Mobile": map[string]interface{}{
		"Native": nil,
		"Hybrid": nil,
	},
	"Full": map[string]interface{}{
		"Client-side": nil,
		"Server-side": nil,
	},
}

var CodeRulesArray = []string{"Modular", "Clean Code", "Easy to Read", "Easy to Customize", "Basic Security (Validation)"}
