import { useEffect, useState } from 'react'

import { useIsMounted } from '@/common/hooks'
import { fetchApi, streamApi } from '@/common/services'
import get from 'lodash/get'

export const useStreamApi = (url, options = {}) => {
  const isMounted = useIsMounted()
  const [state, setState] = useState({
    error: null,
    loading: true,
    data: null,
  })
  const [refreshIndex, setRefreshIndex] = useState(0)
  const apiMethod = get(options, 'streamApi', true) ? streamApi : fetchApi

  useEffect(
    () => {
      ;(async () => {
        try {
          // Set loading to true
          setState({ ...state, loading: true })
          // change of url is watched by effect so base is same and shadow parameter is passed alone
          const { shadowQueryParameter, ...restOptions } = options
          const { data } = await apiMethod(
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
