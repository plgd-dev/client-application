import { WebTracerProvider } from '@opentelemetry/sdk-trace-web'
import {
  ConsoleSpanExporter,
  SimpleSpanProcessor,
} from '@opentelemetry/sdk-trace-base'
import { registerInstrumentations } from '@opentelemetry/instrumentation'
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch'
import { ZoneContextManager } from '@opentelemetry/context-zone'
import { B3Propagator } from '@opentelemetry/propagator-b3'

let webTracer = undefined

const init = () => {
  const provider = new WebTracerProvider()
  provider.addSpanProcessor(new SimpleSpanProcessor(new ConsoleSpanExporter()))
  provider.register({
    contextManager: new ZoneContextManager(),
    propagator: new B3Propagator(),
  })

  registerInstrumentations({
    instrumentations: [
      new FetchInstrumentation({
        clearTimingResources: true,
      }),
    ],
  })

  webTracer = provider.getTracer('client-app-tracer')
}

export const openTelemetry = {
  init: () => init(),
  getWebTracer: () => webTracer,
}
