// Package tag provides a lightweight key-value tagging system for port events.
//
// Tags are immutable sets of "key=value" string pairs that can be attached to
// alert events to carry contextual metadata — such as environment, host role,
// datacenter, or any custom annotation — through the portwatch notification
// pipeline.
//
// # Usage
//
//	s, err := tag.New([]string{"env=prod", "role=web", "dc=us-east"})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	env, _ := s.Get("env")  // "prod"
//	fmt.Println(s.String()) // "dc=us-east,env=prod,role=web"
//
// Sets are immutable; Merge returns a new Set without modifying either operand.
package tag
