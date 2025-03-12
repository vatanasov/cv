package models

import "strings"

type TechnologyMeta struct {
	Pretty string
	Order  int
	Name   string
}
type PrettifiedTechnologies struct {
	internal map[string]TechnologyMeta
	aliases  map[string]string
}

var PrettyTechnologies = PrettifiedTechnologies{
	internal: map[string]TechnologyMeta{
		"elixir": {Pretty: "<i class=\"devicon-elixir-plain colored\"></i>Elixir", Order: 9, Name: "elixir"},
		"golang": {
			Pretty: "<i class=\"devicon-go-original-wordmark colored\"></i>Go",
			Order:  9,
			Name:   "golang",
		},
		"phoenix": {
			Pretty: "<i class=\"devicon-phoenix-original colored\"></i>Phoenix",
			Order:  8,
			Name:   "phoenix",
		},
		"phoenix liveview": {
			Pretty: "<i class=\"devicon-phoenix-original colored\"></i>Phoenix LiveView",
			Order:  7,
			Name:   "phoenix liveview",
		},
		"gcp": {
			Pretty: "<i class=\"devicon-googlecloud-plain colored\"></i>Google Cloud",
			Order:  1,
			Name:   "gcp",
		},
		"graphql": {
			Pretty: "<i class=\"devicon-graphql-plain colored\"></i>GraphQL",
			Order:  1,
			Name:   "graphql",
		},
		"ocaml": {Pretty: "<i class=\"devicon-ocaml-plain colored\"></i>Ocaml", Order: 3, Name: "ocaml"},
		"react": {Pretty: "<i class=\"devicon-react-original colored\"></i>React", Order: 6, Name: "react"},
		"kafka": {
			Pretty: "<i class=\"devicon-apachekafka-original colored\"></i>Kafka",
			Order:  1,
			Name:   "kafka",
		},
		"python": {Pretty: "<i class=\"devicon-python-plain colored\"></i>Python", Order: 8, Name: "python"},
		"typescript": {
			Pretty: "<i class=\"devicon-typescript-plain colored\"></i>TypeScript",
			Order:  8,
			Name:   "typescript",
		},
		"node.js": {
			Pretty: "<i class=\"devicon-nodejs-plain-wordmark colored\"></i>NodeJS",
			Order:  7,
			Name:   "node.js",
		},
		"mongodb": {
			Pretty: "<i class=\"devicon-mongodb-plain colored\"></i>MongoDB",
			Order:  5,
			Name:   "mongodb",
		},
		"postgresql": {
			Pretty: "<i class=\"devicon-postgresql-plain colored\"></i>PostgreSQL",
			Order:  6,
			Name:   "postgresql",
		},
		"django": {Pretty: "<i class=\"devicon-django-plain colored\"></i>Django", Order: 7, Name: "django"},
		"dynamodb": {
			Pretty: "<i class=\"devicon-dynamodb-plain colored\"></i>DynamoDB",
			Order:  1,
			Name:   "dynamodb",
		},
		"sqs": {Pretty: "SQS", Order: 1, Name: "sqs"},
		"rds": {Pretty: "RDS", Order: 1, Name: "rds"},
		"aws": {
			Pretty: "<i class=\"devicon-amazonwebservices-plain-wordmark colored\"></i>AWS",
			Order:  1,
			Name:   "aws",
		},
	},
	aliases: map[string]string{
		"go": "golang",
	},
}

func (t PrettifiedTechnologies) Get(name string) TechnologyMeta {
	lowercasedName := strings.ToLower(name)
	if alias, ok := t.aliases[lowercasedName]; ok {
		lowercasedName = alias
	}

	value, ok := t.internal[lowercasedName]
	if !ok {
		return TechnologyMeta{Pretty: lowercasedName, Order: 0, Name: lowercasedName}
	}
	return value
}
