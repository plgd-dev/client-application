import { useEffect, useState } from 'react'

import { useIsMounted } from '@/common/hooks'
import { fetchApi, streamApi } from '@/common/services'

export const useApi = (url, options = {}) => {
  const isMounted = useIsMounted()
  const [state, setState] = useState({
    error: null,
    loading: true,
    data: null,
  })
  const [refreshIndex, setRefreshIndex] = useState(0)

  useEffect(
    () => {
      ;(async () => {
        try {
          // Set loading to true
          setState({ ...state, loading: true })

          const { data } = await fetchApi(url, options)

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
          }
        }
      })()
    },
    [url, refreshIndex] // eslint-disable-line
  )

  return {
    ...state,
    updateData: updatedData => {
      setState({ ...state, data: updatedData })
    },
    refresh: () => setRefreshIndex(Math.random),
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

  useEffect(
    () => {
      ;(async () => {
        try {
          // Set loading to true
          setState({ ...state, loading: true })
          // change of url is watched by effect so base is same and shadow parameter is passed alone
          const { shadowQueryParameter, ...restOptions } = options
          const { data } = await streamApi(
            shadowQueryParameter ? url + shadowQueryParameter : url,
            restOptions
          )

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
