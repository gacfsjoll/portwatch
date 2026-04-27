// Package labelset provides structured key=value label sets that can be
// attached to port-change events before they are dispatched through the
// alert pipeline.
//
// Labels allow operators to annotate events with contextual metadata such
// as environment, team ownership, or service tier, making downstream
// routing and filtering easier.
//
// Example usage:
//
//	ls, err := labelset.New("env=prod", "owner=platform")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(ls.String()) // env=prod,owner=platform
package labelset
