import { useEffect, useState } from 'react'
import { context, trace } from '@opentelemetry/api'

import { useIsMounted } from '@/common/hooks'
import { fetchApi, streamApi } from '@/common/services'
import { useAppConfig } from '@/containers/app'

export const useApi = (url, options = {}) => {
  const isMounted = useIsMounted()
  const [state, setState] = useState({
    error: null,
    loading: true,
    data: null,
  })
  const [refreshIndex, setRefreshIndex] = useState(0)
  const { telemetryWebTracer } = useAppConfig()
  const { telemetrySpan } = options

  useEffect(
    () => {
      ;(async () => {
        try {
          // Set loading to true
          setState({ ...state, loading: true })
          let data = undefined

          if (telemetryWebTracer && telemetrySpan) {
            const singleSpan = telemetryWebTracer.startSpan(telemetrySpan)
            await context.with(
              trace.setSpan(context.active(), singleSpan),
              async () => {
                data = await fetchApi(url, options).then(result => {
                  trace
                    .getSpan(context.active())
                    .addEvent('fetching-single-span-completed')
                  singleSpan.end()

                  return result.data
                })
              }
            )
          } else {
            const result = await fetchApi(url, options)
            data = result.data
          }

          if (isMounted.current) {
            setState({
              ...state,
              data,
              error: null,
              loading: false,
            })
          }
        } catch (error) {
          if (isMounted.current) {
            setState({
              ...state,
              data: null,
              error,
              loading: false,
            })

            trace
              .getSpan(context.active())
              .addEvent('fetching-single-span-completed')
          }
        }
      })()
    },
    [url, refreshIndex] // eslint-disable-line
  )

  return {
    ...state,
    updateData: updatedData => setState({ ...state, data: updatedData }),
    refresh: () => setRefreshIndex(refreshIndex + 1),
  }
}

export const useStreamApi = (url, options = {}) => {
  const isMounted = useIsMounted()
  const [state, setState] = useState({
    error: null,
    loading: true,
    data: null,
  })
  const [refreshIndex, setRefreshIndex] = useState(0)
  const { telemetryWebTracer } = useAppConfig()
  const { telemetrySpan } = options

  useEffect(
    () => {
      ;(async () => {
        try {
          // Set loading to true
          setState({ ...state, loading: true })
          let data = undefined

          if (telemetryWebTracer && telemetrySpan) {
            const singleSpan = telemetryWebTracer.startSpan(telemetrySpan)
            await context.with(
              trace.setSpan(context.active(), singleSpan),
              async () => {
                data = await streamApi(url, options).then(result => {
                  trace
                    .getSpan(context.active())
                    .addEvent('fetching-single-span-completed')
                  singleSpan.end()

                  return result.data
                })
              }
            )
          } else {
            const result = await streamApi(url, options)
            data = result.data
          }

          if (isMounted.current) {
            setState({
              ...state,
              data,
              error: null,
              loading: false,
            })
          }
        } catch (error) {
          if (isMounted.current) {
            setState({
              ...state,
              data: null,
              error,
              loading: false,
            })

            trace
              .getSpan(context.active())
              .addEvent('fetching-single-span-completed')
          }
        }
      })()
    },
    [url, refreshIndex] // eslint-disable-line
  )

  return {
    ...state,
    updateData: updatedData => setState({ ...state, data: updatedData }),
    refresh: () => setRefreshIndex(Math.random),
  }
}
