# Api Test

## TODO
- Agregar autorización (roles)
- Logica para hace los filtros contra la db (filtros de busqueda)
- Agregar posthog/datadog/newrelic/sentry/etc
- Poner el https://golangci-lint.run/


## Braille ASCII Art

> https://lachlanarthur.github.io/Braille-ASCII-Art/

## Folder Structure

saas/
├── cmd/                   # Puntos de entrada
│   ├── api/               # API principal
│   ├── worker/            # Procesos en segundo plano
│   ├── cli/               # Comandos CLI
│   └── migrate/           # Migraciones de base de datos
├── internal/              # Código interno
│   ├── app/               # Inicialización de la aplicación (configuración)
│   ├── common/            # Componentes reutilizables (logs, auth, middleware)
│   ├── config/            # Configuración central
│   ├── database/          # Configuración y conexión a DB
│   ├── modules/           # Módulos de negocio
│   │   ├── inventario/
│   │   │   ├── domain/    # Definición de entidades y lógica pura
│   │   │   ├── infra/     # Implementaciones concretas (repos, API)
│   │   │   ├── usecase/   # Casos de uso
│   │   │   ├── api/       # Controladores HTTP/gRPC
│   │   │   ├── events/    # Mensajería
│   │   │   └── service.go # Wiring del módulo
├── pkg/                   # Librerías reutilizables
├── scripts/               # Scripts de despliegue y mantenimiento
└── go.mod


