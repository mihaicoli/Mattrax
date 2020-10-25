// TODO: Pagination, Filters
import { errorForStatus } from './errors'

export const actions = {
  getAll(context: any) {
    return new Promise((resolve, reject) => {
      fetch(process.env.baseUrl + '/devices', {
        headers: new Headers({
          Authorization: 'Bearer ' + context.rootState.authentication.authToken,
        }),
      })
        .then(async (res) => {
          if (res.status !== 200) {
            reject(errorForStatus(res, 'Error fetching devices from server'))
            return
          }

          const devices = await res.json()
          resolve(devices)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getByID(context: any, deviceID: string) {
    return new Promise((resolve, reject) => {
      fetch(process.env.baseUrl + '/device/' + encodeURI(deviceID), {
        headers: new Headers({
          Authorization: 'Bearer ' + context.rootState.authentication.authToken,
        }),
      })
        .then(async (res) => {
          if (res.status !== 200) {
            reject(errorForStatus(res, 'Error fetching device from server'))
            return
          }

          const device = await res.json()
          resolve(device)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getInformationByID(context: any, deviceID: string) {
    return new Promise((resolve, reject) => {
      fetch(process.env.baseUrl + '/device/' + encodeURI(deviceID) + '/info', {
        headers: new Headers({
          Authorization: 'Bearer ' + context.rootState.authentication.authToken,
        }),
      })
        .then(async (res) => {
          if (res.status !== 200) {
            reject(
              errorForStatus(
                res,
                'Error fetching device information from server'
              )
            )
            return
          }

          const deviceInfo = await res.json()
          resolve(deviceInfo)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getScopeByID(context: any, deviceID: string) {
    return new Promise((resolve, reject) => {
      fetch(process.env.baseUrl + '/device/' + encodeURI(deviceID) + '/scope', {
        headers: new Headers({
          Authorization: 'Bearer ' + context.rootState.authentication.authToken,
        }),
      })
        .then(async (res) => {
          if (res.status !== 200) {
            reject(
              errorForStatus(res, 'Error fetching device scope from server')
            )
            return
          }

          const deviceScope = await res.json()
          resolve(deviceScope)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
}
