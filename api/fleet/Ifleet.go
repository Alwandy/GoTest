package fleet

// Armament structure, needed another database table to create the armament types and use relationship to assign to a ship armements existing
// Armaments available:
/* "armament": [
	{
		"title": "Turbo Laser",
		"qty": "60"
}, {
	  "title": "Ion Cannons",
	  "qty": "60",
	},
	{
	  "title": "Tractor Beam",
	  "qty": "10",
}, ]
*/
type Armament struct {
	Title string `json:"title"`
	Qty   string `json:"qty"`
}

/*
    The structure to generate a ship
	Ship{
		Name:   "",
		Image:  "",
		Class:  "",
		Crew:   0,
		Status: "",
		Value:  0.0,
		//Armament: nil, - This is disabled, didn't have time to create another table for armament and use ids to add them in
	}

*/
type Ship struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Class string `json:"class"`
	Crew  int    `json:"crew"`
	//Armament map[int]Armament `json:"armament"`
	Image  string  `json:"image"`
	Value  float64 `json:"value"`
	Status string  `json:"status"`
}
