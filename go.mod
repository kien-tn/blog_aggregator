module github.com/kien-tn/blog_aggregator

go 1.23.3

replace github.com/kien-tn/blog_aggregator/internal => ./internal

require (
	github.com/google/uuid v1.6.0
	github.com/kien-tn/blog_aggregator/internal v0.0.0-00010101000000-000000000000
)

require github.com/lib/pq v1.10.9
