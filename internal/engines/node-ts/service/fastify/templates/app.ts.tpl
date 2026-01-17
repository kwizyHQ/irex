import fastify, { FastifyInstance } from 'fastify'

export interface AppConfig {
  port?: number
  host?: string
  logger?: { level?: string }
}

export type StartHook = (app: FastifyInstance) => Promise<void> | void
export type StopHook = (app: FastifyInstance) => Promise<void> | void

const DEFAULT_CONFIG: Required<AppConfig> = {
  port: process.env.PORT ? parseInt(process.env.PORT) : {{ .EnvPort }},
  host: process.env.HOST || "{{ .EnvHost }}",
  logger: { level: process.env.LOG_LEVEL || 'info' }
}

let app: FastifyInstance | null = null
let serverConfig: Required<AppConfig> = DEFAULT_CONFIG
let startHooks: StartHook[] = []
let stopHooks: StopHook[] = []

export function registerStartHook(h: StartHook) {
  startHooks.push(h)
}

export function registerStopHook(h: StopHook) {
  stopHooks.push(h)
}

export function buildApp(config?: AppConfig) {
  serverConfig = { ...DEFAULT_CONFIG, ...(config || {}) } 
  app = fastify({ logger: serverConfig.logger })

  // placeholder for route registration - engine will replace or extend
  try {
    // dynamic import of routes index if it exists in generated project
    // note: template consumers can override this behavior by calling registerStartHook
    // to register custom routes before starting the server.
    // eslint-disable-next-line @typescript-eslint/no-var-requires
    // routes are expected to be an array of register functions: () => plugin
    // We keep this non-failing so generated projects without routes still start.
    // The template engine that scaffolds routes should write a routes/index file.
    // @ts-ignore: optional module
    const routes = require('./routes/index').default
    if (Array.isArray(routes)) {
      routes.forEach((r: any) => app!.register(r))
    }
  } catch (err) {
    // ignore missing routes
  }

  return app
}

export async function start(config?: AppConfig) {
  if (!app) buildApp(config)
  if (!app) throw new Error('app not built')

  for (const h of startHooks) {
    await Promise.resolve(h(app))
  }

  try {
    await app.listen({ port: serverConfig.port, host: serverConfig.host })
    app.log.info(`Server listening at ${serverConfig.host}:${serverConfig.port}`)
  } catch (err) {
    app.log.error(err)
    process.exit(1)
  }
}

export async function stop() {
  if (!app) return

  for (const h of stopHooks) {
    await Promise.resolve(h(app))
  }

  try {
    await app.close()
    app = null
  } catch (err) {
    // ignore close errors
  }
}

export default { buildApp, start, stop, registerStartHook, registerStopHook }
import fastify from 'fastify'
import { config } from './config'
import routes from './routes/index'

export function buildApp() {
  const app = fastify({ logger: config.logger })

  // register routes
  routes.forEach(r => app.register(r))

  return app
}

export default buildApp
