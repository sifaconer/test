package postgres

import "embed"

//go:embed migrations/admin/*.sql
var Admins embed.FS

//go:embed migrations/tenants/*.sql
var Tenants embed.FS

//go:embed migrations/common/*.sql
var Common embed.FS
