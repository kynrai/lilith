package postgres_example

var getThing = `
SELECT * FROM things where id=$1;
`

var setThing = `
INSERT INTO things
(id, data)
VALUES ($1, $2)
ON CONFLICT (id)
DO UPDATE SET data = $3;
`
