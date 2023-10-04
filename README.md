# VoxelAI

Generate voxel art from a prompt.

# Setup

**Install**:
- [flyctl](https://fly.io/docs/flyctl)
- [Turso CLI](https://docs.turso.tech/reference/turso-cli)
- [sqld](https://github.com/libsql/sqld)

**Run:**
- `pnpm i`

## Local development

- Create a `.env.local` file an add:
  - `DB_URL=<local database URL>` (most likely `http://127.0.0.1:1331`)
- Start database & bundlers: `pnpm run dev` (this runs `pnpm run dev:db` and `pnpm run dev:sources`. If you're on Windows, you might want to run `dev:db` separately on WSL)
- Start web server: `go run . -debug`

**Note:** Server will bind to 127.0.0.1:8080 by default unless `BIND` and
`PORT` environment variables are set.

## Deployment

Configurations are provided for deploying to a staging and production
environment on Fly.io.

Change `app` in `fly.dev.toml` and `fly.prod.toml` to you own app's name.

**Staging:**

**Important:** Make sure no public IPs are allocated to this service.

Deploy: `pnpm run deploy:staging`

Assuming WireGuard VPN is connected to Fly.io's tunnel, the app should be
accessible at `http://<app-name>.flycast/`
